package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/product-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProducts godoc
// @Summary      Get all products
// @Description  Get all product information of all products
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Product
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
// @Router       /products/id/{id} [get]
func GetProductByID(context *gin.Context) {
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

// GetProduct godoc
// @Summary      Get single product by SKU
// @Description  Get details of an product by its SKU
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        sku   path      string  true  "Product SKU"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/sku/{sku} [get]
func GetProductBySKU(context *gin.Context) {
	var product *models.Product
	productSku := context.Param("sku")

	l := logger.FromContext(context.Request.Context())
	l.Debug("GetProduct called", "SKU", productSku)

	product, err := models.GetProductBySKU(productSku)
	if err != nil {
		l.Error("failed to fetch product", "SKU", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	l.Info("fetched product", "SKU", productSku)
	//Response in JSON
	context.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create product with optional category assignment. Requires authentication and authorization as role "admin".
// @Tags         Products (Admin)
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Product payload - Optional field: categoryIds (array of integers, e.g. [1, 2, 3])"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/products/create [post]
func CreateProduct(context *gin.Context) {

	l := logger.FromContext(context.Request.Context())
	l.Debug("CreateProduct called")

	var requestBody struct {
		models.Product
		CategoryIds []int64 `json:"categoryIds,omitempty"`
	}
	err := context.ShouldBindJSON(&requestBody)

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := context.GetInt64("userId")
	requestBody.Product.CreatorID = userId

	// Produkt anlegen
	err = requestBody.Product.InsertProduct()

	if err != nil {
		l.Error("failed to save product", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create product.", "error": err.Error()})
		return
	}

	// Kategorien hinzufügen (falls vorhanden)
	if len(requestBody.CategoryIds) > 0 {
		err = requestBody.Product.AddCategories(requestBody.CategoryIds)
		if err != nil {
			l.Error("failed to add categories to product", "product_id", requestBody.Product.ID, "error", err)
			context.JSON(http.StatusInternalServerError, gin.H{"message": "product created but could not add categories.", "error": err.Error()})
			return
		}
		l.Info("added categories to product", "product_id", requestBody.Product.ID, "category_count", len(requestBody.CategoryIds))
	}

	l.Info("created product", "product_id", requestBody.Product.ID, "creator_id", userId)
	context.JSON(http.StatusCreated, gin.H{"message": "Product created", "product": requestBody.Product})
}

// UpdateProduct godoc
// @Summary      Update an existing product
// @Description  Update product by sku. Note: SKU in body must match SKU in URL path. Requires authentication and authorization as role "admin".
// @Tags         Products (Admin)
// @Accept       json
// @Produce      json
// @Param        sku      path      string           true  "Product SKU"
// @Param        product  body      models.Product   true  "Updated product payload"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/products/update/{sku} [put]
func UpdateProduct(context *gin.Context) {
	var product *models.Product
	productSku := context.Param("sku")

	l := logger.FromContext(context.Request.Context())
	l.Debug("UpdateProduct called", "productSku", productSku)

	product, err := models.GetProductBySKU(productSku)
	if err != nil {
		l.Error("failed to fetch product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	var updatedProduct models.Product
	err = context.ShouldBindJSON(&updatedProduct)
	if err != nil {
		l.Warn("failed to parse update payload", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request data.", "error": err.Error()})
		return
	}

	if updatedProduct.SKU != productSku {
		l.Error("sku values are not equal (either route parameter or request sku is false)", "updatedProduct_sku", updatedProduct.SKU, "routeParam sku", productSku, "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "sku values are not equal (either route parameter or request sku is false)"})
		return
	}

	userId, _ := context.Get("userId")
	updatedProduct.ID = product.ID
	updatedProduct.CreatorID = product.CreatorID
	updatedProduct.UpdatorID = userId.(int64)

	err = updatedProduct.UpdateProduct()
	if err != nil {
		l.Error("failed to update product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update product.", "error": err.Error()})
		return
	}
	//Response in JSON
	l.Info("updated product", "productSku", productSku)
	context.JSON(http.StatusOK, gin.H{"message": "updated product successfully", "updatedProduct": updatedProduct})
}

// DeactivateProduct godoc
// @Summary      Deactivate a product
// @Description  Deactivate product by sku. Requires authentication and authorization as role "admin".
// @Tags         Products (Admin)
// @Accept       json
// @Produce      json
// @Param        sku   path  string  true  "Product sku"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/products/deactivate/{sku} [post]
func DeactivateProductBySKU(context *gin.Context) {
	productSku := context.Param("sku")
	l := logger.FromContext(context.Request.Context())
	l.Debug("DeactivateProduct called", "productSku", productSku)

	deactivatedProduct, err := models.GetProductBySKU(productSku)
	if err != nil {
		l.Error("failed to fetch product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	err = deactivatedProduct.DeactivateProductBySKU()
	if err != nil {
		l.Error("failed to deactivate product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not deactivate product.", "error": err.Error()})
		return
	}
	//Response in JSON
	l.Info("deactivated product", "productSku", productSku)
	context.JSON(http.StatusOK, gin.H{"message": "deactivated product successfully", "deactivatedProduct": deactivatedProduct})
}

// DeleteProduct godoc
// @Summary      Delete an product
// @Description  Delete product by sku. Requires authentication and authorization as role "admin".
// @Tags         Products (Admin)
// @Accept       json
// @Produce      json
// @Param        sku   path  string  true  "Product sku"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/products/delete/{sku} [delete]
func DeleteProductBySKU(context *gin.Context) {
	productSku := context.Param("sku")
	l := logger.FromContext(context.Request.Context())
	l.Debug("DeleteProduct called", "productSku", productSku)

	deleteProduct, err := models.GetProductBySKU(productSku)
	if err != nil {
		l.Error("failed to fetch product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	err = deleteProduct.DeleteProductBySKU()
	if err != nil {
		l.Error("failed to delete product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete product.", "error": err.Error()})
		return
	}
	//Response in JSON
	l.Info("deleted product", "productSku", productSku)
	context.JSON(http.StatusOK, gin.H{"message": "deleted product successfully", "deletedProduct": deleteProduct})
}

// GetProductCategories godoc
// @Summary      Get categories of a product
// @Description  Get all categories assigned to a specific product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        sku   path      string  true  "Product SKU"
// @Success      200  {array}   models.Category
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{sku}/categories [get]
func GetProductCategories(context *gin.Context) {
	productSku := context.Param("sku")
	l := logger.FromContext(context.Request.Context())
	l.Debug("GetProductCategories called", "productSku", productSku)

	product, err := models.GetProductBySKU(productSku)
	if err != nil {
		l.Error("failed to fetch product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	categories, err := models.GetProductCategories(product.ID)
	if err != nil {
		l.Error("failed to fetch product categories", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product categories.", "error": err.Error()})
		return
	}

	l.Info("fetched product categories", "productSku", productSku, "count", len(categories))
	context.JSON(http.StatusOK, categories)
}

// AddCategoriesToProduct godoc
// @Summary      Add categories to a product
// @Description  Add one or more categories to a product (as an Array of CategoryIds). Requires authentication and authorization as role "admin".
// @Tags         Products (Admin)
// @Accept       json
// @Produce      json
// @Param        sku   path      string  true  "Product SKU"
// @Param        body  body      object  true  "Category IDs - Example:"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/products/{sku}/categories [post]
func AddCategoriesToProduct(context *gin.Context) {
	productSku := context.Param("sku")
	l := logger.FromContext(context.Request.Context())
	l.Debug("AddCategoriesToProduct called", "productSku", productSku)

	var requestBody struct {
		CategoryIds []int64 `json:"categoryIds" binding:"required"`
	}
	err := context.ShouldBindJSON(&requestBody)
	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := models.GetProductBySKU(productSku)
	if err != nil {
		l.Error("failed to fetch product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	err = product.AddCategories(requestBody.CategoryIds)
	if err != nil {
		l.Error("failed to add categories to product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not add categories to product.", "error": err.Error()})
		return
	}

	l.Info("added categories to product", "productSku", productSku, "category_count", len(requestBody.CategoryIds))
	context.JSON(http.StatusOK, gin.H{"message": "categories added successfully", "productSku": productSku, "categoryIds": requestBody.CategoryIds})
}

// RemoveCategoryFromProduct godoc
// @Summary      Remove a category from a product
// @Description  Remove a category assignment from a product. Requires authentication and authorization as role "admin".
// @Tags         Products (Admin)
// @Accept       json
// @Produce      json
// @Param        sku         path  string  true  "Product SKU"
// @Param        categoryId  path  int     true  "Category ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/products/{sku}/categories/{categoryId} [delete]
func RemoveCategoryFromProduct(context *gin.Context) {
	productSku := context.Param("sku")
	categoryId, err := strconv.ParseInt(context.Param("categoryId"), 10, 64)

	l := logger.FromContext(context.Request.Context())

	if err != nil {
		l.Warn("invalid category id", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse category id.", "error": err.Error()})
		return
	}

	l.Debug("RemoveCategoryFromProduct called", "productSku", productSku, "categoryId", categoryId)

	product, err := models.GetProductBySKU(productSku)
	if err != nil {
		l.Error("failed to fetch product", "productSku", productSku, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch product.", "error": err.Error()})
		return
	}

	err = product.RemoveCategory(categoryId)
	if err != nil {
		l.Error("failed to remove category from product", "productSku", productSku, "categoryId", categoryId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not remove category from product.", "error": err.Error()})
		return
	}

	l.Info("removed category from product", "productSku", productSku, "categoryId", categoryId)
	context.JSON(http.StatusOK, gin.H{"message": "category removed successfully", "productSku": productSku, "categoryId": categoryId})
}

type CheckStockRequest struct {
	ProductID int64 `json:"productId" binding:"required" example:"1"`
	Quantity  int   `json:"quantity" binding:"required,min=1" example:"2"`
}

type CheckStockResponse struct {
	Available    bool  `json:"available" example:"true"`
	RequestedQty int   `json:"requestedQty" example:"2"`
	AvailableQty int   `json:"availableQty" example:"10"`
	ProductID    int64 `json:"productId" example:"1"`
}

// CheckStock godoc
// @Summary      Check stock availability
// @Description  Check if enough stock is available for a product (used by Cart/Order services)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        request  body      CheckStockRequest  true  "Product and quantity to check"
// @Success      200      {object}  CheckStockResponse
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /products/stock/check [post]
func CheckStock(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())

	var req CheckStockRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		l.Error("failed to bind request", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid request.", "error": err.Error()})
		return
	}

	l.Debug("CheckStock called", "productId", req.ProductID, "quantity", req.Quantity)

	available, currentStock, err := models.CheckStockAvailable(req.ProductID, req.Quantity)
	if err != nil {
		l.Error("failed to check stock", "productId", req.ProductID, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not check stock.", "error": err.Error()})
		return
	}

	response := CheckStockResponse{
		Available:    available,
		RequestedQty: req.Quantity,
		AvailableQty: currentStock,
		ProductID:    req.ProductID,
	}

	l.Info("stock checked", "productId", req.ProductID, "available", available, "requested", req.Quantity, "current", currentStock)
	context.JSON(http.StatusOK, response)
}

type ReduceStockRequest struct {
	ProductID int64 `json:"productId" binding:"required" example:"1"`
	Quantity  int   `json:"quantity" binding:"required,min=1" example:"2"`
}

// ReduceStock godoc
// @Summary      Reduce stock quantity
// @Description  Reduce stock when an order is confirmed (used by Order service)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        request  body      ReduceStockRequest  true  "Product and quantity to reduce"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      409      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /products/stock/reduce [post]
func ReduceStock(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())

	var req ReduceStockRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		l.Error("failed to bind request", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid request.", "error": err.Error()})
		return
	}

	l.Debug("ReduceStock called", "productId", req.ProductID, "quantity", req.Quantity)

	err := models.ReduceStock(req.ProductID, req.Quantity)
	if err != nil {
		if _, ok := err.(*models.StockError); ok {
			l.Warn("insufficient stock", "productId", req.ProductID, "quantity", req.Quantity)
			context.JSON(http.StatusConflict, gin.H{"message": "insufficient stock or product not found.", "error": err.Error()})
			return
		}
		l.Error("failed to reduce stock", "productId", req.ProductID, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not reduce stock.", "error": err.Error()})
		return
	}

	l.Info("stock reduced", "productId", req.ProductID, "quantity", req.Quantity)
	context.JSON(http.StatusOK, gin.H{"message": "stock reduced successfully", "productId": req.ProductID, "quantity": req.Quantity})
}
