package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/carbonable-labs/account/internal/db"
	"github.com/carbonable-labs/account/internal/domain"
	"github.com/oklog/ulid/v2"
)

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrAccountNotVerified = errors.New("account not verified")
)

type StarknetAccountManager struct {
	rpc *rpc.Provider
	db  *db.Queries
	txm *db.PgTxManager
}

func AccountFromDb(account []db.GetUserAccountRow) (*domain.Account, error) {
	if len(account) == 0 {
		return &domain.Account{}, ErrAccountNotFound
	}
	var deplogs []domain.DeploymentLog
	for _, d := range account {
		deplogs = append(deplogs, domain.DeploymentLog{AccountID: d.ID_3.String, Message: d.Message.String, Payload: d.Payload})
	}
	return &domain.Account{
		ID:            account[0].ID,
		UserID:        account[0].UserID,
		Address:       account[0].Address.String,
		Status:        domain.DeploymentStatus(account[0].Status.String),
		Configuration: domain.AccountConfiguration{},
		Data:          account[0].Data,
		Logs:          deplogs,
	}, nil
}

func NewStarknetAccountManager(rpc *rpc.Provider, db *db.Queries, txm *db.PgTxManager) *StarknetAccountManager {
	return &StarknetAccountManager{rpc: rpc, db: db, txm: txm}
}

func (m *StarknetAccountManager) DeployAccount(ctx context.Context, req domain.DeployAccountRequest) (domain.DeployAccountResponse, error) {
	// check user is present in our system
	u, err := m.db.GetUserByID(ctx, req.UserID)
	if err != nil {
		return domain.DeployAccountResponse{}, fmt.Errorf("user not found : %w", err)
	}

	// create account_deployment_object
	account, err := m.createAccount(ctx, u)
	if err != nil {
		return domain.DeployAccountResponse{}, fmt.Errorf("failed to create account in db: %w", err)
	}

	// deploy account to starknet
	account, err = m.Deploy(ctx, account)
	if err != nil {
		return domain.DeployAccountResponse{}, fmt.Errorf("failed to deploy account to starknet: %w", err)
	}

	return domain.DeployAccountResponse{
		Address: account.Address,
	}, nil
}

func (m *StarknetAccountManager) Deploy(ctx context.Context, acc *domain.Account) (*domain.Account, error) {
	// Set account.Status to STARTED
	err := m.db.StartAccountDeployment(ctx, acc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to set account deployment status to STARTED: %w", err)
	}

	// Account deployment will be done in a multistep process
	// as documented here : https://github.com/NethermindEth/starknet.go/blob/20e049b218e294daa558788712e9a2cc551f4400/examples/deployAccount/main.go#L4

	// Generate set of keypairs
	// Store keystore in acc.Data (use some kind of https://github.com/ethereum/go-ethereum/tree/master/accounts/keystore encryption)

	// create deployment tx with keypairs, precompiled account class hash (braavos or argent x)
	// sign tx and deploy

	// provide user the address to send funds to

	// Verifier will check account is deployed, has funds and send tx
	// Verifier will change account.Status to DEPLOYED when account tx is sent

	return nil, nil
}

func (m *StarknetAccountManager) Verify(ctx context.Context, acc domain.Account) (domain.VerificationStatus, error) {
	switch acc.Status {
	case domain.DEPLOYED:
	// check account is deployed
	case domain.ERROR:
	// try deploy account
	case domain.INIT:
	// wait, account was not yet picked up to deploy
	case domain.PENDING:
	// try to deploy account
	case domain.STARTED:
	// wait, account was picked up but not deployed yet
	default:
		panic(fmt.Sprintf("unexpected domain.DeploymentStatus: %#v", acc.Status))
	}
	return domain.VerificationStatus{}, ErrAccountNotVerified
}

func (m *StarknetAccountManager) createAccount(ctx context.Context, u db.User) (*domain.Account, error) {
	acc := domain.NewAccount(u.ID, domain.AccountConfiguration{})
	deploymentID := ulid.Make().String()
	deploymentLog := domain.DeploymentLog{AccountID: acc.ID, Message: "Account deployment started", Payload: []byte(`{"account_id":"` + acc.ID + `", "user_id":"` + acc.UserID + `", "account_deployment_id": "` + deploymentID + `"}`)}

	acc.Status = domain.PENDING
	acc.Logs = append(acc.Logs, deploymentLog)

	acc, err := m.accountDbTx(ctx, acc, deploymentID, deploymentLog)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (m *StarknetAccountManager) accountDbTx(ctx context.Context, acc *domain.Account, deploymentID string, deploymentLog domain.DeploymentLog) (*domain.Account, error) {
	err := m.txm.DoTx(ctx, func(qtx *db.Queries) error {
		_, err := m.db.CreateAccount(ctx, db.CreateAccountParams{
			ID:     acc.ID,
			UserID: acc.UserID,
			Data:   []byte("{}"),
		})
		if err != nil {
			return err
		}
		_, err = m.db.CreateAccountDeployment(ctx, db.CreateAccountDeploymentParams{
			ID:        deploymentID,
			UserID:    acc.UserID,
			AccountID: acc.ID,
			Status:    string(acc.Status),
		})
		if err != nil {
			return err
		}
		_, err = m.db.CreateDeploymentLog(ctx, db.CreateDeploymentLogParams{
			ID:        ulid.Make().String(),
			AccountID: acc.ID,
			Message:   deploymentLog.Message,
			Payload:   deploymentLog.Payload,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create account : %w", err)
	}
	return acc, nil
}
