package domain

import "context"

type (
	DeployAccountRequest struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
	}
	DeployAccountResponse struct {
		Address string `json:"address"`
	}
	AccountManager interface {
		DeployAccount(ctx context.Context, req DeployAccountRequest) (DeployAccountResponse, error)
	}
)
