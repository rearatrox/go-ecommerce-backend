package handlers

import (
	"net/http"
	"rearatrox/event-booking-api/pkg/logger"
	"rearatrox/event-booking-api/services/user-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
