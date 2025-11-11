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

	// Get JWT token for service-to-service calls
	jwtToken := context.GetHeader("Authorization")
	if jwtToken == "" {
		l.Error("missing authorization token", "user_id", userId)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "authorization required."})
		return
	}

	// Verify address ownership if addresses are provided
	if req.ShippingAddressID != nil {
		if err := verifyAddressOwnership(*req.ShippingAddressID, userId, jwtToken); err != nil {
			l.Warn("invalid shipping address", "user_id", userId, "address_id", *req.ShippingAddressID, "error", err)
			context.JSON(http.StatusForbidden, gin.H{"message": "invalid shipping address.", "error": err.Error()})
			return
		}
	}

	if req.BillingAddressID != nil {
		if err := verifyAddressOwnership(*req.BillingAddressID, userId, jwtToken); err != nil {
			l.Warn("invalid billing address", "user_id", userId, "address_id", *req.BillingAddressID, "error", err)
			context.JSON(http.StatusForbidden, gin.H{"message": "invalid billing address.", "error": err.Error()})
			return
		}
	}

	// Get cart items to check stock before creating order
	cartItems, err := models.GetCartItemsForUser(userId)
	if err != nil {
		l.Error("failed to get cart items", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch cart items.", "error": err.Error()})
		return
	}

	if len(cartItems) == 0 {
		l.Warn("empty cart", "user_id", userId)
		context.JSON(http.StatusBadRequest, gin.H{"message": "cannot create order from empty cart."})
		return
	}

	// Check stock availability for all items
	for _, item := range cartItems {
		stockResp, err := checkStockAvailability(item.ProductID, item.Quantity)
		if err != nil {
			l.Error("failed to check stock", "product_id", item.ProductID, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not check stock availability.", "error": err.Error()})
			return
		}

		if !stockResp.Available {
			l.Warn("insufficient stock for order", "product_id", item.ProductID, "requested", item.Quantity, "available", stockResp.AvailableQty)
			context.JSON(http.StatusConflict, gin.H{
				"message":     "insufficient stock",
				"productId":   item.ProductID,
				"productName": item.ProductName,
				"requested":   item.Quantity,
				"available":   stockResp.AvailableQty,
			})
			return
		}
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

	// If status changes to 'confirmed', reduce stock
	if req.Status == "confirmed" && order.Status != "confirmed" {
		l.Debug("order confirmed, reducing stock", "order_id", orderId)

		if err := reduceStockForOrder(order.Items); err != nil {
			l.Error("failed to reduce stock", "order_id", orderId, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not reduce stock.", "error": err.Error()})
			return
		}
		l.Info("stock reduced", "order_id", orderId, "items_count", len(order.Items))
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

// CancelOrder godoc
// @Summary      Cancel an order
// @Description  Cancel an order (only possible for pending orders)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  models.Order
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /orders/{id}/cancel [patch]
func CancelOrder(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	orderIdStr := context.Param("id")

	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		l.Error("invalid order ID", "order_id", orderIdStr, "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid order ID."})
		return
	}

	l.Debug("CancelOrder called", "user_id", userId, "order_id", orderId)

	// Get order
	order, err := models.GetOrderByID(orderId, userId)
	if err != nil {
		l.Error("failed to get order", "user_id", userId, "order_id", orderId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "order not found."})
		return
	}

	// Only allow cancellation of pending or confirmed orders (not shipped/delivered)
	if order.Status != "pending" && order.Status != "confirmed" {
		l.Warn("cannot cancel order in current state", "order_id", orderId, "status", order.Status)
		context.JSON(http.StatusConflict, gin.H{
			"message": "order cannot be cancelled in current state",
			"status":  order.Status,
		})
		return
	}

	// Update status to cancelled
	if err := order.UpdateStatus("cancelled"); err != nil {
		l.Error("failed to cancel order", "user_id", userId, "order_id", orderId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not cancel order.", "error": err.Error()})
		return
	}

	l.Info("cancelled order", "user_id", userId, "order_id", orderId)
	context.JSON(http.StatusOK, order)
}

// InternalUpdateOrderStatus godoc
// @Summary      Internal order status update
// @Description  Updates order status without authentication (for service-to-service calls from payment-service)
// @Tags         Internal
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "Order ID"
// @Param        request  body      UpdateStatusRequest  true  "New status"
// @Success      200      {object}  models.Order
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /internal/orders/{id}/status [patch]
func InternalUpdateOrderStatus(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
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

	l.Debug("InternalUpdateOrderStatus called", "order_id", orderId, "new_status", req.Status)

	// Get order without user validation (internal call)
	order, err := models.GetOrderByIDInternal(orderId)
	if err != nil {
		l.Error("failed to get order", "order_id", orderId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "order not found."})
		return
	}

	// If status changes to 'confirmed', reduce stock
	if req.Status == "confirmed" && order.Status != "confirmed" {
		l.Debug("order confirmed, reducing stock (internal)", "order_id", orderId)

		if err := reduceStockForOrder(order.Items); err != nil {
			l.Error("failed to reduce stock", "order_id", orderId, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "could not reduce stock.", "error": err.Error()})
			return
		}
		l.Info("stock reduced (internal)", "order_id", orderId, "items_count", len(order.Items))
	}

	// Update status
	if err := order.UpdateStatus(req.Status); err != nil {
		l.Error("failed to update order status", "order_id", orderId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update order status.", "error": err.Error()})
		return
	}

	l.Info("updated order status (internal)", "order_id", orderId, "new_status", req.Status)
	context.JSON(http.StatusOK, order)
}
