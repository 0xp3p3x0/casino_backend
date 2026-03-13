package db

import (
	"casino-backend/internal/models"
	"fmt"
	"log"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg *models.Config) *gorm.DB {
	sslMode := "require"
	if cfg.DBHost == "localhost" || cfg.DBHost == "127.0.0.1" {
		sslMode = "disable"
	}

	// Safely encode username and password
	user := url.UserPassword(cfg.DBUser, cfg.DBPassword)

	// Build URL
	dbURL := &url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
		RawQuery: url.Values{
			"sslmode": []string{sslMode},
		}.Encode(),
	}

	dsn := dbURL.String()
	log.Printf("Connecting to DB: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	// Auto-migrate the tables we care about
	if err := db.AutoMigrate(&models.User{}, &models.Wallet{}); err != nil {
		log.Fatalf("failed to automigrate: %v", err)
	}

	return db
}
