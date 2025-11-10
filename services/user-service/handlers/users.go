package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/user-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary      Get all users
// @Description  Retrieve a list of all users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]interface{}
// @Router       /users [get]
func GetUsers(context *gin.Context) {

	l := logger.FromContext(context.Request.Context())
	l.Debug("GetUsers called")

	users, err := models.GetUsers()

	if err != nil {
		l.Error("failed to fetch events", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch events.", "error": err.Error()})
		return
	}

	l.Info("fetched users", "count", len(users))
	//Response in JSON
	context.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get single user by ID
// @Description  Get details of a user by its ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id} [get]
func GetUser(context *gin.Context) {

	var user *models.User
	userId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	l := logger.FromContext(context.Request.Context())
	l.Debug("GetUser called")

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id.", "error": err.Error()})
		return
	}

	user, err = models.GetUserById(userId)
	if err != nil {
		l.Error("could not fetch event", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event.", "error": err.Error()})
		return
	}

	l.Info("fetched user", "user", user)
	//Response in JSON
	context.JSON(http.StatusOK, user)
}

// GetMyProfile godoc
// @Summary      Get own profile
// @Description  Get profile information of the authenticated user
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /users/me [get]
func GetMyProfile(context *gin.Context) {
	userId := context.GetInt64("userId")
	l := logger.FromContext(context.Request.Context())
	l.Debug("GetMyProfile called", "user_id", userId)

	user, err := models.GetUserById(userId)
	if err != nil {
		l.Error("failed to fetch user profile", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch profile.", "error": err.Error()})
		return
	}

	// Don't send password hash to client
	user.Password = ""

	l.Info("fetched user profile", "user_id", userId)
	context.JSON(http.StatusOK, user)
}

// UpdateMyProfile godoc
// @Summary      Update own profile
// @Description  Update profile information of the authenticated user (firstName, lastName, phone)
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        profile  body      object  true  "Profile payload - Example: {\"firstName\": \"Max\", \"lastName\": \"Mustermann\", \"phone\": \"+49 123 456789\"}"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /users/me [put]
func UpdateMyProfile(context *gin.Context) {
	userId := context.GetInt64("userId")
	l := logger.FromContext(context.Request.Context())
	l.Debug("UpdateMyProfile called", "user_id", userId)

	var profileUpdate struct {
		FirstName *string `json:"firstName"`
		LastName  *string `json:"lastName"`
		Phone     *string `json:"phone"`
	}

	err := context.ShouldBindJSON(&profileUpdate)
	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing user
	user, err := models.GetUserById(userId)
	if err != nil {
		l.Error("failed to fetch user", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch user.", "error": err.Error()})
		return
	}

	// Update fields
	user.FirstName = profileUpdate.FirstName
	user.LastName = profileUpdate.LastName
	user.Phone = profileUpdate.Phone

	err = user.UpdateProfile()
	if err != nil {
		l.Error("failed to update profile", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update profile.", "error": err.Error()})
		return
	}

	// Don't send password hash to client
	user.Password = ""

	l.Info("updated user profile", "user_id", userId)
	context.JSON(http.StatusOK, gin.H{"message": "profile updated successfully", "user": user})
}
