package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/ftamy9/dj_shorturl/internal/auth"
	"github.com/ftamy9/dj_shorturl/internal/config"
	"github.com/ftamy9/dj_shorturl/internal/database"
	"github.com/ftamy9/dj_shorturl/internal/shortener"
)

func main() {
	cfg := config.Load()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		pool, err := database.Connect(cfg)
		if err != nil {
			slog.Error("database connection failed", "error", err)
			os.Exit(1)
		}
		defer pool.Close()
		if err := runMigrations(pool); err != nil {
			slog.Error("migration failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations completed")
		return
	}

	slog.Info("connecting to database", "host", cfg.DBHost)
	pool, err := database.Connect(cfg)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("database connected")

	if err := runMigrations(pool); err != nil {
		slog.Error("auto-migration failed", "error", err)
		os.Exit(1)
	}

	userRepo := auth.NewRepository(pool)
	addrRepo := shortener.NewRepository(pool)

	authHandler := auth.NewHandler(userRepo, cfg)
	addrHandler := shortener.NewHandler(addrRepo)
	authMw := auth.Middleware(userRepo, cfg)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/signup", authHandler.ServeHTTP)
	mux.HandleFunc("PUT /user/signup", authHandler.ServeHTTP)

	mux.HandleFunc("POST /shorter/url", authMw(http.HandlerFunc(addrHandler.HandleCreate)).ServeHTTP)

	mux.HandleFunc("GET /shorter/url/", addrHandler.HandleRedirect)

	addr := ":" + cfg.ServerPort
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func runMigrations(pool interface{}) error {
	slog.Info("running migrations")

	type executor interface {
		Exec(ctx context.Context, sql string, args ...any) (any, error)
	}
	db, ok := pool.(executor)
	if !ok {
		return fmt.Errorf("pool does not implement Exec")
	}

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS base_users (
			id VARCHAR(10) PRIMARY KEY,
			password_hash VARCHAR(128) NOT NULL,
			create_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
		`CREATE TABLE IF NOT EXISTS addresses (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			url VARCHAR(32779) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(context.Background(), m); err != nil {
			return fmt.Errorf("migration failed: %s: %w", truncate(m, 60), err)
		}
	}

	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return strings.TrimSpace(s[:max]) + "..."
}
