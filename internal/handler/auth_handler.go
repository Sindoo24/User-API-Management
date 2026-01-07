package handler

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"BACKEND/internal/middleware"
	"BACKEND/internal/models"
	"BACKEND/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService  service.AuthServiceInterface
	validate     *validator.Validate
	logger       *zap.Logger
	cookieSecure bool
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService service.AuthServiceInterface, logger *zap.Logger, cookieSecure bool) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		validate:     validator.New(),
		logger:       logger,
		cookieSecure: cookieSecure,
	}
}

// Signup handles user registration
// POST /auth/signup
func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req models.SignupRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		middleware.GetRequestLogger(c).Error("failed to parse signup request", zap.Error(err))
		return models.SendBadRequest(c, "Invalid request body", middleware.GetRequestID(c))
	}

	// Validate request using go-playground validator
	if err := h.validate.Struct(req); err != nil {
		middleware.GetRequestLogger(c).Error("signup validation failed", zap.Error(err))
		return models.SendError(c, fiber.StatusBadRequest, err.Error(), models.ErrCodeValidationFailed, middleware.GetRequestID(c))
	}

	// Validate password strength
	if err := h.authService.ValidatePasswordStrength(req.Password); err != nil {
		middleware.GetRequestLogger(c).Warn("weak password attempt", zap.String("email", req.Email), zap.Error(err))
		return models.SendError(c, fiber.StatusBadRequest, err.Error(), models.ErrCodeValidationFailed, middleware.GetRequestID(c))
	}

	// Create user (password will be hashed inside the service)
	user, err := h.authService.CreateUser(
		c.Context(),
		req.Name,
		req.Email,
		req.Password,
		req.Dob,
		"user", // default role
	)
	if err != nil {
		// Handle duplicate email error
		if err == service.ErrEmailAlreadyExists {
			middleware.GetRequestLogger(c).Warn("signup attempt with existing email", zap.String("email", req.Email))
			return models.SendConflict(c, "Email already exists", middleware.GetRequestID(c))
		}

		// Handle duplicate email from database error message
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			middleware.GetRequestLogger(c).Warn("signup attempt with existing email (db error)", zap.String("email", req.Email))
			return models.SendConflict(c, "Email already exists", middleware.GetRequestID(c))
		}

		middleware.GetRequestLogger(c).Error("failed to create user", zap.Error(err))
		return models.SendInternalError(c, "Failed to create user", middleware.GetRequestID(c))
	}

	middleware.GetRequestLogger(c).Info("user signed up successfully",
		zap.Int32("user_id", user.ID),
		zap.String("email", user.Email),
	)

	// Return success response (no tokens as per requirements)
	return c.Status(fiber.StatusCreated).JSON(models.SignupResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// Login handles user authentication
// POST /auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		middleware.GetRequestLogger(c).Error("failed to parse login request", zap.Error(err))
		return models.SendBadRequest(c, "Invalid request body", middleware.GetRequestID(c))
	}

	// Validate request using go-playground validator
	if err := h.validate.Struct(req); err != nil {
		middleware.GetRequestLogger(c).Error("login validation failed", zap.Error(err))
		return models.SendError(c, fiber.StatusBadRequest, err.Error(), models.ErrCodeValidationFailed, middleware.GetRequestID(c))
	}

	// Authenticate user and generate token
	user, token, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			middleware.GetRequestLogger(c).Warn("invalid login attempt", zap.String("email", req.Email))
			return models.SendError(c, fiber.StatusUnauthorized, "Invalid email or password", models.ErrCodeInvalidCredentials, middleware.GetRequestID(c))
		}
		middleware.GetRequestLogger(c).Error("failed to login", zap.Error(err))
		return models.SendInternalError(c, "Failed to authenticate user", middleware.GetRequestID(c))
	}

	// Set JWT token in http-only secure cookie with SameSite=Strict
	cookie := &fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.authService.GetJWTExpiry().Seconds()),
		HTTPOnly: true,
		Secure:   h.cookieSecure,
		SameSite: "Strict",
	}
	c.Cookie(cookie)

	middleware.GetRequestLogger(c).Info("user logged in successfully",
		zap.Int32("user_id", user.ID),
		zap.String("email", user.Email),
	)

	// Return success response
	return c.Status(fiber.StatusOK).JSON(models.LoginResponse{
		Message: "Login successful",
		User: struct {
			ID    int32  `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Role  string `json:"role"`
		}{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}
