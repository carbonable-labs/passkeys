-- name: CreateUser :one
INSERT INTO users (
  id, email, session
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
--
-- name: GetUserByID :one
SELECT * FROM users WHERE email = $1 OR id = $1;

-- name: VerifyUser :one
UPDATE users SET verified = true, credentials = $1 WHERE email = $2 RETURNING *;

-- name: UpdateUserCredentials :one
UPDATE users SET credentials = $1 WHERE email = $2 RETURNING *;

-- name: AuthenticateUser :one
UPDATE users SET session = $1, last_login_at = now() WHERE email = $2 RETURNING *;
