package websocket

import (
	"log"
	"net/http"
	"time"

	"casino-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebsocketHandler upgrades an HTTP connection and pushes profile and
// periodic health messages to the client. Clients must supply a valid
// JWT via the "token" query parameter.
func WebsocketHandler(authService services.AuthServiceInterface, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token query parameter required"})
			return
		}
		user, err := authService.GetUserFromToken(token, false)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("websocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		// send initial profile
		if profile, err := userService.GetByID(user.ID); err == nil {
			conn.WriteJSON(gin.H{"type": "profile", "data": profile})
		}

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				conn.WriteJSON(gin.H{"type": "health", "status": "ok", "time": time.Now()})
			}
		}
	}
}
