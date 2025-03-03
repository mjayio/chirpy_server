-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, token)
VALUES ($1, NOW(), NOW(), $2, $3, $4)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: ListChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users u
JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = $1
AND rt.revoked_at IS NULL
AND rt.created_at > NOW() - INTERVAL '60 days'
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1;

-- name: UpdateUserEmailPassword :one
UPDATE users
SET email = $1, hashed_password = $2
WHERE id = $3
RETURNING *;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;

-- name: MakeChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;

-- name: ListChirpsByAuthor :many
SELECT * FROM chirps WHERE user_id = $1
Order by created_at ASC;