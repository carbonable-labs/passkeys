package domain

import (
	"context"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
)

type (
	RegisterRequestRequest struct {
		Email string `json:"email"`
	}
	RegisterRequestResponse struct {
		Options *protocol.CredentialCreation `json:"options"`
	}
	RegisterRequest struct {
		Email string `json:"email"`
	}
	RegisterResponse struct{}

	Registrator interface {
		BeginRegistration(ctx context.Context, req RegisterRequestRequest) (RegisterRequestResponse, error)
		FinishRegistration(ctx context.Context, req RegisterRequest, httpReq *http.Request) (RegisterResponse, error)
	}
)

func HandleRegisterRequest(ctx context.Context, authManager Registrator, req RegisterRequestRequest) (RegisterRequestResponse, error) {
	res, err := authManager.BeginRegistration(ctx, req)
	if err != nil {
		return RegisterRequestResponse{}, err
	}

	return res, nil
}

func HandleRegister(ctx context.Context, authManager Registrator, req RegisterRequest, httpReq *http.Request) (RegisterResponse, error) {
	res, err := authManager.FinishRegistration(ctx, req, httpReq)
	if err != nil {
		return RegisterResponse{}, err
	}
	return res, nil
}
