package service

import (
	"errors"
	"sync"
	"time"
	repository "github.com/mohamed2394/sahla/internal/repositories"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, string, error)
	Logout(token string) error
	RevokeToken(token string) error
	IsTokenRevoked(token string) bool
	RefreshToken(refreshToken string) (string, string, error)
}

type authService struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	refreshSecret string
	revokedTokens map[string]bool
	tokenMutex    sync.RWMutex
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret, refreshSecret string) AuthService {
	return &authService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		refreshSecret: refreshSecret,
		revokedTokens: make(map[string]bool),
	}
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := s.generateAccessToken(user.UniversalId)
	if err != nil {
		return "", "", err
	}

	refreshToken, refreshTokenExpiresAt, err := s.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Store refresh token in the database
	user.RefreshToken = refreshToken
	user.RefreshTokenExpiresAt = refreshTokenExpiresAt
	if err := s.userRepo.Update(user); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Logout(token string) error {
	// Parse the access token to get the user ID
	claims := &jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		return errors.New("invalid token")
	}

	userID := (*claims)["user_id"].(string)
	user, err := s.userRepo.GetByID(uuid.FromStringOrNil(userID))
	if err != nil {
		return errors.New("user not found")
	}

	// Remove the refresh token and zero out its expiration date
	user.RefreshToken = ""
	user.RefreshTokenExpiresAt = time.Time{}
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Optionally, revoke the access token
	s.tokenMutex.Lock()
	defer s.tokenMutex.Unlock()
	s.revokedTokens[token] = true

	return nil
}

func (s *authService) RevokeToken(token string) error {
	s.tokenMutex.Lock()
	defer s.tokenMutex.Unlock()
	s.revokedTokens[token] = true
	return nil
}

func (s *authService) IsTokenRevoked(token string) bool {
	s.tokenMutex.RLock()
	defer s.tokenMutex.RUnlock()
	return s.revokedTokens[token]
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	// Parse the refresh token
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	// Get the user from the database using the user ID in the claims
	userID := (*claims)["user_id"].(uuid.UUID)
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	// Check if the refresh token matches the one in the database and is not expired
	if user.RefreshToken != refreshToken || time.Now().After(user.RefreshTokenExpiresAt) {
		return "", "", errors.New("invalid or expired refresh token")
	}

	// Generate new access and refresh tokens
	newAccessToken, err := s.generateAccessToken(user.UniversalId)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, refreshTokenExpiresAt, err := s.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Update the user's refresh token in the database
	user.RefreshToken = newRefreshToken
	user.RefreshTokenExpiresAt = refreshTokenExpiresAt
	if err := s.userRepo.Update(user); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *authService) generateAccessToken(userID uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(15 * time.Minute).Unix()

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) generateRefreshToken() (string, time.Time, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &jwt.MapClaims{
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}
