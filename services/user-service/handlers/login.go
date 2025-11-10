package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/user-service/models"
	"rearatrox/go-ecommerce-backend/services/user-service/utils"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      Authenticate user
// @Description  Authenticate a user and return a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.User  true  "User credentials (email + password)"
// @Success      200          {object}  map[string]interface{}
// @Failure      400          {object}  map[string]interface{}
// @Failure      401          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Router       /auth/login [post]
func Login(context *gin.Context) {
	var user models.User
	l := logger.FromContext(context.Request.Context())

	l.Debug("Login called")

	//Payload in user Struct speichern
	err := context.ShouldBindJSON(&user)
	if err != nil {
		l.Warn("invalid request payload", "error", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data.", "error": err.Error()})
		return
	}

	// inside Scan(hash, id, role) --> id + role werden an User-Struct übergeben als Pointer
	err = user.ValidateCredentials()

	if err != nil {
		l.Error("login failed", "error", err)
		context.JSON(http.StatusUnauthorized, gin.H{"message": "login failed", "error": err.Error()})
		return
	}

	//ID + Role aus der DB geholt
	token, err := utils.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		l.Error("login failed", "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token", "error": err.Error()})
		return
	}

	//Bearer Token Format
	token = "Bearer " + token

	l.Info("Login successful", "token", token, "userId", user.ID, "userRole", user.Role)
	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
