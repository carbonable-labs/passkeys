package domain

import (
	"context"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type (
	LoginRequestRequest struct {
		Email string `json:"email"`
	}
	LoginRequestResponse struct {
		Options *protocol.CredentialAssertion `json:"options"`
	}
	LoginRequest struct {
		Email string `json:"email"`
	}
	LoginResponse struct {
		Credentials webauthn.Credential `json:"credentials"`
	}

	Authenticator interface {
		BeginLogin(ctx context.Context, req LoginRequestRequest) (LoginRequestResponse, error)
		FinishLogin(ctx context.Context, req LoginRequest, httpReq *http.Request) (LoginResponse, error)
	}
)

func HandleLoginRequest(ctx context.Context, authManager Authenticator, req LoginRequestRequest) (LoginRequestResponse, error) {
	res, err := authManager.BeginLogin(ctx, req)
	if err != nil {
		return LoginRequestResponse{}, err
	}
	return res, nil
}

func HandleLogin(ctx context.Context, authManager Authenticator, req LoginRequest, httpReq *http.Request) (LoginResponse, error) {
	res, err := authManager.FinishLogin(ctx, req, httpReq)
	if err != nil {
		return LoginResponse{}, err
	}

	return res, nil
}
