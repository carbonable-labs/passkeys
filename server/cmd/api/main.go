package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/carbonable-labs/account/internal/db"
	"github.com/carbonable-labs/account/internal/domain"
	"github.com/carbonable-labs/account/internal/infrastructure"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func main() {
	wconfig := &webauthn.Config{
		RPDisplayName: "Carbonable",                                                        // Display Name for your site
		RPID:          "carbonable.io",                                                     // Generally the FQDN for your site
		RPOrigins:     []string{"https://www.carbonable.io", "https://auth.carbonable.io"}, // The origin URLs allowed for WebAuthn requests
	}

	webAuthn, err := webauthn.New(wconfig)
	if err != nil {
		slog.Error("Failed to create WebAuthn instance", "err", err)
		return
	}

	pgdb, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		return
	}
	dbClient := db.New(pgdb)
	authManager := infrastructure.NewWebAuthnManager(webAuthn, dbClient)

	e := echo.New()
	account := e.Group("/account")

	account.POST("/register-request", func(c echo.Context) error {
		var req domain.RegisterRequestRequest
		err := c.Bind(&req)
		if err != nil {
			return fmt.Errorf("failed to bind request: %w", err)
		}
		res, err := domain.HandleRegisterRequest(c.Request().Context(), authManager, req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})
	account.POST("/register", func(c echo.Context) error {
		var req domain.RegisterRequest
		err := c.Bind(&req)
		if err != nil {
			return fmt.Errorf("failed to bind request: %w", err)
		}
		res, err := domain.HandleRegister(c.Request().Context(), req)
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
		res, err := domain.HandleLoginRequest(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})
	account.POST("/login", func(c echo.Context) error {
		var req domain.LoginRequest
		err := c.Bind(&req)
		if err != nil {
			return fmt.Errorf("failed to bind request: %w", err)
		}
		res, err := domain.HandleLogin(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
