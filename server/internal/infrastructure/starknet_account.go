package infrastructure

import (
	"context"
	"fmt"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/carbonable-labs/account/internal/db"
	"github.com/carbonable-labs/account/internal/domain"
)

type StarknetAccountManager struct {
	rpc *rpc.Provider
	db  *db.Queries
}

func NewStarknetAccountManager(rpc *rpc.Provider, db *db.Queries) *StarknetAccountManager {
	return &StarknetAccountManager{rpc: rpc, db: db}
}

func (m *StarknetAccountManager) DeployAccount(ctx context.Context, req domain.DeployAccountRequest) (domain.DeployAccountResponse, error) {
	_, err := m.db.GetUserByID(ctx, req.UserID)
	if err != nil {
		return domain.DeployAccountResponse{}, fmt.Errorf("user not found : %w", err)
	}

	return domain.DeployAccountResponse{}, nil
}
