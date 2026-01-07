package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger(l *zap.Logger) {
	logger = l
}

// GetRequestID extracts the request ID from fiber context
// Returns empty string if request ID is not found
func GetRequestID(c *fiber.Ctx) string {
	requestID, _ := c.Locals("requestID").(string)
	return requestID
}

// GetRequestLogger returns a logger with request_id field pre-populated
// This ensures all logs within a request context include the request ID
func GetRequestLogger(c *fiber.Ctx) *zap.Logger {
	if logger == nil {
		return zap.NewNop()
	}

	requestID := GetRequestID(c)
	if requestID != "" {
		return logger.With(zap.String("request_id", requestID))
	}

	return logger
}

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestID := GetRequestID(c)

		err := c.Next()

		duration := time.Since(start)

		if logger != nil {
			logger.Info("request completed",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Int("status", c.Response().StatusCode()),
				zap.Duration("duration", duration),
				zap.String("request_id", requestID),
			)
		}

		return err
	}
}
