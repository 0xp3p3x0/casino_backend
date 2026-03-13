package models

import "time"

// Role constants used throughout the application
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User represents an application user. PasswordHash is stored
// for authentication and is omitted from any JSON responses.
type User struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey"`
	Username     string    `json:"username" gorm:"unique;not null"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role" gorm:"not null;default:user"`
	CreatedAt    time.Time `json:"created_at"`
}
