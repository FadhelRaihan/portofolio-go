package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) FindByEmail(ctx context.Context, email string) (*User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, email, password_hash FROM users WHERE email = $1", email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
		return nil, err
	}

	return &u, nil
}
