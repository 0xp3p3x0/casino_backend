package models

import "time"

type Wallet struct {
	ID     string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uint   `gorm:"not null;uniqueIndex"`

	User User `gorm:"foreignKey:UserID"`

	Balance float64 `gorm:"type:numeric(18,2);default:0"`
	Points  int64   `gorm:"default:0"`

	Currency string `gorm:"size:10;default:'USD'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
