package server

import (
	"casino-backend/internal/config"
	"casino-backend/internal/handlers"
	"casino-backend/internal/models"
	"casino-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	engine   *gin.Engine
	handlers *Handlers
	db       *gorm.DB
}

func NewServer(db *gorm.DB, cfg *models.Config) (*Server, error) {
	config.SetDB(db)

	authService := services.NewAuthService(db, cfg)
	userService := services.NewUserService(db)

	// previous Handlers struct is kept for compatibility but not used
	authHandler := handlers.NewAuthHandler(authService)
	handlers := &Handlers{
		Auth: authHandler,
	}

	server := &Server{
		engine:   gin.Default(),
		handlers: handlers,
		db:       db,
	}

	server.engine = SetupRouter(authService, userService)
	return server, nil
}

func (s *Server) Run(port string) error {
	return s.engine.Run(":" + port)
}
