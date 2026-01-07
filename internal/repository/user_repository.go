package repository

import (
	"BACKEND/db/sqlc/generated"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	queries *generated.Queries
}

func NewUserRepository(q *generated.Queries) *UserRepository {
	return &UserRepository{queries: q}
}

func (r *UserRepository) Create(ctx context.Context, name string, dob time.Time) (generated.CreateUserRow, error) {
	return r.queries.CreateUser(ctx, generated.CreateUserParams{
		Name: name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
	})
}

// CreateWithAuth creates a new user with authentication fields
func (r *UserRepository) CreateWithAuth(ctx context.Context, name, email, passwordHash, role string, dob time.Time) (generated.CreateUserRow, error) {
	return r.queries.CreateUser(ctx, generated.CreateUserParams{
		Name: name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
		Email:        email,
		PasswordHash: passwordHash,
		Column5:      role, // role parameter (COALESCE in SQL)
	})
}

func (r *UserRepository) GetByID(ctx context.Context, id int32) (generated.GetUserByIDRow, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *UserRepository) List(ctx context.Context) ([]generated.ListUsersRow, error) {
	return r.queries.ListUsers(ctx)
}

func (r *UserRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (generated.UpdateUserRow, error) {
	return r.queries.UpdateUser(ctx, generated.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
	})
}

func (r *UserRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *UserRepository) ListPaginated(ctx context.Context, limit, offset int32) ([]generated.ListUsersPaginatedRow, error) {
	return r.queries.ListUsersPaginated(ctx, generated.ListUsersPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

// GetByEmail retrieves a user by email address
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (generated.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}
