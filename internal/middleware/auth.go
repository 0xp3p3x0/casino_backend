package middleware

import (
	"casino-backend/internal/logger"
	"casino-backend/internal/models"
	"casino-backend/internal/services"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.AuthServiceInterface, isAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("AuthMiddleware: No Authorization header found")
			logger.Error("AuthMiddleware: No Authorization header found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("AuthMiddleware: Invalid Authorization header format")
			logger.Error("AuthMiddleware: Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in the format: Bearer <token>"})
			return
		}

		token := parts[1]
		log.Printf("AuthMiddleware: Token extracted from header")
		logger.Info("AuthMiddleware: Token extracted from header")

		// Get user from token
		user, err := authService.GetUserFromToken(token, isAdmin)
		if err != nil {
			log.Printf("AuthMiddleware: Failed to get user from token: %v", err)
			logger.Error("AuthMiddleware: Failed to get user from token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		log.Printf("AuthMiddleware: User successfully extracted from token")
		logger.Info("AuthMiddleware: User successfully extracted from token")
		// Store user in context
		c.Set("user", user)
		c.Next()
	}
}

func GetUserFromContext(ctx context.Context) *models.User {
	// Try to get user from Gin context first
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if userVal, exists := ginCtx.Get("user"); exists {
			if user, ok := userVal.(*models.User); ok {
				return user
			}
		}
	}

	// Fallback to standard context
	userVal := ctx.Value("user")
	if userVal == nil {
		return nil
	}

	user, ok := userVal.(*models.User)
	if !ok {
		return nil
	}

	return user
}
