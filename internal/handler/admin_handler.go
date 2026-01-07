package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"BACKEND/internal/middleware"
	"BACKEND/internal/models"
	"BACKEND/internal/repository"
)

// AdminHandler handles admin-only operations
type AdminHandler struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(repo *repository.UserRepository, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetAllUsers returns all users (admin only)
// GET /admin/users
func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	authUser := middleware.GetAuthUser(c)

	middleware.GetRequestLogger(c).Info("admin accessing all users",
		zap.Int32("admin_id", authUser.ID),
	)

	users, err := h.repo.List(c.Context())
	if err != nil {
		middleware.GetRequestLogger(c).Error("failed to list all users", zap.Error(err))
		return models.SendInternalError(c, "Failed to retrieve users", middleware.GetRequestID(c))
	}

	return c.JSON(fiber.Map{
		"total": len(users),
		"users": users,
	})
}

// GetStats returns system statistics (admin only)
// GET /admin/stats
func (h *AdminHandler) GetStats(c *fiber.Ctx) error {
	authUser := middleware.GetAuthUser(c)

	middleware.GetRequestLogger(c).Info("admin accessing stats",
		zap.Int32("admin_id", authUser.ID),
	)

	count, err := h.repo.Count(c.Context())
	if err != nil {
		middleware.GetRequestLogger(c).Error("failed to get user count", zap.Error(err))
		return models.SendInternalError(c, "Failed to retrieve statistics", middleware.GetRequestID(c))
	}

	return c.JSON(fiber.Map{
		"total_users": count,
		"message":     "Admin statistics",
	})
}
