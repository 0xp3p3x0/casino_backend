package middleware

import (
	"casino-backend/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs API requests and responses
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log the request
		logger.LogAPIRequest(c.Request, c.Writer.Status(), duration)

		// Log any errors
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.LogAPIError(c.Request, e.Err)
			}
		}
	}
}
