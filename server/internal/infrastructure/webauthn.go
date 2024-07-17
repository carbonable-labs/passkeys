package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/carbonable-labs/account/internal/db"
	"github.com/carbonable-labs/account/internal/domain"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnManager struct {
	a  *webauthn.WebAuthn
	db *db.Queries
}

type WebAuthnUser struct {
	user    domain.AuthUser
	session *webauthn.SessionData
}

func (u WebAuthnUser) WebAuthnID() []byte {
	return []byte(u.user.Email)
}

func (u WebAuthnUser) WebAuthnName() string {
	return u.user.Email
}

func (u WebAuthnUser) WebAuthnDisplayName() string {
	return u.user.Email
}

func (u WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return []webauthn.Credential{}
}

func (u WebAuthnUser) WebAuthnIcon() string {
	return ""
}

func NewWebAuthnManager(a *webauthn.WebAuthn, dbClient *db.Queries) *WebAuthnManager {
	return &WebAuthnManager{a: a, db: dbClient}
}

func (m *WebAuthnManager) BeginRegistration(ctx context.Context, req domain.RegisterRequestRequest) (domain.RegisterRequestResponse, error) {
	u, err := domain.NewUser(req.Email)
	if err != nil {
		return domain.RegisterRequestResponse{}, fmt.Errorf("failed to create user: %w", err)
	}
	wu := WebAuthnUser{user: *u}
	options, session, err := m.a.BeginRegistration(wu)
	if err != nil {
		return domain.RegisterRequestResponse{}, fmt.Errorf("failed to begin registration: %w", err)
	}

	wu.session = session
	by, err := json.Marshal(wu.session)
	if err != nil {
		return domain.RegisterRequestResponse{}, fmt.Errorf("failed to marshal session data : %w", err)
	}

	_, err = m.db.CreateUser(ctx, db.CreateUserParams{
		ID:      wu.user.Id,
		Email:   wu.user.Email,
		Session: json.RawMessage(by),
	})
	if err != nil {
		return domain.RegisterRequestResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	return domain.RegisterRequestResponse{Options: options}, nil
}

func (m *WebAuthnManager) FinishRegistration(ctx context.Context, req domain.RegisterRequest) (domain.RegisterResponse, error) {
	return domain.RegisterResponse{}, nil
}
