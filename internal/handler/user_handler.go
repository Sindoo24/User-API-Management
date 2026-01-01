package handler

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

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
		h.logger.Error("failed to parse request body", zap.Error(err))
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
	}

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		h.logger.Error("invalid date format", zap.Error(err))
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid date format, use YYYY-MM-DD"})
	}

	user, err := h.repo.Create(c.Context(), req.Name, dob)
	if err != nil {
		h.logger.Error("create user failed", zap.Error(err))
		return c.Status(500).JSON(models.ErrorResponse{Error: "Failed to create user"})
	}

	h.logger.Info("user created", zap.Int32("id", user.ID))

	return c.Status(201).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
	})
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid user ID"})
	}

	resp, err := h.service.GetUserWithAge(c.Context(), int32(id))
	if err != nil {
		h.logger.Error("get user failed", zap.Error(err))
		return c.Status(404).JSON(models.ErrorResponse{Error: "User not found"})
	}

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
			h.logger.Error("list users paginated failed", zap.Error(err))
			return c.Status(500).JSON(models.ErrorResponse{Error: "Failed to list users"})
		}

		return c.JSON(paginatedResp)
	}

	users, err := h.service.ListUsersWithAge(c.Context())
	if err != nil {
		h.logger.Error("list users failed", zap.Error(err))
		return c.Status(500).JSON(models.ErrorResponse{Error: "Failed to list users"})
	}

	return c.JSON(users)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid user ID"})
	}

	var req models.UserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("failed to parse request body", zap.Error(err))
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		return c.Status(400).JSON(models.ErrorResponse{Error: err.Error()})
	}

	dob, err := time.Parse("2006-01-02", req.Dob)
	if err != nil {
		h.logger.Error("invalid date format", zap.Error(err))
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid date format, use YYYY-MM-DD"})
	}

	user, err := h.repo.Update(c.Context(), int32(id), req.Name, dob)
	if err != nil {
		h.logger.Error("update user failed", zap.Error(err))
		return c.Status(404).JSON(models.ErrorResponse{Error: "User not found"})
	}

	h.logger.Info("user updated", zap.Int32("id", user.ID))

	return c.JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid user ID"})
	}

	if err := h.repo.Delete(c.Context(), int32(id)); err != nil {
		h.logger.Error("delete user failed", zap.Error(err))
		return c.Status(404).JSON(models.ErrorResponse{Error: "User not found"})
	}

	h.logger.Info("user deleted", zap.Int("id", id))

	return c.SendStatus(204)
}
