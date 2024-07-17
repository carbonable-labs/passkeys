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
RETURNING id, email, session
`

type CreateUserParams struct {
	ID      string
	Email   string
	Session json.RawMessage
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.ID, arg.Email, arg.Session)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.Session)
	return i, err
}