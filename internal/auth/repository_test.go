package auth

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	f, err := os.CreateTemp("", "test-auth-*.db")
	if err != nil {
		t.Fatalf("create temp db: %v", err)
	}
	f.Close()
	db, err := sql.Open("sqlite3", f.Name()+"?_journal_mode=WAL")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	_, err = db.ExecContext(context.Background(),
		`CREATE TABLE base_users (
			id TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL,
			create_date TEXT NOT NULL DEFAULT (datetime('now'))
		)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	t.Cleanup(func() { db.Close(); os.Remove(f.Name()); os.Remove(f.Name() + "-shm"); os.Remove(f.Name() + "-wal") })
	return db
}

func TestUserSignupAndLogin(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	user := &BaseUser{ID: "alice", PasswordHash: "hashed_value"}
	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if user.CreateDate == "" {
		t.Fatal("CreateDate not set")
	}

	fetched, err := repo.GetByID(context.Background(), "alice")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched.ID != "alice" {
		t.Fatalf("expected id alice, got %s", fetched.ID)
	}
	if fetched.PasswordHash != "hashed_value" {
		t.Fatalf("expected hash hashed_value, got %s", fetched.PasswordHash)
	}
}
