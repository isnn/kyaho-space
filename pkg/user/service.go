package user

import (
	"errors"
	"kyaho-space/pkg/entities"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(req entities.SignupRequest) (*entities.User, error)
	Login(req entities.LoginRequest, secret string) (accessToken string, refreshToken string, err error)
	RefreshToken(refreshToken string, secret string) (newAccessToken string, err error)
	Logout(refreshToken string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

// Register creates a new user after validating uniqueness and hashing the password.
func (s *service) Register(req entities.SignupRequest) (*entities.User, error) {
	// Check duplicate username
	if existing, _ := s.repository.FindByUsername(req.Username); existing != nil {
		return nil, errors.New("username already exists")
	}

	// Check duplicate email
	if existing, _ := s.repository.FindByEmail(req.Email); existing != nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &entities.User{
		Username: req.Username,
		Email:    strings.ToLower(req.Email),
		Password: string(hashed),
	}

	if err := s.repository.CreateUser(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// Login authenticates a user by username or email and returns dual JWT tokens.
func (s *service) Login(req entities.LoginRequest, secret string) (string, string, error) {
	var user *entities.User
	var err error

	// Try username first, then email
	user, err = s.repository.FindByUsername(req.Identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = s.repository.FindByEmail(req.Identifier)
			if err != nil {
				return "", "", errors.New("invalid credentials")
			}
		} else {
			return "", "", errors.New("invalid credentials")
		}
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Generate access token (15 minutes)
	accessToken, err := generateToken(user.ID, secret, 15*time.Minute)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	// Generate refresh token (7 days)
	refreshToken, err := generateToken(user.ID, secret, 7*24*time.Hour)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	// Persist refresh token in DB
	if err := s.repository.UpdateToken(user.ID, refreshToken); err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates a refresh token against DB and issues a new access token.
func (s *service) RefreshToken(refreshToken string, secret string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired refresh token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return "", errors.New("invalid refresh token")
	}

	user, err := s.repository.FindByToken(refreshToken)
	if err != nil || user.ID != userID {
		return "", errors.New("refresh token not recognized")
	}
	newAccess, err := generateToken(userID, secret, 7*24*time.Hour)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return newAccess, nil
}

func (s *service) Logout(refreshToken string) error {
	user, err := s.repository.FindByToken(refreshToken)
	if err != nil {
		return nil
	}

	return s.repository.UpdateToken(user.ID, "")
}

func generateToken(userID string, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
