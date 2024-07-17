// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package db

import (
	"context"
	"encoding/json"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  id, email, session
) VALUES (
  $1, $2, $3
)
RETURNING id, email, session, credentials, verified
`

type CreateUserParams struct {
	ID      string
	Email   string
	Session json.RawMessage
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.ID, arg.Email, arg.Session)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Session,
		&i.Credentials,
		&i.Verified,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, session, credentials, verified FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Session,
		&i.Credentials,
		&i.Verified,
	)
	return i, err
}

const verifyUser = `-- name: VerifyUser :one
UPDATE users SET verified = true, credentials = $1 WHERE email = $2 RETURNING id, email, session, credentials, verified
`

type VerifyUserParams struct {
	Credentials json.RawMessage
	Email       string
}

func (q *Queries) VerifyUser(ctx context.Context, arg VerifyUserParams) (User, error) {
	row := q.db.QueryRow(ctx, verifyUser, arg.Credentials, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Session,
		&i.Credentials,
		&i.Verified,
	)
	return i, err
}
