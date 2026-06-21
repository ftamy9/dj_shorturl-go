package shortener

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	f, err := os.CreateTemp("", "test-short-*.db")
	if err != nil {
		t.Fatalf("create temp db: %v", err)
	}
	f.Close()
	db, err := sql.Open("sqlite3", f.Name()+"?_journal_mode=WAL")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	_, err = db.ExecContext(context.Background(),
		`CREATE TABLE addresses (
			id TEXT PRIMARY KEY,
			url TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now'))
		)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	_, err = db.ExecContext(context.Background(),
		`CREATE UNIQUE INDEX idx_addresses_url ON addresses(url)`)
	if err != nil {
		t.Fatalf("create index: %v", err)
	}
	t.Cleanup(func() { db.Close(); os.Remove(f.Name()); os.Remove(f.Name() + "-shm"); os.Remove(f.Name() + "-wal") })
	return db
}

func TestShortURLCreateAndDedup(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	addr := &Address{URL: "https://example.com"}
	if err := repo.Create(context.Background(), addr); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if addr.ID == "" {
		t.Fatal("ID not set")
	}
	if _, err := uuid.Parse(addr.ID); err != nil {
		t.Fatalf("ID is not a valid UUID: %v", err)
	}
	if addr.CreatedAt == "" {
		t.Fatal("CreatedAt not set")
	}

	fetched, err := repo.GetByID(context.Background(), addr.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.URL != "https://example.com" {
		t.Fatalf("expected url https://example.com, got %s", fetched.URL)
	}

	existing, err := repo.GetByURL(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("GetByURL failed: %v", err)
	}
	if existing.ID != addr.ID {
		t.Fatal("GetByURL returned different ID")
	}
}
