package domain

import "context"

type (
	LoginRequestRequest struct {
		Email string `json:"email"`
	}
	LoginRequestResponse struct{}
	LoginRequest         struct{}
	LoginResponse        struct{}

	Authenticator interface {
		BeginLogin(ctx context.Context, req LoginRequestRequest) (LoginRequestResponse, error)
		FinishLogin(ctx context.Context, req LoginRequestRequest) (LoginResponse, error)
	}
)

func HandleLoginRequest(ctx context.Context, req LoginRequestRequest) (LoginRequestResponse, error) {
	return LoginRequestResponse{}, nil
}

func HandleLogin(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	return LoginResponse{}, nil
}
