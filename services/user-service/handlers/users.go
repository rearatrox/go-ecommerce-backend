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
