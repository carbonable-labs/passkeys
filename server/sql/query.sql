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

-- name: GetUserAccount :many
SELECT a.*, ad.*, dl.* FROM accounts a 
LEFT JOIN account_deployments ad ON a.id = ad.account_id 
LEFT JOIN account_deployment_logs dl ON a.id = dl.account_id
WHERE a.user_id = $1;

-- name: CreateAccount :one
INSERT INTO accounts (
  id, user_id, data
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: CreateAccountDeployment :one
INSERT INTO account_deployments (
  id, user_id, account_id, status
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: CreateDeploymentLog :one
INSERT INTO account_deployment_logs (
  id, account_id, message, payload
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: StartAccountDeployment :exec
UPDATE account_deployments SET status = 'STARTED' WHERE account_id = $1;
