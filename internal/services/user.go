package services

import (
	"casino-backend/internal/models"

	"gorm.io/gorm"
)

// UserService handles read operations on users.
type UserService struct {
	db *gorm.DB
}

// NewUserService constructs a UserService.
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetByID returns a user by their ID.
func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll returns all users in the database.
func (s *UserService) GetAll() ([]models.User, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
