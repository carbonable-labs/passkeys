package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/carbonable-labs/account/internal/db"
	"github.com/carbonable-labs/account/internal/domain"
	"github.com/carbonable-labs/account/internal/infrastructure"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

func main() {
	rpOrigins := strings.Split(os.Getenv("RELYING_PARTY_ORIGINS"), ",")
	wconfig := &webauthn.Config{
		RPDisplayName: os.Getenv("RELYING_PARTY_NAME"), // Display Name for your site
		RPID:          os.Getenv("RELYING_PARTY_ID"),   // Generally the FQDN for your site
		RPOrigins:     rpOrigins,                       // The origin URLs allowed for WebAuthn requests
	}

	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		slog.Error("Failed to create WebAuthn instance", "err", err)
		panic(err)
	}

	rpcClient, err := rpc.NewProvider(os.Getenv("RPC_URL"), ethrpc.WithHeader("x-apikey", os.Getenv("RPC_API_KEY")))
	if err != nil {
		slog.Error("failed dialing into rpc provider", "err", err)
		panic(err)
	}

	pgdb, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		panic(err)
	}

	dbClient := db.New(pgdb)
	txManager := db.NewTxManager(pgdb, dbClient)
	authManager := infrastructure.NewWebAuthnManager(webAuthn, dbClient)

	registrationHandler := domain.NewRegistrationHandler(authManager, infrastructure.NewStarknetAccountManager(rpcClient, dbClient, txManager))

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	account := e.Group("/account")

	account.POST("/register-request", func(c echo.Context) error {
		var req domain.RegisterRequestRequest
		err := c.Bind(&req)
		if err != nil {
			return fmt.Errorf("failed to bind request: %w", err)
		}
		res, err := registrationHandler.HandleRegisterRequest(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})
	account.POST("/register/:email", func(c echo.Context) error {
		req := domain.RegisterRequest{Email: c.QueryParam("email")}
		res, err := registrationHandler.HandleRegister(c.Request().Context(), req, c.Request())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})

	account.POST("/login-request", func(c echo.Context) error {
		var req domain.LoginRequestRequest
		err := c.Bind(&req)
		if err != nil {
			return fmt.Errorf("failed to bind request: %w", err)
		}
		res, err := domain.HandleLoginRequest(c.Request().Context(), authManager, req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})
	account.POST("/login", func(c echo.Context) error {
		req := domain.LoginRequest{Email: c.QueryParam("email")}
		res, err := domain.HandleLogin(c.Request().Context(), authManager, req, c.Request())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
