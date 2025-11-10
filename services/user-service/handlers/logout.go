package handlers

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"rearatrox/go-ecommerce-backend/services/user-service/models"

	"github.com/gin-gonic/gin"
)

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate all existing JWT tokens by incrementing the token version
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /auth/logout [post]
func Logout(context *gin.Context) {
	l := logger.FromContext(context.Request.Context())
	userId := context.GetInt64("userId")

	l.Debug("Logout called", "user_id", userId)

	// Get user to increment token version
	user, err := models.GetUserById(userId)
	if err != nil {
		l.Error("failed to get user", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not logout user.", "error": err.Error()})
		return
	}

	// Increment token version to invalidate all existing tokens
	err = user.IncrementTokenVersion()
	if err != nil {
		l.Error("failed to increment token version", "user_id", userId, "error", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not logout user.", "error": err.Error()})
		return
	}

	l.Info("Logout successful", "user_id", userId, "new_token_version", user.TokenVersion)
	context.JSON(http.StatusOK, gin.H{"message": "Logout successful. All tokens have been invalidated."})
}
