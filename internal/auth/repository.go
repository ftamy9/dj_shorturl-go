package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *BaseUser) error
	GetByID(ctx context.Context, id string) (*BaseUser, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *BaseUser) error {
	user.CreateDate = time.Now()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO base_users (id, password_hash, create_date) VALUES ($1, $2, $3)`,
		user.ID, user.PasswordHash, user.CreateDate,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*BaseUser, error) {
	user := &BaseUser{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, password_hash, create_date FROM base_users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.PasswordHash, &user.CreateDate)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
