package models

// Role is a simple lookup for user roles. We only store the name since
// roles are static (admin/user).
type Role struct {
	Name string `json:"name" gorm:"primaryKey"`
}
