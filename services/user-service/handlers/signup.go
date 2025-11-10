package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/user-service/models"

	"github.com/gin-gonic/gin"
)

// User-Handlers, die beim Aufruf von Routen /users aufgerufen werden (Verarbeitung der Requests)
// Signup godoc
// @Summary      Create a new user
// @Description  Register a new user account
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User payload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users/signup [post]
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
