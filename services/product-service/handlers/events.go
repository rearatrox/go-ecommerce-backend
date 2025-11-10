package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/product-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Product-Handlers, die beim Aufruf von Routen /products aufgerufen werden (Verarbeitung der Requests)
// GetProducts godoc
// @Summary      Get all products
// @Description  Get all product information of all products
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Product
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products [get]
func GetProducts(context *gin.Context) {

	l := logger.FromContext(context.Request.Context())
	l.Debug("GetProducts called")

	products, err := models.GetProducts()

	if err != nil {
		l.Error("failed to fetch products", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch products.", "error": err.Error()})
		return
	}

	l.Info("fetched products", "count", len(products))
	//Response in JSON
	context.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary      Get single product by ID
// @Description  Get details of an product by its ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func GetProduct(context *gin.Context) {
	var product *models.Product
	productId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	l := logger.FromContext(context.Request.Context())

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse product id.", "error": err.Error()})
		return
	}

	l.Debug("GetProduct called", "product_id", productId)

	product, err = models.GetProductByID(productId)
	if err != nil {
		l.Error("failed to fetch product", "product_id", productId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	l.Info("fetched product", "product_id", productId)
	//Response in JSON
	context.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create product. Requires authentication.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Product payload"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     JWT
// @Router       /products [post]
func CreateProduct(context *gin.Context) {

	l := logger.FromContext(context.Request.Context())
	l.Debug("CreateProduct called")

	var product models.Product
	err := context.ShouldBindJSON(&product)

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := context.GetInt64("userId")
	product.CreatorID = userId

	err = product.InsertProduct()

	if err != nil {
		l.Error("failed to save product", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create product.", "error": err.Error()})
		return
	}
	l.Info("created product", "product_id", product.ID, "creator_id", userId)
	context.JSON(http.StatusCreated, gin.H{"message": "Product created", "product": product})
}

// UpdateProduct godoc
// @Summary      Update an existing product
// @Description  Update product by ID. Requires authentication and ownership.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id     path      int           true  "Product ID"
// @Param        product  body      models.Product  true  "Updated product payload"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     JWT
// @Router       /products/{id} [put]
func UpdateProduct(context *gin.Context) {
	var product *models.Product
	productId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	l := logger.FromContext(context.Request.Context())

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse product id to update product.", "error": err.Error()})
		return
	}

	l.Debug("UpdateProduct called", "product_id", productId)

	product, err = models.GetProductByID(productId)
	if err != nil {
		l.Error("failed to fetch product", "product_id", productId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	userId, _ := context.Get("userId")

	if product.CreatorID != userId {
		l.Warn("unauthorized update attempt", "product_id", productId, "user", userId)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized to update that product"})
		return
	}

	var updatedProduct models.Product
	err = context.ShouldBindJSON(&updatedProduct)
	if err != nil {
		l.Warn("failed to parse update payload", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request data.", "error": err.Error()})
		return
	}

	updatedProduct.ID = product.ID
	err = updatedProduct.UpdateProduct()
	if err != nil {
		l.Error("failed to update product", "product_id", productId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update product.", "error": err.Error()})
		return
	}
	//Response in JSON
	l.Info("updated product", "product_id", productId)
	context.JSON(http.StatusOK, gin.H{"message": "updated product successfully", "updatedProduct": updatedProduct})
}

// DeleteProduct godoc
// @Summary      Delete an product
// @Description  Delete product by ID. Requires authentication and ownership.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Product ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     JWT
// @Router       /products/{id} [delete]
func DeleteProductByID(context *gin.Context) {
	productId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	l := logger.FromContext(context.Request.Context())

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse product id to update product.", "error": err.Error()})
		return
	}

	l.Debug("DeleteProduct called", "product_id", productId)

	deleteProduct, err := models.GetProductByID(productId)
	if err != nil {
		l.Error("failed to fetch product", "product_id", productId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	userId, _ := context.Get("userId")

	if deleteProduct.CreatorID != userId {
		l.Warn("unauthorized delete attempt", "product_id", productId, "user", userId)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized to delete that product"})
		return
	}

	err = deleteProduct.DeleteProductByID()
	if err != nil {
		l.Error("failed to delete product", "product_id", productId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete product.", "error": err.Error()})
		return
	}
	//Response in JSON
	l.Info("deleted product", "product_id", productId)
	context.JSON(http.StatusOK, gin.H{"message": "deleted product successfully", "deletedProduct": deleteProduct})
}
