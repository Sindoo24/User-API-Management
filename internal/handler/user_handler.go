package handler

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"BACKEND/internal/middleware"
	"BACKEND/internal/models"
	"BACKEND/internal/repository"
	"BACKEND/internal/service"
)

type UserHandler struct {
	repo     *repository.UserRepository
	service  *service.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

func NewUserHandler(r *repository.UserRepository, s *service.UserService, l *zap.Logger) *UserHandler {
	return &UserHandler{
		repo:     r,
		service:  s,
		validate: validator.New(), // initialize validator once
		logger:   l,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req models.UserRequest

	if err := c.BodyParser(&req); err != nil {
		middleware.GetRequestLogger(c).Error("failed to parse request body", zap.Error(err))
		return models.SendBadRequest(c, "Invalid request body", middleware.GetRequestID(c))
	}

	if err := h.validate.Struct(req); err != nil {
		middleware.GetRequestLogger(c).Error("validation failed", zap.Error(err))
		return models.SendError(c, fiber.StatusBadRequest, err.Error(), models.ErrCodeValidationFailed, middleware.GetRequestID(c))
	}

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		middleware.GetRequestLogger(c).Error("invalid date format", zap.Error(err))
		return models.SendBadRequest(c, "Invalid date format, use YYYY-MM-DD", middleware.GetRequestID(c))
	}

	user, err := h.repo.Create(c.Context(), req.Name, dob)
	if err != nil {
		middleware.GetRequestLogger(c).Error("create user failed", zap.Error(err))
		return models.SendInternalError(c, "Failed to create user", middleware.GetRequestID(c))
	}

	middleware.GetRequestLogger(c).Info("user created", zap.Int32("id", user.ID))

	return c.Status(201).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
	})
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return models.SendBadRequest(c, "Invalid user ID", middleware.GetRequestID(c))
	}

	resp, err := h.service.GetUserWithAge(c.Context(), int32(id))
	if err != nil {
		middleware.GetRequestLogger(c).Error("get user failed", zap.Error(err))
		return models.SendNotFound(c, "User not found", middleware.GetRequestID(c))
	}

	return c.JSON(resp)
}

// GetCurrentUser returns the authenticated user's details
// GET /users/me
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	// Extract authenticated user from context (set by auth middleware)
	authUser := middleware.GetAuthUser(c)
	if authUser == nil {
		middleware.GetRequestLogger(c).Error("auth user not found in context")
		return models.SendUnauthorized(c, "Unauthorized", middleware.GetRequestID(c))
	}

	// Fetch full user details from database using user_id from context
	resp, err := h.service.GetUserWithAge(c.Context(), authUser.ID)
	if err != nil {
		middleware.GetRequestLogger(c).Error("get current user failed", zap.Int32("user_id", authUser.ID), zap.Error(err))
		return models.SendNotFound(c, "User not found", middleware.GetRequestID(c))
	}

	middleware.GetRequestLogger(c).Info("current user retrieved", zap.Int32("user_id", authUser.ID))

	return c.JSON(resp)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr != "" || limitStr != "" {
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}

		paginatedResp, err := h.service.ListUsersWithAgePaginated(c.Context(), page, limit)
		if err != nil {
			middleware.GetRequestLogger(c).Error("list users paginated failed", zap.Error(err))
			return models.SendInternalError(c, "Failed to list users", middleware.GetRequestID(c))
		}

		return c.JSON(paginatedResp)
	}

	users, err := h.service.ListUsersWithAge(c.Context())
	if err != nil {
		middleware.GetRequestLogger(c).Error("list users failed", zap.Error(err))
		return models.SendInternalError(c, "Failed to list users", middleware.GetRequestID(c))
	}

	return c.JSON(users)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return models.SendBadRequest(c, "Invalid user ID", middleware.GetRequestID(c))
	}

	var req models.UserRequest
	if err := c.BodyParser(&req); err != nil {
		middleware.GetRequestLogger(c).Error("failed to parse request body", zap.Error(err))
		return models.SendBadRequest(c, "Invalid request body", middleware.GetRequestID(c))
	}

	if err := h.validate.Struct(req); err != nil {
		middleware.GetRequestLogger(c).Error("validation failed", zap.Error(err))
		return models.SendError(c, fiber.StatusBadRequest, err.Error(), models.ErrCodeValidationFailed, middleware.GetRequestID(c))
	}

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		middleware.GetRequestLogger(c).Error("invalid date format", zap.Error(err))
		return models.SendBadRequest(c, "Invalid date format, use YYYY-MM-DD", middleware.GetRequestID(c))
	}

	user, err := h.repo.Update(c.Context(), int32(id), req.Name, dob)
	if err != nil {
		middleware.GetRequestLogger(c).Error("update user failed", zap.Error(err))
		return models.SendNotFound(c, "User not found", middleware.GetRequestID(c))
	}

	middleware.GetRequestLogger(c).Info("user updated", zap.Int32("id", user.ID))

	return c.JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return models.SendBadRequest(c, "Invalid user ID", middleware.GetRequestID(c))
	}

	if err := h.repo.Delete(c.Context(), int32(id)); err != nil {
		middleware.GetRequestLogger(c).Error("delete user failed", zap.Error(err))
		return models.SendNotFound(c, "User not found", middleware.GetRequestID(c))
	}

	middleware.GetRequestLogger(c).Info("user deleted", zap.Int("id", id))

	return c.SendStatus(204)
}
