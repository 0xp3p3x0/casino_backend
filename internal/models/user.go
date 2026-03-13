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
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	Username string `json:"username" gorm:"size:50;not null"` // duplicate allowed

	Email       string `json:"email" gorm:"size:120;uniqueIndex;not null"`
	PhoneNumber string `json:"phone_number" gorm:"size:20;uniqueIndex"`
	Nickname    string `json:"nickname" gorm:"size:50;uniqueIndex"`

	Level int `json:"level" gorm:"not null;default:1"`

	ReferralCode string `json:"referral_code" gorm:"size:20;uniqueIndex"`
	RefererCode  string `json:"referrer_code" gorm:"size:20;index"`

	PasswordHash string `json:"-" gorm:"not null"`

	Role   string `json:"role" gorm:"size:20;not null;default:user"`
	Status string `json:"status" gorm:"size:20;default:active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
