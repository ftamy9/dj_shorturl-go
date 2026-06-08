package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/ftamy9/dj_shorturl-go/internal/auth"
	"github.com/ftamy9/dj_shorturl-go/internal/config"
	"github.com/ftamy9/dj_shorturl-go/internal/database"
	"github.com/ftamy9/dj_shorturl-go/internal/shortener"
)

func main() {
	cfg := config.Load()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		db, err := openDB(cfg)
		if err != nil {
			slog.Error("database connection failed", "error", err)
			os.Exit(1)
		}
		defer db.Close()
		if err := runMigrations(db, cfg); err != nil {
			slog.Error("migration failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations completed")
		return
	}

	slog.Info("connecting to database", "driver", cfg.DBDriver, "host", cfg.DBHost)
	db, err := openDB(cfg)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("database connected")

	if err := runMigrations(db, cfg); err != nil {
		slog.Error("auto-migration failed", "error", err)
		os.Exit(1)
	}

	userRepo := auth.NewRepository(db)
	addrRepo := shortener.NewRepository(db)

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

func openDB(cfg *config.Config) (*sql.DB, error) {
	if cfg.DBDriver == "sqlite3" {
		return database.OpenSQLite(cfg.DBName)
	}
	return database.OpenPostgres(cfg)
}

func runMigrations(db *sql.DB, cfg *config.Config) error {
	slog.Info("running migrations")

	isSQLite := cfg.DBDriver == "sqlite3"

	if isSQLite {
		migrations := []string{
			`CREATE TABLE IF NOT EXISTS base_users (
				id TEXT PRIMARY KEY,
				password_hash TEXT NOT NULL,
				create_date TEXT NOT NULL DEFAULT (datetime('now'))
			)`,
			`CREATE TABLE IF NOT EXISTS addresses (
				id TEXT PRIMARY KEY,
				url TEXT NOT NULL,
				created_at TEXT NOT NULL DEFAULT (datetime('now'))
			)`,
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_addresses_url ON addresses(url)`,
		}
		for _, m := range migrations {
			if _, err := db.ExecContext(context.Background(), m); err != nil {
				return fmt.Errorf("migration failed: %s: %w", truncate(m, 60), err)
			}
		}
	} else {
		migrations := []string{
			`CREATE TABLE IF NOT EXISTS base_users (
				id VARCHAR(10) PRIMARY KEY,
				password_hash VARCHAR(128) NOT NULL,
				create_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)`,
			`CREATE TABLE IF NOT EXISTS addresses (
				id UUID PRIMARY KEY,
				url VARCHAR(32779) NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)`,
			`CREATE UNIQUE INDEX IF NOT EXISTS idx_addresses_url ON addresses(url)`,
		}
		for _, m := range migrations {
			if _, err := db.ExecContext(context.Background(), m); err != nil {
				return fmt.Errorf("migration failed: %s: %w", truncate(m, 60), err)
			}
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
