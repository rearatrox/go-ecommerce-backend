package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/cart-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCart godoc
// @Summary      Get user's cart
// @Description  Get the active cart for the authenticated user with all items
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Cart
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /cart [get]
func GetCart(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("GetCart called", "user_id", userId)

	cart, err := models.GetOrCreateCart(userId)
	if err != nil {
		l.Error("failed to get cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch cart.", "error": err.Error()})
		return
	}

	l.Info("fetched cart", "user_id", userId, "cart_id", cart.ID, "items_count", len(cart.Items))
	context.JSON(http.StatusOK, cart)
}

// AddItem godoc
// @Summary      Add item to cart
// @Description  Add a product to the cart or update quantity if it already exists
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        request  body      models.AddItemRequest  true  "Product and quantity to add"
// @Success      200      {object}  models.Cart
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /cart/items [post]
func AddItem(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("AddItem called", "user_id", userId)

	var req models.AddItemRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		l.Error("failed to bind request", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid request.", "error": err.Error()})
		return
	}

	// Check stock availability
	stockResp, err := checkStockAvailability(int64(req.ProductID), req.Quantity)
	if err != nil {
		l.Error("failed to check stock", "product_id", req.ProductID, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not check stock availability.", "error": err.Error()})
		return
	}

	if !stockResp.Available {
		l.Warn("insufficient stock", "product_id", req.ProductID, "requested", req.Quantity, "available", stockResp.AvailableQty)
		context.JSON(http.StatusConflict, gin.H{
			"message":   "insufficient stock",
			"requested": req.Quantity,
			"available": stockResp.AvailableQty,
			"productId": req.ProductID,
		})
		return
	}

	// Get or create cart
	cart, err := models.GetOrCreateCart(userId)
	if err != nil {
		l.Error("failed to get cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch cart.", "error": err.Error()})
		return
	}

	// Add item to cart
	cartItem := &models.CartItem{
		CartID:    cart.ID,
		ProductID: int64(req.ProductID),
		Quantity:  req.Quantity,
	}

	if err := cartItem.AddOrUpdate(); err != nil {
		l.Error("failed to add item to cart", "user_id", userId, "product_id", req.ProductID, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not add item to cart.", "error": err.Error()})
		return
	}

	// Reload cart with updated items
	if err := cart.Reload(); err != nil {
		l.Error("failed to reload cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not reload cart.", "error": err.Error()})
		return
	}

	l.Info("added item to cart", "user_id", userId, "cart_id", cart.ID, "product_id", req.ProductID, "quantity", req.Quantity)
	context.JSON(http.StatusOK, cart)
}

// UpdateItem godoc
// @Summary      Update cart item quantity
// @Description  Update the quantity of a specific product in the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        productId  path      int                       true  "Product ID"
// @Param        request    body      models.UpdateItemRequest  true  "New quantity"
// @Success      200        {object}  models.Cart
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /cart/items/{productId} [put]
func UpdateItem(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("UpdateItem called", "user_id", userId)

	productId, err := strconv.ParseInt(context.Param("productId"), 10, 64)
	if err != nil {
		l.Error("invalid product ID", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid product ID.", "error": err.Error()})
		return
	}

	var req models.UpdateItemRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		l.Error("failed to bind request", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid request.", "error": err.Error()})
		return
	}

	// Get cart
	cart, err := models.GetOrCreateCart(userId)
	if err != nil {
		l.Error("failed to get cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch cart.", "error": err.Error()})
		return
	}

	// Update item quantity
	cartItem := &models.CartItem{
		CartID:    cart.ID,
		ProductID: productId,
		Quantity:  req.Quantity,
	}

	if err := cartItem.UpdateQuantity(); err != nil {
		l.Error("failed to update item quantity", "user_id", userId, "product_id", productId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "item not found in cart.", "error": err.Error()})
		return
	}

	// Reload cart with updated items
	if err := cart.Reload(); err != nil {
		l.Error("failed to reload cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not reload cart.", "error": err.Error()})
		return
	}

	l.Info("updated item in cart", "user_id", userId, "cart_id", cart.ID, "product_id", productId, "quantity", req.Quantity)
	context.JSON(http.StatusOK, cart)
}

// RemoveItem godoc
// @Summary      Remove item from cart
// @Description  Remove a product from the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        productId  path      int  true  "Product ID"
// @Success      200        {object}  models.Cart
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /cart/items/{productId} [delete]
func RemoveItem(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("RemoveItem called", "user_id", userId)

	productId, err := strconv.ParseInt(context.Param("productId"), 10, 64)
	if err != nil {
		l.Error("invalid product ID", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid product ID.", "error": err.Error()})
		return
	}

	// Get cart
	cart, err := models.GetOrCreateCart(userId)
	if err != nil {
		l.Error("failed to get cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch cart.", "error": err.Error()})
		return
	}

	// Remove item
	cartItem := &models.CartItem{
		CartID:    cart.ID,
		ProductID: productId,
	}

	if err := cartItem.Remove(); err != nil {
		l.Error("failed to remove item from cart", "user_id", userId, "product_id", productId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "item not found in cart.", "error": err.Error()})
		return
	}

	// Reload cart with updated items
	if err := cart.Reload(); err != nil {
		l.Error("failed to reload cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not reload cart.", "error": err.Error()})
		return
	}

	l.Info("removed item from cart", "user_id", userId, "cart_id", cart.ID, "product_id", productId)
	context.JSON(http.StatusOK, cart)
}

// ClearCart godoc
// @Summary      Clear cart
// @Description  Remove all items from the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Cart
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /cart [delete]
func ClearCart(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("ClearCart called", "user_id", userId)

	// Get cart
	cart, err := models.GetOrCreateCart(userId)
	if err != nil {
		l.Error("failed to get cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch cart.", "error": err.Error()})
		return
	}

	// Clear cart
	if err := cart.Clear(); err != nil {
		l.Error("failed to clear cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not clear cart.", "error": err.Error()})
		return
	}

	// Reload cart (now empty)
	if err := cart.Reload(); err != nil {
		l.Error("failed to reload cart", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not reload cart.", "error": err.Error()})
		return
	}

	l.Info("cleared cart", "user_id", userId, "cart_id", cart.ID)
	context.JSON(http.StatusOK, cart)
}
