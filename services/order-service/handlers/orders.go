package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/order-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	ShippingAddressID *int64 `json:"shippingAddressId" example:"1"`
	BillingAddressID  *int64 `json:"billingAddressId" example:"2"`
}

// CreateOrder godoc
// @Summary      Create order from cart
// @Description  Creates a new order from the user's active cart and marks cart as ordered
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request  body      CreateOrderRequest  true  "Address IDs (optional)"
// @Success      201      {object}  models.Order
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /orders [post]
func CreateOrder(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("CreateOrder called", "user_id", userId)

	var req CreateOrderRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		l.Error("failed to bind request", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid request.", "error": err.Error()})
		return
	}

	// Create order from active cart
	order, err := models.CreateFromCart(userId, req.ShippingAddressID, req.BillingAddressID)
	if err != nil {
		l.Error("failed to create order", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create order.", "error": err.Error()})
		return
	}

	l.Info("created order", "user_id", userId, "order_id", order.ID, "total_cents", order.TotalCents)
	context.JSON(http.StatusCreated, order)
}

// GetOrder godoc
// @Summary      Get order by ID
// @Description  Get details of a specific order including items and addresses
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  models.Order
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /orders/{id} [get]
func GetOrder(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	orderIdStr := context.Param("id")

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		l.Error("invalid order ID", "order_id", orderIdStr, "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid order ID."})
		return
	}

	l.Debug("GetOrder called", "user_id", userId, "order_id", orderId)

	order, err := models.GetOrderByID(orderId, userId)
	if err != nil {
		l.Error("failed to get order", "user_id", userId, "order_id", orderId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "order not found."})
		return
	}

	l.Info("fetched order", "user_id", userId, "order_id", orderId)
	context.JSON(http.StatusOK, order)
}

// ListOrders godoc
// @Summary      List user's orders
// @Description  Get all orders for the authenticated user
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Order
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /orders [get]
func ListOrders(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("ListOrders called", "user_id", userId)

	orders, err := models.GetUserOrders(userId)
	if err != nil {
		l.Error("failed to get orders", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch orders.", "error": err.Error()})
		return
	}

	l.Info("fetched orders", "user_id", userId, "count", len(orders))
	context.JSON(http.StatusOK, orders)
}

type UpdateStatusRequest struct {
	Status string `json:"status" example:"confirmed" binding:"required"`
}

// UpdateOrderStatus godoc
// @Summary      Update order status
// @Description  Update the status of an order (for future payment integration)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Order ID"
// @Param        request  body      UpdateStatusRequest  true  "New status"
// @Success      200      {object}  models.Order
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /orders/{id}/status [patch]
func UpdateOrderStatus(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	orderIdStr := context.Param("id")

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		l.Error("invalid order ID", "order_id", orderIdStr, "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid order ID."})
		return
	}

	var req UpdateStatusRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		l.Error("failed to bind request", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid request.", "error": err.Error()})
		return
	}

	l.Debug("UpdateOrderStatus called", "user_id", userId, "order_id", orderId, "new_status", req.Status)

	// Get order
	order, err := models.GetOrderByID(orderId, userId)
	if err != nil {
		l.Error("failed to get order", "user_id", userId, "order_id", orderId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "order not found."})
		return
	}

	// Update status
	if err := order.UpdateStatus(req.Status); err != nil {
		l.Error("failed to update order status", "user_id", userId, "order_id", orderId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update order status.", "error": err.Error()})
		return
	}

	l.Info("updated order status", "user_id", userId, "order_id", orderId, "new_status", req.Status)
	context.JSON(http.StatusOK, order)
}
