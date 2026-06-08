package shortener

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, addr *Address) error
	GetByID(ctx context.Context, id string) (*Address, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, addr *Address) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO addresses (url, created_at) VALUES ($1, NOW()) RETURNING id, created_at`,
		addr.URL,
	).Scan(&addr.ID, &addr.CreatedAt)
	if err != nil {
		return fmt.Errorf("create address: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Address, error) {
	addr := &Address{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, url, created_at FROM addresses WHERE id = $1`,
		id,
	).Scan(&addr.ID, &addr.URL, &addr.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get address: %w", err)
	}
	return addr, nil
}
