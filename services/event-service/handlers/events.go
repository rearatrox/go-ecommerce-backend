package handlers

import (
	"net/http"
	"rearatrox/event-booking-api/pkg/logger"
	"rearatrox/event-booking-api/services/event-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Event-Handlers, die beim Aufruf von Routen /events aufgerufen werden (Verarbeitung der Requests)

func GetEvents(context *gin.Context) {

	l := logger.FromContext(context.Request.Context())
	l.Debug("GetEvents called")

	events, err := models.GetEvents()

	if err != nil {
		l.Error("failed to fetch events", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch events.", "error": err.Error()})
		return
	}

	l.Info("fetched events", "count", len(events))
	//Response in JSON
	context.JSON(http.StatusOK, events)
}

func GetEvent(context *gin.Context) {
	var event *models.Event
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id.", "error": err.Error()})
		return
	}

	l := logger.FromContext(context.Request.Context())
	l.Debug("GetEvent called", "event_id", eventId)

	event, err = models.GetEventByID(eventId)
	if err != nil {
		l.Error("failed to fetch event", "event_id", eventId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event.", "error": err.Error()})
		return
	}

	l.Info("fetched event", "event_id", eventId)
	//Response in JSON
	context.JSON(http.StatusOK, event)
}

func CreateEvent(context *gin.Context) {

	l := logger.FromContext(context.Request.Context())
	l.Debug("CreateEvent called")

	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := context.GetInt64("userId")
	event.CreatorID = userId

	err = event.SaveEvent()

	if err != nil {
		l.Error("failed to save event", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create event.", "error": err.Error()})
		return
	}
	l.Info("created event", "event_id", event.ID, "creator_id", userId)
	context.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}

func UpdateEvent(context *gin.Context) {
	var event *models.Event
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id to update event.", "error": err.Error()})
		return
	}

	l := logger.FromContext(context.Request.Context())
	l.Debug("UpdateEvent called", "event_id", eventId)

	event, err = models.GetEventByID(eventId)
	if err != nil {
		l.Error("failed to fetch event", "event_id", eventId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event.", "error": err.Error()})
		return
	}

	userId, _ := context.Get("userId")

	if event.CreatorID != userId {
		l.Warn("unauthorized update attempt", "event_id", eventId, "user", userId)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized to update that event"})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		l.Warn("failed to parse update payload", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request data.", "error": err.Error()})
		return
	}

	updatedEvent.ID = event.ID
	err = updatedEvent.UpdateEvent()
	if err != nil {
		l.Error("failed to update event", "event_id", eventId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update event.", "error": err.Error()})
		return
	}
	//Response in JSON
	l.Info("updated event", "event_id", eventId)
	context.JSON(http.StatusOK, gin.H{"message": "updated event successfully", "updatedEvent": updatedEvent})
}

func DeleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event id to update event.", "error": err.Error()})
		return
	}

	l := logger.FromContext(context.Request.Context())
	l.Debug("DeleteEvent called", "event_id", eventId)

	deleteEvent, err := models.GetEventByID(eventId)
	if err != nil {
		l.Error("failed to fetch event", "event_id", eventId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event.", "error": err.Error()})
		return
	}

	userId, _ := context.Get("userId")

	if deleteEvent.CreatorID != userId {
		l.Warn("unauthorized delete attempt", "event_id", eventId, "user", userId)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized to delete that event"})
		return
	}

	err = deleteEvent.DeleteEvent()
	if err != nil {
		l.Error("failed to delete event", "event_id", eventId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete event.", "error": err.Error()})
		return
	}
	//Response in JSON
	l.Info("deleted event", "event_id", eventId)
	context.JSON(http.StatusOK, gin.H{"message": "deleted event successfully", "deletedEvent": deleteEvent})
}
