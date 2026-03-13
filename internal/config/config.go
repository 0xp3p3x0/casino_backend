package config

import (
	"casino-backend/internal/models"
	"os"
	"sync"

	"gorm.io/gorm"
)

var (
	cfg  *models.Config
	once sync.Once
	db   *gorm.DB
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDB() *gorm.DB {
	return db
}
func SetDB(database *gorm.DB) {
	db = database
}

func LoadConfig() *models.Config {
	once.Do(func() {
		cfg = &models.Config{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", ""),
			DBName:     getEnv("DB_NAME", "casino_db"),
			Port:       getEnv("PORT", "8080"),
			JWTSecret:  getEnv("JWT_SECRET", "supersecret"),
		}
	})
	return cfg
}
