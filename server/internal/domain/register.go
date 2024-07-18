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
	RegisterResponse struct {
		// we do not expose this field to the client
		// but we need it to deploy the account
		userID string
	}

	Registrator interface {
		BeginRegistration(ctx context.Context, req RegisterRequestRequest) (RegisterRequestResponse, error)
		FinishRegistration(ctx context.Context, req RegisterRequest, httpReq *http.Request) (RegisterResponse, error)
	}

	RegistrationHandler struct {
		r  Registrator
		am AccountManager
	}
)

func (regRes *RegisterResponse) WithUserID(userID string) *RegisterResponse {
	regRes.userID = userID
	return regRes
}

func NewRegistrationHandler(r Registrator, sn AccountManager) *RegistrationHandler {
	return &RegistrationHandler{r: r, am: sn}
}

func (h *RegistrationHandler) HandleRegisterRequest(ctx context.Context, req RegisterRequestRequest) (RegisterRequestResponse, error) {
	res, err := h.r.BeginRegistration(ctx, req)
	if err != nil {
		return RegisterRequestResponse{}, err
	}

	return res, nil
}

func (h *RegistrationHandler) HandleRegister(ctx context.Context, req RegisterRequest, httpReq *http.Request) (RegisterResponse, error) {
	res, err := h.r.FinishRegistration(ctx, req, httpReq)
	if err != nil {
		return RegisterResponse{}, err
	}
	// user sucessfully registered
	// deploy starknet account
	go h.am.DeployAccount(ctx, DeployAccountRequest{UserID: res.userID, Email: req.Email})

	return res, nil
}
