package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carbonable-labs/account/internal/db"
	"github.com/carbonable-labs/account/internal/domain"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnManager struct {
	a  *webauthn.WebAuthn
	db *db.Queries
}

type WebAuthnUser struct {
	user        domain.AuthUser
	session     *webauthn.SessionData
	credentials webauthn.Credential
}

func WebAuthnUserFromDb(user db.User) WebAuthnUser {
	var sess *webauthn.SessionData
	data, _ := user.Session.MarshalJSON()
	_ = json.Unmarshal(data, &sess)
	var creds webauthn.Credential
	c, _ := user.Credentials.MarshalJSON()
	_ = json.Unmarshal(c, &creds)
	return WebAuthnUser{user: domain.AuthUser{Id: user.ID, Email: user.Email}, session: sess, credentials: creds}
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
	return []webauthn.Credential{u.credentials}
}

func (u WebAuthnUser) WebAuthnIcon() string {
	return ""
}

func (u *WebAuthnUser) AddCredential(cred webauthn.Credential) {
	u.credentials = cred
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

func (m *WebAuthnManager) FinishRegistration(ctx context.Context, req domain.RegisterRequest, httpReq *http.Request) (domain.RegisterResponse, error) {
	user, err := m.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return domain.RegisterResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	wu := WebAuthnUserFromDb(user)
	creds, err := m.a.FinishRegistration(wu, *wu.session, httpReq)
	if err != nil {
		return domain.RegisterResponse{}, fmt.Errorf("failed to finish registration: %w", err)
	}
	wu.AddCredential(*creds)

	by, err := json.Marshal(creds)
	if err != nil {
		return domain.RegisterResponse{}, fmt.Errorf("failed to marshal credentials data : %w", err)
	}
	_, err = m.db.VerifyUser(ctx, db.VerifyUserParams{
		Credentials: json.RawMessage(by),
		Email:       wu.user.Email,
	})
	if err != nil {
		return domain.RegisterResponse{}, fmt.Errorf("failed to save user credentials : %w", err)
	}

	return domain.RegisterResponse{}, nil
}

func (m *WebAuthnManager) BeginLogin(ctx context.Context, req domain.LoginRequestRequest) (domain.LoginRequestResponse, error) {
	user, err := m.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return domain.LoginRequestResponse{}, fmt.Errorf("failed to get user: %w", err)
	}
	wu := WebAuthnUserFromDb(user)
	fmt.Printf("user: %+v\n", wu)
	options, session, err := m.a.BeginLogin(wu)
	if err != nil {
		return domain.LoginRequestResponse{}, fmt.Errorf("failed to begin authentication: %w", err)
	}

	wu.session = session
	by, err := json.Marshal(wu.session)
	if err != nil {
		return domain.LoginRequestResponse{}, fmt.Errorf("failed to marshal session data : %w", err)
	}
	_, err = m.db.AuthenticateUser(ctx, db.AuthenticateUserParams{
		Session: json.RawMessage(by),
		Email:   wu.user.Email,
	})
	if err != nil {
		return domain.LoginRequestResponse{}, fmt.Errorf("failed to save user session : %w", err)
	}

	return domain.LoginRequestResponse{Options: options}, nil
}

func (m *WebAuthnManager) FinishLogin(ctx context.Context, req domain.LoginRequest, httpReq *http.Request) (domain.LoginResponse, error) {
	user, err := m.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("failed to get user: %w", err)
	}
	wu := WebAuthnUserFromDb(user)
	creds, err := m.a.FinishLogin(wu, *wu.session, httpReq)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("failed to finish login: %w", err)
	}

	wu.AddCredential(*creds)

	by, err := json.Marshal(creds)
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("failed to marshal credentials data : %w", err)
	}
	_, err = m.db.VerifyUser(ctx, db.VerifyUserParams{
		Credentials: json.RawMessage(by),
		Email:       wu.user.Email,
	})
	if err != nil {
		return domain.LoginResponse{}, fmt.Errorf("failed to save user credentials : %w", err)
	}

	return domain.LoginResponse{Credentials: *creds}, nil
}
