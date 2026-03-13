package server

import (
	"casino-backend/internal/handlers"
	"casino-backend/internal/middleware"
	"casino-backend/internal/services"
	"casino-backend/internal/websocket"
	"net/http"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Handlers contains all the handlers for the application
type Handlers struct {
	Auth *handlers.AuthHandler
}

func SetupRouter(
	authService *services.AuthService,
	userService *services.UserService,
) *gin.Engine {
	router := gin.Default()

	// Add logging middleware
	router.Use(middleware.LoggingMiddleware())

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://dev-exchange.she.io", "https://www.she.io", "https://she.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoints
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Service is running"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Service is healthy"})
	})

	// create handler instances
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	// Protected routes (accessible to any authenticated user)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authService, false))
	{
		protected.GET("/profile", userHandler.Profile)
	}

	// Admin-only routes
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(authService, true))
	{
		admin.GET("/users", userHandler.GetAllUsers)
	}

	// WebSocket endpoint
	router.GET("/ws", websocket.WebsocketHandler(authService, userService))

	return router
}
