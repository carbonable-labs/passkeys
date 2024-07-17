package domain

import (
	"context"

	"github.com/go-webauthn/webauthn/protocol"
)

type (
	RegisterRequestRequest struct {
		Email string `json:"email"`
	}
	RegisterRequestResponse struct {
		Options *protocol.CredentialCreation `json:"options"`
	}
	RegisterRequest  struct{}
	RegisterResponse struct{}

	Registrator interface {
		BeginRegistration(ctx context.Context, req RegisterRequestRequest) (RegisterRequestResponse, error)
		FinishRegistration(ctx context.Context, req RegisterRequest) (RegisterResponse, error)
	}
)

func HandleRegisterRequest(ctx context.Context, authManager Registrator, req RegisterRequestRequest) (RegisterRequestResponse, error) {
	res, err := authManager.BeginRegistration(ctx, req)
	if err != nil {
		return RegisterRequestResponse{}, err
	}

	return res, nil
}

func HandleRegister(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	return RegisterResponse{}, nil
}
