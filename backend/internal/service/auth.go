// Package service holds the business logic, sitting between the HTTP handlers
// and the repositories. It depends only on repository interfaces.
package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"creditanalysis/internal/repository"
)

// ErrInvalidCredentials is returned when login fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthService authenticates users and issues JWT tokens.
type AuthService struct {
	users     repository.UserRepository
	jwtSecret []byte
	tokenTTL  time.Duration
}

// NewAuthService builds an AuthService.
func NewAuthService(users repository.UserRepository, secret string) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(secret),
		tokenTTL:  24 * time.Hour,
	}
}

// Login validates the credentials and returns a signed JWT on success.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(u.ID, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}
