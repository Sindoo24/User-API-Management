package service

import (
	"context"
	"time"

	"BACKEND/internal/models"
	"BACKEND/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func calculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}

	return age
}

func (s *UserService) GetUserWithAge(ctx context.Context, id int32) (*models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.UserWithAgeResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Time.Format("2006-01-02"),
		Age:  calculateAge(user.Dob.Time),
	}, nil
}

func (s *UserService) ListUsersWithAge(ctx context.Context) ([]models.UserWithAgeResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Convert DB models to response models with ages
	result := make([]models.UserWithAgeResponse, len(users))
	for i, user := range users {
		result[i] = models.UserWithAgeResponse{
			ID:   user.ID,
			Name: user.Name,
			Dob:  user.Dob.Time.Format("2006-01-02"),
			Age:  calculateAge(user.Dob.Time),
		}
	}

	return result, nil
}

func (s *UserService) ListUsersWithAgePaginated(ctx context.Context, page, limit int) (*models.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	users, err := s.repo.ListPaginated(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}
	data := make([]models.UserWithAgeResponse, len(users))
	for i, user := range users {
		data[i] = models.UserWithAgeResponse{
			ID:   user.ID,
			Name: user.Name,
			Dob:  user.Dob.Time.Format("2006-01-02"),
			Age:  calculateAge(user.Dob.Time),
		}
	}
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &models.PaginatedUsersResponse{
		Data: data,
		Pagination: models.PaginationMeta{
			Total:       total,
			Page:        page,
			Limit:       limit,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
	}, nil
}
