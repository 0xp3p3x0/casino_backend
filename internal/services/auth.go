package services

import (
	"casino-backend/internal/models"
	"crypto/rand"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const referralCodeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const referralCodeLen = 8

func generateReferralCode() string {
	b := make([]byte, referralCodeLen)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = referralCodeChars[int(b[i])%len(referralCodeChars)]
	}
	return string(b)
}

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

// RegisterInput holds optional fields for registration.
type RegisterInput struct {
	Username     string
	Email        string
	Password     string
	Nickname     string
	PhoneNumber  string
	ReferrerCode string
	Role         string
}

// Register creates a new user with a hashed password. Email and username are required.
// Role is sanitized to one of the supported values. Assigns Level 1, Status active,
// and a generated referral code.
func (s *AuthService) Register(in RegisterInput) (*models.User, error) {
	role := in.Role
	if role != models.RoleAdmin && role != models.RoleUser {
		role = models.RoleUser
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user := &models.User{
		Username:     in.Username,
		Nickname:     in.Nickname,
		Email:        in.Email,
		PhoneNumber:  in.PhoneNumber,
		Level:        1,
		ReferralCode: generateReferralCode(),
		RefererCode:  in.ReferrerCode,
		PasswordHash: string(hash),
		Role:         role,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Login checks supplied credentials (username or email) and returns a signed JWT when valid.
func (s *AuthService) Login(usernameOrEmail, password string) (string, error) {
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error; err != nil {
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
	sub, ok := (*claims)["sub"]
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	var uid uint
	switch v := sub.(type) {
	case float64:
		uid = uint(v)
	default:
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
