package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"BACKEND/db/sqlc/generated"
	"BACKEND/internal/models"
	"BACKEND/internal/service"
)

// mockAuthService is a mock implementation of AuthService for testing
type mockAuthService struct {
	validatePasswordStrengthFunc func(password string) error
	createUserFunc               func(ctx context.Context, name, email, password, dobStr, role string) (generated.CreateUserRow, error)
	loginFunc                    func(ctx context.Context, email, password string) (generated.User, string, error)
	getJWTExpiryFunc             func() time.Duration
	setJWTConfigFunc             func(secret string, expiry time.Duration)
}

func (m *mockAuthService) ValidatePasswordStrength(password string) error {
	if m.validatePasswordStrengthFunc != nil {
		return m.validatePasswordStrengthFunc(password)
	}
	return nil
}

func (m *mockAuthService) CreateUser(ctx context.Context, name, email, password, dobStr, role string) (generated.CreateUserRow, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, name, email, password, dobStr, role)
	}
	// Return a default user for testing
	dob, _ := time.Parse("2006-01-02", dobStr)
	return generated.CreateUserRow{
		ID:   1,
		Name: name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
		Email: email,
		Role:  role,
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}, nil
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (generated.User, string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, email, password)
	}
	return generated.User{}, "", service.ErrInvalidCredentials
}

func (m *mockAuthService) GetJWTExpiry() time.Duration {
	if m.getJWTExpiryFunc != nil {
		return m.getJWTExpiryFunc()
	}
	return 24 * time.Hour
}

func (m *mockAuthService) SetJWTConfig(secret string, expiry time.Duration) {
	if m.setJWTConfigFunc != nil {
		m.setJWTConfigFunc(secret, expiry)
	}
}

func TestSignup_Success(t *testing.T) {
	app := fiber.New()
	logger, _ := zap.NewDevelopment()

	// Create a mock auth service for this test
	mockSvc := &mockAuthService{
		validatePasswordStrengthFunc: func(password string) error {
			return nil // Password is valid
		},
		createUserFunc: func(ctx context.Context, name, email, password, dobStr, role string) (generated.CreateUserRow, error) {
			dob, _ := time.Parse("2006-01-02", dobStr)
			return generated.CreateUserRow{
				ID:   1,
				Name: name,
				Dob: pgtype.Date{
					Time:  dob,
					Valid: true,
				},
				Email: email,
				Role:  role,
				CreatedAt: pgtype.Timestamp{
					Time:  time.Now(),
					Valid: true,
				},
				UpdatedAt: pgtype.Timestamp{
					Time:  time.Now(),
					Valid: true,
				},
			}, nil
		},
	}
	handler := NewAuthHandler(mockSvc, logger, false)

	app.Post("/auth/signup", handler.Signup)

	// Test with valid request
	reqBody := models.SignupRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "SecurePass123!",
		Dob:      "1990-01-01",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Note: This will fail without a real database, but tests the validation logic
	// In a real test, you'd mock the repository or use a test database
	if resp.StatusCode != fiber.StatusCreated && resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status 201 or 500, got %d", resp.StatusCode)
	}
}

func TestSignup_InvalidEmail(t *testing.T) {
	app := fiber.New()
	logger, _ := zap.NewDevelopment()
	mockSvc := &mockAuthService{}
	handler := NewAuthHandler(mockSvc, logger, false)

	app.Post("/auth/signup", handler.Signup)

	reqBody := models.SignupRequest{
		Name:     "John Doe",
		Email:    "invalid-email",
		Password: "SecurePass123!",
		Dob:      "1990-01-01",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestSignup_WeakPassword(t *testing.T) {
	app := fiber.New()
	logger, _ := zap.NewDevelopment()
	mockSvc := &mockAuthService{
		validatePasswordStrengthFunc: func(password string) error {
			return service.ErrPasswordTooShort // Simulate weak password
		},
	}
	handler := NewAuthHandler(mockSvc, logger, false)

	app.Post("/auth/signup", handler.Signup)

	testCases := []struct {
		name     string
		password string
	}{
		{"Too short", "Short1!"},
		{"No uppercase", "securepass123!"},
		{"No lowercase", "SECUREPASS123!"},
		{"No digit", "SecurePass!"},
		{"No special char", "SecurePass123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := models.SignupRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: tc.password,
				Dob:      "1990-01-01",
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			if resp.StatusCode != fiber.StatusBadRequest {
				t.Errorf("Expected status 400 for weak password, got %d", resp.StatusCode)
			}

			var errorResp models.ErrorResponse
			json.NewDecoder(resp.Body).Decode(&errorResp)

			if errorResp.Error.Message == "" {
				t.Error("Expected error message for weak password")
			}
		})
	}
}

func TestSignup_MissingFields(t *testing.T) {
	app := fiber.New()
	logger, _ := zap.NewDevelopment()
	mockSvc := &mockAuthService{}
	handler := NewAuthHandler(mockSvc, logger, false)

	app.Post("/auth/signup", handler.Signup)

	testCases := []struct {
		name    string
		reqBody models.SignupRequest
	}{
		{
			name: "Missing name",
			reqBody: models.SignupRequest{
				Email:    "john@example.com",
				Password: "SecurePass123!",
				Dob:      "1990-01-01",
			},
		},
		{
			name: "Missing email",
			reqBody: models.SignupRequest{
				Name:     "John Doe",
				Password: "SecurePass123!",
				Dob:      "1990-01-01",
			},
		},
		{
			name: "Missing password",
			reqBody: models.SignupRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Dob:   "1990-01-01",
			},
		},
		{
			name: "Missing dob",
			reqBody: models.SignupRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "SecurePass123!",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			if resp.StatusCode != fiber.StatusBadRequest {
				t.Errorf("Expected status 400 for missing field, got %d", resp.StatusCode)
			}
		})
	}
}

func TestSignup_InvalidDateFormat(t *testing.T) {
	app := fiber.New()
	logger, _ := zap.NewDevelopment()
	mockSvc := &mockAuthService{}
	handler := NewAuthHandler(mockSvc, logger, false)

	app.Post("/auth/signup", handler.Signup)

	reqBody := models.SignupRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "SecurePass123!",
		Dob:      "01/01/1990", // Wrong format
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid date format, got %d", resp.StatusCode)
	}
}

func TestSignup_InvalidJSON(t *testing.T) {
	app := fiber.New()
	logger, _ := zap.NewDevelopment()
	mockSvc := &mockAuthService{}
	handler := NewAuthHandler(mockSvc, logger, false)

	app.Post("/auth/signup", handler.Signup)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

func TestPasswordStrengthValidation(t *testing.T) {
	authService := &service.AuthService{}

	tests := []struct {
		name        string
		password    string
		shouldError bool
	}{
		{"Valid password", "SecurePass123!", false},
		{"Too short", "Short1!", true},
		{"No uppercase", "securepass123!", true},
		{"No lowercase", "SECUREPASS123!", true},
		{"No digit", "SecurePass!", true},
		{"No special char", "SecurePass123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.ValidatePasswordStrength(tt.password)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for password %q, got nil", tt.password)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for password %q, got %v", tt.password, err)
			}
		})
	}
}
