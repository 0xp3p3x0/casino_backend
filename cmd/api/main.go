package main

import (
	"casino-backend/internal/config"
	"casino-backend/internal/db"
	"casino-backend/internal/logger"
	"casino-backend/internal/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		// Don't fatal here, as .env might not exist in production
	}

	// Initialize logger
	if err := logger.Init(); err != nil {
		logger.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database := db.Init(cfg)

	// Create and start server
	srv, err := server.NewServer(database, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
		logger.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		logger.Fatalf("Failed to start server: %v", err)
	}
}
