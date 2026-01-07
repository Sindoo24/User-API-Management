package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"BACKEND/db/sqlc/generated"
	"BACKEND/internal/repository"
)


type AuthServiceInterface interface {
	ValidatePasswordStrength(password string) error
	CreateUser(ctx context.Context, name, email, password, dobStr, role string) (generated.CreateUserRow, error)
	Login(ctx context.Context, email, password string) (generated.User, string, error)
	GetJWTExpiry() time.Duration
	SetJWTConfig(secret string, expiry time.Duration)
}


type AuthService struct {
	repo       *repository.UserRepository
	jwtSecret  string
	jwtExpiry  time.Duration
}


func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}


func (s *AuthService) SetJWTConfig(secret string, expiry time.Duration) {
	s.jwtSecret = secret
	s.jwtExpiry = expiry
}


func (s *AuthService) GetJWTExpiry() time.Duration {
	return s.jwtExpiry
}


var (
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters long")
	ErrPasswordNoUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoDigit     = errors.New("password must contain at least one digit")
	ErrPasswordNoSpecial   = errors.New("password must contain at least one special character")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
)


type JWTClaims struct {
	UserID int32  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}


func (s *AuthService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return ErrPasswordNoUppercase
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return ErrPasswordNoLowercase
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return ErrPasswordNoDigit
	}

	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?~` + "`" + `]`).MatchString(password)
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}


func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}


func (s *AuthService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *AuthService) CreateUser(ctx context.Context, name, email, password, dobStr, role string) (generated.CreateUserRow, error) {
	if err := s.ValidatePasswordStrength(password); err != nil {
		return generated.CreateUserRow{}, err
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return generated.CreateUserRow{}, err
	}

	
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return generated.CreateUserRow{}, fmt.Errorf("invalid date format: %w", err)
	}

	if role == "" {
		role = "user"
	}

	user, err := s.repo.CreateWithAuth(ctx, name, email, hashedPassword, role, dob)
	if err != nil {
		
		if err.Error() == "duplicate key value violates unique constraint" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return generated.CreateUserRow{}, ErrEmailAlreadyExists
		}
		return generated.CreateUserRow{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) GenerateJWT(userID int32, role string) (string, error) {
	if s.jwtSecret == "" {
		return "", fmt.Errorf("JWT secret not configured")
	}

	expiryTime := time.Now().Add(s.jwtExpiry)
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}


func (s *AuthService) Login(ctx context.Context, email, password string) (generated.User, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return generated.User{}, "", ErrInvalidCredentials
	}

	if err := s.ComparePassword(user.PasswordHash, password); err != nil {
		return generated.User{}, "", ErrInvalidCredentials
	}

	token, err := s.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return generated.User{}, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}
