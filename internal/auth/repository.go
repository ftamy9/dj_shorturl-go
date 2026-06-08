package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *BaseUser) error
	GetByID(ctx context.Context, id string) (*BaseUser, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, user *BaseUser) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO base_users (id, password_hash, create_date) VALUES ($1, $2, NOW())`,
		user.ID, user.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*BaseUser, error) {
	user := &BaseUser{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, password_hash, create_date FROM base_users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.PasswordHash, &user.CreateDate)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
