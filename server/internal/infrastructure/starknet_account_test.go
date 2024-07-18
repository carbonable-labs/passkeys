package infrastructure

import (
	"context"
	"testing"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/carbonable-labs/account/internal/db"
	"github.com/carbonable-labs/account/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestDeployAccountUserExist(t *testing.T) {
	rpcClient, err := rpc.NewProvider("https://free-rpc.nethermind.io/sepolia-juno")
	if err != nil {
		t.Errorf("failed to dial in rpc provider : %s", err)
	}

	db := db.NewTestDB(t)

	am := NewStarknetAccountManager(rpcClient, db)
	_, err = am.DeployAccount(context.Background(), domain.DeployAccountRequest{UserID: "fake-user-id", Email: "fake@email.com"})

	// we should have an error if user is not present in our system
	assert.NotNil(t, err)
}
