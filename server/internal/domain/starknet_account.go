package domain

import (
	"context"
	"encoding/json"

	"github.com/oklog/ulid/v2"
)

type (
	DeploymentStatus string

	DeployAccountRequest struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
	}
	DeployAccountResponse struct {
		Address string `json:"address"`
	}
	VerificationStatus struct {
		AccountID string `json:"account_id"`
		Status    string `json:"status"`
	}

	AccountManager interface {
		DeployAccount(ctx context.Context, req DeployAccountRequest) (DeployAccountResponse, error)
	}
	AccountDeployer interface {
		Deploy(ctx context.Context, acc *Account) (*Account, error)
	}
	AccountVerifier interface {
		Verify(ctx context.Context, acc Account) (VerificationStatus, error)
	}
)

const (
	INIT     DeploymentStatus = "INIT"
	PENDING  DeploymentStatus = "PENDING"
	STARTED  DeploymentStatus = "STARTED"
	DEPLOYED DeploymentStatus = "DEPLOYED"
	ERROR    DeploymentStatus = "ERROR"
)

type AccountConfiguration struct{}

type DeploymentLog struct {
	AccountID string          `json:"account_id"`
	Message   string          `json:"message"`
	Payload   json.RawMessage `json:"payload"`
}

type Account struct {
	ID            string               `json:"id"`
	UserID        string               `json:"user_id"`
	Address       string               `json:"address"`
	Status        DeploymentStatus     `json:"status"`
	Configuration AccountConfiguration `json:"configuration"`
	Data          json.RawMessage      `json:"data"`
	Logs          []DeploymentLog      `json:"logs"`
}

func NewAccount(userID string, configuration AccountConfiguration) *Account {
	return &Account{
		ID:            ulid.Make().String(),
		UserID:        userID,
		Configuration: configuration,
		Status:        INIT,
	}
}
