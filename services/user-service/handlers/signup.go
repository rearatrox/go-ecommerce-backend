package handlers

import (
	"net/http"
	"rearatrox/event-booking-api/pkg/logger"
	"rearatrox/event-booking-api/services/user-service/models"

	"github.com/gin-gonic/gin"
)

// User-Handlers, die beim Aufruf von Routen /users aufgerufen werden (Verarbeitung der Requests)

func Signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	l := logger.FromContext(context.Request.Context())
	l.Debug("Signup called")

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data.", "error": err.Error()})
		return
	}

	err = user.SaveUser()
	if err != nil {
		l.Error("could not create event", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create event.", "error": err.Error()})
		return
	}
	l.Info("signup successful")
	context.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}
