package shortener

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, addr *Address) error
	GetByID(ctx context.Context, id string) (*Address, error)
	GetByURL(ctx context.Context, url string) (*Address, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, addr *Address) error {
	addr.ID = uuid.New().String()
	addr.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO addresses (id, url, created_at) VALUES ($1, $2, $3)`,
		addr.ID, addr.URL, addr.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create address: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Address, error) {
	addr := &Address{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, url, created_at FROM addresses WHERE id = $1`,
		id,
	).Scan(&addr.ID, &addr.URL, &addr.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get address: %w", err)
	}
	return addr, nil
}

func (r *PostgresRepository) GetByURL(ctx context.Context, url string) (*Address, error) {
	addr := &Address{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, url, created_at FROM addresses WHERE url = $1`,
		url,
	).Scan(&addr.ID, &addr.URL, &addr.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get address by url: %w", err)
	}
	return addr, nil
}
