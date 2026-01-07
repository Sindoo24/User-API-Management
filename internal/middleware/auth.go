package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"BACKEND/internal/models"
	"BACKEND/internal/service"
)

const (
	// AuthUserKey is the key used to store authenticated user in fiber context
	AuthUserKey = "authUser"
)

// GetAuthUser extracts the authenticated user from fiber context
// Returns nil if user is not authenticated
func GetAuthUser(c *fiber.Ctx) *models.AuthUser {
	user, ok := c.Locals(AuthUserKey).(models.AuthUser)
	if !ok {
		return nil
	}
	return &user
}

// Auth creates a middleware that validates JWT tokens from Authorization header
func Auth(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			if logger != nil {
				logger.Warn("missing authorization header", zap.String("path", c.Path()))
			}
			return models.SendUnauthorized(c, "Missing authorization header", GetRequestID(c))
		}

		// Check if header starts with "Bearer " (case-insensitive)
		authHeaderLower := strings.ToLower(authHeader)
		if !strings.HasPrefix(authHeaderLower, "bearer ") {
			if logger != nil {
				logger.Warn("invalid authorization header format", zap.String("path", c.Path()))
			}
			return models.SendUnauthorized(c, "Invalid authorization header format. Expected: Bearer <token>", GetRequestID(c))
		}

		// Extract token - find space after "bearer" (case-insensitive)
		// Split by space and take everything after the first word
		parts := strings.Fields(authHeader)
		if len(parts) < 2 {
			if logger != nil {
				logger.Warn("empty token", zap.String("path", c.Path()))
			}
			return models.SendUnauthorized(c, "Token is required", GetRequestID(c))
		}

		// Check if first part is "bearer" (case-insensitive)
		if !strings.EqualFold(parts[0], "bearer") {
			if logger != nil {
				logger.Warn("invalid authorization header format", zap.String("path", c.Path()))
			}
			return models.SendUnauthorized(c, "Invalid authorization header format. Expected: Bearer <token>", GetRequestID(c))
		}

		// Join remaining parts in case token contains spaces (though JWT tokens shouldn't)
		tokenString := strings.Join(parts[1:], " ")
		if tokenString == "" {
			if logger != nil {
				logger.Warn("empty token", zap.String("path", c.Path()))
			}
			return models.SendUnauthorized(c, "Token is required", GetRequestID(c))
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &service.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			if logger != nil {
				logger.Warn("token validation failed", zap.Error(err), zap.String("path", c.Path()))
			}
			return models.SendError(c, fiber.StatusUnauthorized, "Invalid or expired token", models.ErrCodeInvalidToken, GetRequestID(c))
		}

		// Extract claims
		claims, ok := token.Claims.(*service.JWTClaims)
		if !ok || !token.Valid {
			if logger != nil {
				logger.Warn("invalid token claims", zap.String("path", c.Path()))
			}
			return models.SendError(c, fiber.StatusUnauthorized, "Invalid token claims", models.ErrCodeInvalidToken, GetRequestID(c))
		}

		// Create AuthUser from claims
		authUser := models.AuthUser{
			ID:   claims.UserID,
			Role: claims.Role,
		}

		// Inject user into context
		c.Locals(AuthUserKey, authUser)

		if logger != nil {
			logger.Info("user authenticated",
				zap.Int32("user_id", authUser.ID),
				zap.String("role", authUser.Role),
				zap.String("path", c.Path()),
			)
		}

		return c.Next()
	}
}

// RequireRole creates a middleware that checks if the authenticated user has one of the required roles
// This middleware must be used AFTER the Auth middleware as it depends on the authenticated user in context
// Returns 403 Forbidden if user doesn't have the required role
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authenticated user from context
		authUser := GetAuthUser(c)
		if authUser == nil {
			if logger != nil {
				logger.Warn("role check failed: no authenticated user in context",
					zap.String("path", c.Path()),
				)
			}
			return models.SendUnauthorized(c, "Unauthorized", GetRequestID(c))
		}

		// Check if user's role is in the allowed roles
		hasRole := false
		for _, role := range allowedRoles {
			if authUser.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			if logger != nil {
				logger.Warn("role check failed: insufficient permissions",
					zap.Int32("user_id", authUser.ID),
					zap.String("user_role", authUser.Role),
					zap.Strings("required_roles", allowedRoles),
					zap.String("path", c.Path()),
				)
			}
			return models.SendError(c, fiber.StatusForbidden, "Forbidden: insufficient permissions", models.ErrCodeInsufficientPerms, GetRequestID(c))
		}

		if logger != nil {
			logger.Info("role check passed",
				zap.Int32("user_id", authUser.ID),
				zap.String("role", authUser.Role),
				zap.String("path", c.Path()),
			)
		}

		return c.Next()
	}
}
