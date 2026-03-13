package services

import (
	"casino-backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthServiceInterface exposes just enough methods for the
// authentication middleware. The concrete service provides
// additional helpers used by handlers.
type AuthServiceInterface interface {
	GetUserFromToken(token string, isAdmin bool) (*models.User, error)
}

// AuthService holds dependencies needed for authentication
// and user-related operations.
type AuthService struct {
	db  *gorm.DB
	cfg *models.Config
}

// NewAuthService creates an instance backed by a gorm.DB and
// configuration.
func NewAuthService(db *gorm.DB, cfg *models.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// Register creates a new user with a hashed password. The role
// field is sanitized to one of the supported values.
func (s *AuthService) Register(username, password, role string) (*models.User, error) {
	if role != models.RoleAdmin && role != models.RoleUser {
		role = models.RoleUser
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
	}
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Login checks supplied credentials and returns a signed JWT when
// they are valid.
func (s *AuthService) Login(username, password string) (string, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	return GenerateJWT(&user, s.cfg.JWTSecret)
}

// GetUserFromToken validates a token string and optionally checks
// that the embedded role is "admin" when isAdmin is true.
// It returns the corresponding user record from the database.
func (s *AuthService) GetUserFromToken(tokenString string, isAdmin bool) (*models.User, error) {
	claims, err := ValidateJWT(tokenString, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}
	uid, ok := (*claims)["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	if isAdmin {
		if role, _ := (*claims)["role"].(string); role != models.RoleAdmin {
			return nil, errors.New("admin access required")
		}
	}
	var user models.User
	if err := s.db.First(&user, "id = ?", uid).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
