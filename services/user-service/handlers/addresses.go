package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/user-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserAddresses godoc
// @Summary      Get user's addresses
// @Description  Get all addresses of the authenticated user
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Address
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /users/me/addresses [get]
func GetUserAddresses(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("GetUserAddresses called", "user_id", userId)

	addresses, err := models.GetUserAddresses(userId)
	if err != nil {
		l.Error("failed to fetch user addresses", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch addresses.", "error": err.Error()})
		return
	}

	l.Info("fetched user addresses", "user_id", userId, "count", len(addresses))
	context.JSON(http.StatusOK, addresses)
}

// GetAddressByID godoc
// @Summary      Get single address by ID
// @Description  Get details of a specific address of the authenticated user
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Address ID"
// @Success      200  {object}  models.Address
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /users/me/addresses/{id} [get]
func GetAddressByID(context *gin.Context) {
	userId := context.GetInt64("userId")
	addressId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	l := logger.FromContext(context.Request.Context())

	if err != nil {
		l.Warn("invalid address id", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse address id.", "error": err.Error()})
		return
	}

	l.Debug("GetAddressByID called", "user_id", userId, "address_id", addressId)

	address, err := models.GetAddressByID(addressId, userId)
	if err != nil {
		l.Error("failed to fetch address", "user_id", userId, "address_id", addressId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "address not found.", "error": err.Error()})
		return
	}

	l.Info("fetched address", "user_id", userId, "address_id", addressId)
	context.JSON(http.StatusOK, address)
}

// CreateAddress godoc
// @Summary      Create a new address
// @Description  Create a new address for the authenticated user
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        address  body      models.Address  true  "Address payload"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /users/me/addresses [post]
func CreateAddress(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")
	l.Debug("CreateAddress called", "user_id", userId)

	var address models.Address
	err := context.ShouldBindJSON(&address)

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address.UserID = userId

	err = address.InsertAddress()
	if err != nil {
		l.Error("failed to save address", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create address.", "error": err.Error()})
		return
	}

	l.Info("created address", "user_id", userId, "address_id", address.ID)
	context.JSON(http.StatusCreated, gin.H{"message": "Address created", "address": address})
}

// UpdateAddress godoc
// @Summary      Update an address
// @Description  Update an existing address of the authenticated user
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        id       path      int             true  "Address ID"
// @Param        address  body      models.Address  true  "Updated address payload"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /users/me/addresses/{id} [put]
func UpdateAddress(context *gin.Context) {
	userId := context.GetInt64("userId")
	addressId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	l := logger.FromContext(context.Request.Context())

	if err != nil {
		l.Warn("invalid address id", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse address id.", "error": err.Error()})
		return
	}

	l.Debug("UpdateAddress called", "user_id", userId, "address_id", addressId)

	// Check if address exists and belongs to user
	existingAddress, err := models.GetAddressByID(addressId, userId)
	if err != nil {
		l.Error("failed to fetch address", "user_id", userId, "address_id", addressId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "address not found.", "error": err.Error()})
		return
	}

	var updatedAddress models.Address
	err = context.ShouldBindJSON(&updatedAddress)
	if err != nil {
		l.Warn("failed to parse update payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data.", "error": err.Error()})
		return
	}

	updatedAddress.ID = existingAddress.ID
	updatedAddress.UserID = userId

	err = updatedAddress.UpdateAddress()
	if err != nil {
		l.Error("failed to update address", "user_id", userId, "address_id", addressId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update address.", "error": err.Error()})
		return
	}

	l.Info("updated address", "user_id", userId, "address_id", addressId)
	context.JSON(http.StatusOK, gin.H{"message": "updated address successfully", "address": updatedAddress})
}

// DeleteAddress godoc
// @Summary      Delete an address
// @Description  Delete an address of the authenticated user
// @Tags         Addresses
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Address ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /users/me/addresses/{id} [delete]
func DeleteAddress(context *gin.Context) {
	userId := context.GetInt64("userId")
	addressId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	l := logger.FromContext(context.Request.Context())

	if err != nil {
		l.Warn("invalid address id", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse address id.", "error": err.Error()})
		return
	}

	l.Debug("DeleteAddress called", "user_id", userId, "address_id", addressId)

	// Check if address exists and belongs to user
	address, err := models.GetAddressByID(addressId, userId)
	if err != nil {
		l.Error("failed to fetch address", "user_id", userId, "address_id", addressId, "error", err)
		context.JSON(http.StatusNotFound, gin.H{"message": "address not found.", "error": err.Error()})
		return
	}

	err = address.DeleteAddress()
	if err != nil {
		l.Error("failed to delete address", "user_id", userId, "address_id", addressId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete address.", "error": err.Error()})
		return
	}

	l.Info("deleted address", "user_id", userId, "address_id", addressId)
	context.JSON(http.StatusOK, gin.H{"message": "deleted address successfully"})
}
