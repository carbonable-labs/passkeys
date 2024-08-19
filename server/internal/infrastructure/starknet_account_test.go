package infrastructure

import (
	"context"
	"testing"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/carbonable-labs/account/internal/db"
	"github.com/carbonable-labs/account/internal/domain"

	"github.com/stretchr/testify/assert"
)

func createUser(t *testing.T, ctx context.Context, query *db.Queries, u domain.AuthUser) error {
	t.Helper()
	_, err := query.CreateUser(ctx, db.CreateUserParams{
		ID:      u.Id,
		Email:   u.Email,
		Session: []byte("{}"),
	})
	if err != nil {
		return err
	}
	return nil
}

func TestDeployAccountUserExist(t *testing.T) {
	rpcClient, err := rpc.NewProvider("https://free-rpc.nethermind.io/sepolia-juno")
	if err != nil {
		t.Errorf("failed to dial in rpc provider : %s", err)
	}

	db, txm := db.NewTestDB(t)

	am := NewStarknetAccountManager(rpcClient, db, txm)
	_, err = am.DeployAccount(context.Background(), domain.DeployAccountRequest{UserID: "fake-user-id", Email: "fake@email.com"})

	// we should have an error if user is not present in our system
	assert.NotNil(t, err)
}

func TestCreateAccount(t *testing.T) {
	rpcClient, err := rpc.NewProvider("https://free-rpc.nethermind.io/sepolia-juno")
	if err != nil {
		t.Errorf("failed to dial in rpc provider : %s", err)
	}

	d, txm := db.NewTestDB(t)
	am := NewStarknetAccountManager(rpcClient, d, txm)

	ctx := context.Background()
	user := domain.AuthUser{Id: "fake-user-id", Email: "fake@email.com"}
	err = createUser(t, ctx, d, user)
	if err != nil {
		t.Fatal(err)
	}

	account := domain.NewAccount(user.Id, domain.AccountConfiguration{})
	assert.Equal(t, domain.INIT, account.Status)

	acc, err := am.createAccount(ctx, db.User{ID: "fake-user-id", Email: "fake@email.com"})
	if err != nil {
		t.Fatalf("failed to deploy account : %s", err)
	}

	assert.Nil(t, err)
	assert.Equal(t, domain.PENDING, acc.Status)

	a, err := d.GetUserAccount(ctx, "fake-user-id")
	if err != nil {
		t.Fatalf("failed to get account from db : %s", err)
	}
	acc, err = AccountFromDb(a)
	if err != nil {
		t.Fatalf("failed convert db model to domain model : %s", err)
	}
	assert.Nil(t, err)
	assert.Equal(t, domain.PENDING, acc.Status)
	assert.Equal(t, 1, len(acc.Logs))
}

func TestDeployAccount(t *testing.T) {
	rpcClient, err := rpc.NewProvider("https://free-rpc.nethermind.io/sepolia-juno")
	if err != nil {
		t.Errorf("failed to dial in rpc provider : %s", err)
	}

	d, txm := db.NewTestDB(t)
	am := NewStarknetAccountManager(rpcClient, d, txm)

	// create user + account
	ctx := context.Background()
	user := domain.AuthUser{Id: "fake-user-id", Email: "fake@email.com"}
	err = createUser(t, ctx, d, user)
	if err != nil {
		t.Fatal(err)
	}
	account := domain.NewAccount(user.Id, domain.AccountConfiguration{})

	// deploy account
	_, err = am.Deploy(ctx, account)
	if err != nil {
		t.Fatalf("failed to deploy account : %s", err)
	}
}
