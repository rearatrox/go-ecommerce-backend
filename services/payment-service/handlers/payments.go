package handlers

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"strconv"

	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/payment-service/models"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/webhook"
)

type CreatePaymentIntentRequest struct {
	OrderID int64 `json:"orderId" binding:"required" example:"1"`
}

type CreatePaymentIntentResponse struct {
	PaymentID     int64  `json:"paymentId" example:"1"`
	ClientSecret  string `json:"clientSecret" example:"pi_xxx_secret_xxx"`
	PaymentIntent string `json:"paymentIntentId" example:"pi_1234567890"`
}

// CreatePaymentIntent godoc
// @Summary      Create payment intent
// @Description  Creates a Stripe Payment Intent for an order
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request  body      CreatePaymentIntentRequest  true  "Order ID"
// @Success      201      {object}  CreatePaymentIntentResponse
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      409      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /payment-intents [post]
func CreatePaymentIntent(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("CreatePaymentIntent called", "user_id", userId)

	var req CreatePaymentIntentRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		l.Error("failed to bind request", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid request.", "error": err.Error()})
		return
	}

	// Check if payment already exists for this order
	existingPayment, err := models.GetByOrderID(req.OrderID)
	if err == nil && existingPayment != nil {
		l.Warn("payment already exists for order", "order_id", req.OrderID, "payment_id", existingPayment.ID)
		context.JSON(http.StatusConflict, gin.H{
			"message":         "payment already exists for this order",
			"paymentId":       existingPayment.ID,
			"clientSecret":    existingPayment.StripeClientSecret,
			"paymentIntentId": existingPayment.StripePaymentIntentID,
		})
		return
	}

	// Get order details from order-service
	order, err := getOrderDetails(req.OrderID)
	if err != nil {
		l.Error("failed to get order details", "order_id", req.OrderID, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch order details.", "error": err.Error()})
		return
	}

	// Verify user owns this order
	if order.UserID != userId {
		l.Warn("unauthorized order access", "order_id", req.OrderID, "user_id", userId, "order_user_id", order.UserID)
		context.JSON(http.StatusForbidden, gin.H{"message": "access denied."})
		return
	}

	// Verify order is in pending state
	if order.Status != "pending" {
		l.Warn("order not in pending state", "order_id", req.OrderID, "status", order.Status)
		context.JSON(http.StatusConflict, gin.H{"message": "order is not in pending state.", "status": order.Status})
		return
	}

	amountCents := order.TotalCents

	// Create Stripe Payment Intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amountCents)),
		Currency: stripe.String("eur"),
		Params: stripe.Params{
			Metadata: map[string]string{
				"order_id": strconv.FormatInt(req.OrderID, 10),
				"user_id":  strconv.FormatInt(userId, 10),
			},
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		l.Error("failed to create stripe payment intent", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create payment intent.", "error": err.Error()})
		return
	}

	// Save payment to database
	payment := &models.Payment{
		OrderID:               req.OrderID,
		UserID:                userId,
		AmountCents:           amountCents,
		Currency:              "EUR",
		Status:                "pending",
		StripePaymentIntentID: &pi.ID,
		StripeClientSecret:    &pi.ClientSecret,
	}

	if err := models.Create(payment); err != nil {
		l.Error("failed to create payment record", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save payment.", "error": err.Error()})
		return
	}

	l.Info("created payment intent", "payment_id", payment.ID, "order_id", req.OrderID, "amount_cents", amountCents)
	context.JSON(http.StatusCreated, CreatePaymentIntentResponse{
		PaymentID:     payment.ID,
		ClientSecret:  pi.ClientSecret,
		PaymentIntent: pi.ID,
	})
}

// GetPaymentStatus godoc
// @Summary      Get payment status
// @Description  Get the status of a payment by ID
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Payment ID"
// @Success      200  {object}  models.Payment
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /payments/{id} [get]
func GetPaymentStatus(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")

	paymentID, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		l.Error("invalid payment id", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid payment id."})
		return
	}

	payment, err := models.GetByID(paymentID)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Warn("payment not found", "payment_id", paymentID)
			context.JSON(http.StatusNotFound, gin.H{"message": "payment not found."})
			return
		}
		l.Error("failed to get payment", "payment_id", paymentID, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get payment.", "error": err.Error()})
		return
	}

	// Verify user owns this payment
	if payment.UserID != userId {
		l.Warn("unauthorized payment access", "payment_id", paymentID, "user_id", userId, "payment_user_id", payment.UserID)
		context.JSON(http.StatusForbidden, gin.H{"message": "access denied."})
		return
	}

	l.Debug("retrieved payment", "payment_id", paymentID, "status", payment.Status)
	context.JSON(http.StatusOK, payment)
}

// WebhookHandler godoc
// @Summary      Stripe webhook handler
// @Description  Handles Stripe webhook events for payment updates
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /webhooks/stripe [post]
func WebhookHandler(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())

	payload, err := io.ReadAll(context.Request.Body)
	if err != nil {
		l.Error("failed to read webhook payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload."})
		return
	}

	// Verify webhook signature
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(payload, context.GetHeader("Stripe-Signature"), webhookSecret)
	if err != nil {
		l.Error("failed to verify webhook signature", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid signature."})
		return
	}

	l.Info("received stripe webhook", "event_type", event.Type)

	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := pi.UnmarshalJSON(event.Data.Raw); err != nil {
			l.Error("failed to unmarshal payment intent", "error", err)
			context.JSON(http.StatusBadRequest, gin.H{"message": "invalid event data."})
			return
		}

		// Update payment status
		payment, err := models.GetByStripePaymentIntentID(pi.ID)
		if err != nil {
			l.Error("failed to find payment", "payment_intent_id", pi.ID, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "payment not found."})
			return
		}

		if err := models.UpdateStatus(payment.ID, "succeeded"); err != nil {
			l.Error("failed to update payment status", "payment_id", payment.ID, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update payment."})
			return
		}

		l.Info("payment succeeded", "payment_id", payment.ID, "order_id", payment.OrderID)

		// Update order status to "confirmed"
		if err := updateOrderStatus(payment.OrderID, "confirmed"); err != nil {
			l.Error("failed to update order status", "order_id", payment.OrderID, "error", err)
			// Don't return error - payment was successful, this is just a notification issue
		} else {
			l.Info("order confirmed", "order_id", payment.OrderID)
		}

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := pi.UnmarshalJSON(event.Data.Raw); err != nil {
			l.Error("failed to unmarshal payment intent", "error", err)
			context.JSON(http.StatusBadRequest, gin.H{"message": "invalid event data."})
			return
		}

		payment, err := models.GetByStripePaymentIntentID(pi.ID)
		if err != nil {
			l.Error("failed to find payment", "payment_intent_id", pi.ID, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "payment not found."})
			return
		}

		if err := models.UpdateStatus(payment.ID, "failed"); err != nil {
			l.Error("failed to update payment status", "payment_id", payment.ID, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update payment."})
			return
		}

		l.Info("payment failed", "payment_id", payment.ID, "order_id", payment.OrderID)

	case "payment_intent.canceled":
		var pi stripe.PaymentIntent
		if err := pi.UnmarshalJSON(event.Data.Raw); err != nil {
			l.Error("failed to unmarshal payment intent", "error", err)
			context.JSON(http.StatusBadRequest, gin.H{"message": "invalid event data."})
			return
		}

		payment, err := models.GetByStripePaymentIntentID(pi.ID)
		if err != nil {
			l.Error("failed to find payment", "payment_intent_id", pi.ID, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "payment not found."})
			return
		}

		if err := models.UpdateStatus(payment.ID, "cancelled"); err != nil {
			l.Error("failed to update payment status", "payment_id", payment.ID, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update payment."})
			return
		}

		l.Info("payment cancelled", "payment_id", payment.ID, "order_id", payment.OrderID)

	default:
		l.Debug("unhandled webhook event type", "event_type", event.Type)
	}

	context.JSON(http.StatusOK, gin.H{"received": true})
}
