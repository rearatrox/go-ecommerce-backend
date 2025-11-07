package handlers

import (
	"net/http"
	"rearatrox/event-booking-api/pkg/logger"
	"rearatrox/event-booking-api/services/event-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddRegistrationForEvent(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	l.Debug("AddRegistrationForEvent called")

	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		l.Warn("invalid event id", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id.", "error": err.Error()})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		l.Error("could not fetch event", "event_id", eventId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event."})
		return
	}

	err = event.Register(userId)
	if err != nil {
		l.Error("failed to register user", "event_id", eventId, "user", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not register user for event."})
		return
	}

	l.Info("registered user for event", "event_id", eventId, "user", userId)
	context.JSON(http.StatusCreated, gin.H{"message": "successfully registered user for event."})

}

func DeleteRegistrationForEvent(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	l.Debug("DeleteRegistrationForEvent called")

	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		l.Warn("invalid event id", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id.", "error": err.Error()})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		l.Error("failed to fetch event", "event_id", eventId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event."})
		return
	}

	err = event.DeleteRegistration(userId)
	if err != nil {
		l.Error("failed to delete registration", "event_id", eventId, "user", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete user registration for event."})
		return
	}

	l.Info("deleted registration for event", "event_id", eventId, "user", userId)
	context.JSON(http.StatusCreated, gin.H{"message": "successfully deleted user registration for event."})

}
