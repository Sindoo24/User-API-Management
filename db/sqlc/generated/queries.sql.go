package generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countUsers = `-- name: CountUsers :one
SELECT COUNT(*) 
FROM users
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, dob) 
VALUES ($1, $2) 
RETURNING id, name, dob
`

type CreateUserParams struct {
	Name string      `json:"name"`
	Dob  pgtype.Date `json:"dob"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Dob)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Dob)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, dob 
FROM users 
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Dob)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, dob 
FROM users 
ORDER BY id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Name, &i.Dob); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsersPaginated = `-- name: ListUsersPaginated :many
SELECT id, name, dob 
FROM users 
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListUsersPaginatedParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsersPaginated(ctx context.Context, arg ListUsersPaginatedParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsersPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Name, &i.Dob); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users 
SET name = $2, dob = $3 
WHERE id = $1 
RETURNING id, name, dob
`

type UpdateUserParams struct {
	ID   int32       `json:"id"`
	Name string      `json:"name"`
	Dob  pgtype.Date `json:"dob"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser, arg.ID, arg.Name, arg.Dob)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Dob)
	return i, err
}
