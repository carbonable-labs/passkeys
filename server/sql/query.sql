-- name: CreateUser :one
INSERT INTO users (
  id, email, session
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserByEmail :one
