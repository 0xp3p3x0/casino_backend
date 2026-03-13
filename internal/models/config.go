package models

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
	JWTSecret  string
}

const (
	ErrUnauthorized     = "Unauthorized user"
	ErrMissingAuthToken = "Authorization token is required"
	ErrInvalidToken     = "Invalid or expired token"
)
