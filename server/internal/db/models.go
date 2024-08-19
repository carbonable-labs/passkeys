// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        string
	UserID    string
	Address   pgtype.Text
	Data      json.RawMessage
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type AccountDeployment struct {
	ID        string
	UserID    string
	AccountID string
	Status    string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type AccountDeploymentLog struct {
	ID        string
	AccountID string
	Message   string
	Payload   json.RawMessage
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type User struct {
	ID          string
	Email       string
	Session     json.RawMessage
	Credentials json.RawMessage
	Verified    bool
	LastLoginAt pgtype.Timestamp
}
