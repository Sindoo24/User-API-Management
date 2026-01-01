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

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		requestID, _ := c.Locals("requestID").(string)

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
