package models

import "github.com/gofiber/fiber/v2"

// Error codes for standardized error responses
const (
	// Authentication errors
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeMissingAuth        = "MISSING_AUTH_HEADER"
	ErrCodeInvalidToken       = "INVALID_TOKEN"
	ErrCodeExpiredToken       = "EXPIRED_TOKEN"

	// Authorization errors
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeInsufficientPerms = "INSUFFICIENT_PERMISSIONS"

	// Validation errors
	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeInvalidFormat    = "INVALID_FORMAT"

	// Resource errors
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"

	// Server errors
	ErrCodeInternalError = "INTERNAL_ERROR"
	ErrCodeDatabaseError = "DATABASE_ERROR"
)

// NewErrorResponse creates a standardized error response
func NewErrorResponse(message, code, requestID string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Message:   message,
			Code:      code,
			RequestID: requestID,
		},
	}
}

// SendError sends a standardized error response with the given status code
func SendError(c *fiber.Ctx, status int, message, code, requestID string) error {
	return c.Status(status).JSON(NewErrorResponse(message, code, requestID))
}

// Common error response helpers
func SendBadRequest(c *fiber.Ctx, message, requestID string) error {
	return SendError(c, fiber.StatusBadRequest, message, ErrCodeInvalidInput, requestID)
}

func SendUnauthorized(c *fiber.Ctx, message, requestID string) error {
	return SendError(c, fiber.StatusUnauthorized, message, ErrCodeUnauthorized, requestID)
}

func SendForbidden(c *fiber.Ctx, message, requestID string) error {
	return SendError(c, fiber.StatusForbidden, message, ErrCodeForbidden, requestID)
}

func SendNotFound(c *fiber.Ctx, message, requestID string) error {
	return SendError(c, fiber.StatusNotFound, message, ErrCodeNotFound, requestID)
}

func SendConflict(c *fiber.Ctx, message, requestID string) error {
	return SendError(c, fiber.StatusConflict, message, ErrCodeAlreadyExists, requestID)
}

func SendInternalError(c *fiber.Ctx, message, requestID string) error {
	return SendError(c, fiber.StatusInternalServerError, message, ErrCodeInternalError, requestID)
}
