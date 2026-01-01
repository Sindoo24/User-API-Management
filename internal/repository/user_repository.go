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

func (r *UserRepository) Create(ctx context.Context, name string, dob time.Time) (generated.User, error) {
	return r.queries.CreateUser(ctx, generated.CreateUserParams{
		Name: name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
	})
}

func (r *UserRepository) GetByID(ctx context.Context, id int32) (generated.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

func (r *UserRepository) List(ctx context.Context) ([]generated.User, error) {
	return r.queries.ListUsers(ctx)
}

func (r *UserRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (generated.User, error) {
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

func (r *UserRepository) ListPaginated(ctx context.Context, limit, offset int32) ([]generated.User, error) {
	return r.queries.ListUsersPaginated(ctx, generated.ListUsersPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}
