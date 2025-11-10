package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/product-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCategories godoc
// @Summary      Get all categories
// @Description  Get all category information of all categories
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Category
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /categories [get]
func GetCategories(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	l.Debug("GetCategories called")

	categories, err := models.GetCategories()

	if err != nil {
		l.Error("failed to fetch categories", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch categories.", "error": err.Error()})
		return
	}

	l.Info("fetched categories", "count", len(categories))
	context.JSON(http.StatusOK, categories)
}

// GetCategoryByID godoc
// @Summary      Get single category by ID
// @Description  Get details of a category by its ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /categories/id/{id} [get]
func GetCategoryByID(context *gin.Context) {
	categoryId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	l := logger.FromContext(context.Request.Context())

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse category id.", "error": err.Error()})
		return
	}

	l.Debug("GetCategoryByID called", "category_id", categoryId)

	category, err := models.GetCategoryByID(categoryId)
	if err != nil {
		l.Error("failed to fetch category", "category_id", categoryId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch category.", "error": err.Error()})
		return
	}

	l.Info("fetched category", "category_id", categoryId)
	context.JSON(http.StatusOK, category)
}

// GetCategoryBySlug godoc
// @Summary      Get single category by slug
// @Description  Get details of a category by its slug
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        slug   path      string  true  "Category Slug"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /categories/slug/{slug} [get]
func GetCategoryBySlug(context *gin.Context) {
	categorySlug := context.Param("slug")

	l := logger.FromContext(context.Request.Context())
	l.Debug("GetCategoryBySlug called", "slug", categorySlug)

	category, err := models.GetCategoryBySlug(categorySlug)
	if err != nil {
		l.Error("failed to fetch category", "slug", categorySlug, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch category.", "error": err.Error()})
		return
	}

	l.Info("fetched category", "slug", categorySlug)
	context.JSON(http.StatusOK, category)
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create category. Requires authentication and authorization as role "admin".
// @Tags         Categories (Admin)
// @Accept       json
// @Produce      json
// @Param        category  body      models.Category  true  "Category payload - Example:"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/categories/create [post]
func CreateCategory(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	l.Debug("CreateCategory called")

	var category models.Category
	err := context.ShouldBindJSON(&category)

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = category.InsertCategory()

	if err != nil {
		l.Error("failed to save category", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create category.", "error": err.Error()})
		return
	}
	l.Info("created category", "category_id", category.ID, "slug", category.Slug)
	context.JSON(http.StatusCreated, gin.H{"message": "Category created", "category": category})
}

// UpdateCategory godoc
// @Summary      Update an existing category
// @Description  Update category by slug. Requires authentication and authorization as role "admin". Note: slug in body must match slug in URL path.
// @Tags         Categories (Admin)
// @Accept       json
// @Produce      json
// @Param        slug     path      string           true  "Category Slug"
// @Param        category  body      models.Category  true  "Updated category payload - Example:"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/categories/update/{slug} [put]
func UpdateCategory(context *gin.Context) {
	categorySlug := context.Param("slug")

	l := logger.FromContext(context.Request.Context())
	l.Debug("UpdateCategory called", "categorySlug", categorySlug)

	category, err := models.GetCategoryBySlug(categorySlug)
	if err != nil {
		l.Error("failed to fetch category", "categorySlug", categorySlug, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch category.", "error": err.Error()})
		return
	}

	var updatedCategory models.Category
	err = context.ShouldBindJSON(&updatedCategory)
	if err != nil {
		l.Warn("failed to parse update payload", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request data.", "error": err.Error()})
		return
	}

	if updatedCategory.Slug != categorySlug {
		l.Error("slug values are not equal (either route parameter or request slug is false)", "updatedCategory_slug", updatedCategory.Slug, "routeParam slug", categorySlug)
		context.JSON(http.StatusBadRequest, gin.H{"message": "slug values are not equal (either route parameter or request slug is false)"})
		return
	}

	updatedCategory.ID = category.ID

	err = updatedCategory.UpdateCategory()
	if err != nil {
		l.Error("failed to update category", "categorySlug", categorySlug, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update category.", "error": err.Error()})
		return
	}
	l.Info("updated category", "categorySlug", categorySlug)
	context.JSON(http.StatusOK, gin.H{"message": "updated category successfully", "updatedCategory": updatedCategory})
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Delete category by slug. Requires authentication and authorization as role "admin".
// @Tags         Categories (Admin)
// @Accept       json
// @Produce      json
// @Param        slug   path  string  true  "Category slug"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/categories/delete/{slug} [delete]
func DeleteCategoryBySlug(context *gin.Context) {
	categorySlug := context.Param("slug")
	l := logger.FromContext(context.Request.Context())
	l.Debug("DeleteCategory called", "categorySlug", categorySlug)

	deleteCategory, err := models.GetCategoryBySlug(categorySlug)
	if err != nil {
		l.Error("failed to fetch category", "categorySlug", categorySlug, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch category.", "error": err.Error()})
		return
	}

	err = deleteCategory.DeleteCategory()
	if err != nil {
		l.Error("failed to delete category", "categorySlug", categorySlug, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete category.", "error": err.Error()})
		return
	}
	l.Info("deleted category", "categorySlug", categorySlug)
	context.JSON(http.StatusOK, gin.H{"message": "deleted category successfully", "deletedCategory": deleteCategory})
}

// GetProductsByCategory godoc
// @Summary      Get products by category
// @Description  Get all products assigned to a specific category
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        slug   path      string  true  "Category Slug"
// @Success      200  {array}   models.Product
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /categories/{slug}/products [get]
func GetProductsByCategory(context *gin.Context) {
	categorySlug := context.Param("slug")
	l := logger.FromContext(context.Request.Context())
	l.Debug("GetProductsByCategory called", "categorySlug", categorySlug)

	category, err := models.GetCategoryBySlug(categorySlug)
	if err != nil {
		l.Error("failed to fetch category", "categorySlug", categorySlug, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch category.", "error": err.Error()})
		return
	}

	products, err := models.GetProductsByCategory(category.ID)
	if err != nil {
		l.Error("failed to fetch products by category", "categorySlug", categorySlug, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch products by category.", "error": err.Error()})
		return
	}

	l.Info("fetched products by category", "categorySlug", categorySlug, "count", len(products))
	context.JSON(http.StatusOK, products)
}
