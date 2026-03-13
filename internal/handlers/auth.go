package handlers

import (
	"casino-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler groups authentication related endpoints.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a handler with the given service.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// registerRequest is the expected body for /auth/register.
type registerRequest struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Nickname     string `json:"nickname" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	ReferrerCode string `json:"referrer_code"`
}

// loginRequest is the expected body for /auth/login. Username can be username or email.
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles new user signups.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.authService.Register(services.RegisterInput{
		Username:     req.Username,
		Email:        req.Email,
		Password:     req.Password,
		Nickname:     req.Nickname,
		PhoneNumber:  req.PhoneNumber,
		ReferrerCode: req.ReferrerCode,
		Role:         "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":            user.ID,
			"username":      user.Username,
			"nickname":      user.Nickname,
			"email":         user.Email,
			"phone_number":  user.PhoneNumber,
			"referral_code": user.ReferralCode,
			"level":         user.Level,
			"role":          user.Role,
			"status":        user.Status,
			"created_at":    user.CreatedAt,
		},
	})
}

// Login authenticates a user and returns a JWT.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
