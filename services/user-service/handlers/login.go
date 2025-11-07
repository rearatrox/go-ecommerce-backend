package handlers

import (
	"net/http"
	"rearatrox/event-booking-api/pkg/logger"
	"rearatrox/event-booking-api/services/user-service/models"
	"rearatrox/event-booking-api/services/user-service/utils"

	"github.com/gin-gonic/gin"
)

func Login(context *gin.Context) {
	var user models.User
	l := logger.FromContext(context.Request.Context())

	l.Debug("Login called")
	err := context.ShouldBindJSON(&user)
	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data.", "error": err.Error()})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		l.Error("login failed", "error", err)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "login failed", "error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		l.Error("login failed", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token", "error": err.Error()})
		return
	}

	l.Info("Login successful", "token", token)
	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
