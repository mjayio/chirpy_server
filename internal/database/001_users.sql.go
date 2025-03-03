// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: 001_users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createChirp = `-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING id, created_at, updated_at, user_id, body
`

type CreateChirpParams struct {
	UserID uuid.UUID
	Body   string
}

func (q *Queries) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, createChirp, arg.UserID, arg.Body)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Body,
	)
	return i, err
}

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRefreshTokenParams struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.Token, arg.UserID, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, token)
VALUES ($1, NOW(), NOW(), $2, $3, $4)
RETURNING id, created_at, updated_at, email, hashed_password, token, is_chirpy_red
`

type CreateUserParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
	Token          string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.HashedPassword,
		arg.Token,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.Token,
		&i.IsChirpyRed,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const deleteChirp = `-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1
`

func (q *Queries) DeleteChirp(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteChirp, id)
	return err
}

const getChirp = `-- name: GetChirp :one
SELECT id, created_at, updated_at, user_id, body FROM chirps WHERE id = $1
`

func (q *Queries) GetChirp(ctx context.Context, id uuid.UUID) (Chirp, error) {
	row := q.db.QueryRowContext(ctx, getChirp, id)
	var i Chirp
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Body,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, token, is_chirpy_red FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.Token,
		&i.IsChirpyRed,
	)
	return i, err
}

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT u.id, u.created_at, u.updated_at, u.email, u.hashed_password, u.token, u.is_chirpy_red FROM users u
JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = $1
AND rt.revoked_at IS NULL
AND rt.created_at > NOW() - INTERVAL '60 days'
LIMIT 1
`

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.Token,
		&i.IsChirpyRed,
	)
	return i, err
}

const listChirps = `-- name: ListChirps :many
SELECT id, created_at, updated_at, user_id, body FROM chirps
ORDER BY created_at ASC
`

func (q *Queries) ListChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, listChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listChirpsByAuthor = `-- name: ListChirpsByAuthor :many
SELECT id, created_at, updated_at, user_id, body FROM chirps WHERE user_id = $1
Order by created_at ASC
`

func (q *Queries) ListChirpsByAuthor(ctx context.Context, userID uuid.UUID) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, listChirpsByAuthor, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const makeChirpyRed = `-- name: MakeChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, token, is_chirpy_red
`

func (q *Queries) MakeChirpyRed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, makeChirpyRed, id)
	return err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, revokeRefreshToken, token)
	return err
}

const updateUserEmailPassword = `-- name: UpdateUserEmailPassword :one
UPDATE users
SET email = $1, hashed_password = $2
WHERE id = $3
RETURNING id, created_at, updated_at, email, hashed_password, token, is_chirpy_red
`

type UpdateUserEmailPasswordParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

func (q *Queries) UpdateUserEmailPassword(ctx context.Context, arg UpdateUserEmailPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserEmailPassword, arg.Email, arg.HashedPassword, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.Token,
		&i.IsChirpyRed,
	)
	return i, err
}
