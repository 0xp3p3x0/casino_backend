package handlers

import (
	"casino-backend/internal/middleware"
	"casino-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler contains handlers that operate on users.
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(us *services.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// Profile returns the profile of the authenticated user.
func (h *UserHandler) Profile(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetAllUsers returns a list of all users (admin only).
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
