package user

import (
	"backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, u *models.User) (*models.User, error)
	Update(ctx context.Context, u *models.User) (*models.User, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT id, email, password, full_name, created_at
		 FROM users
		 WHERE email = $1`,
		email,
	)
	return scanUser(row)
}

func (r *repository) FindByID(ctx context.Context, id string) (*models.User, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT id, email, password, full_name, created_at
		 FROM users
		 WHERE id = $1`,
		id,
	)
	return scanUser(row)
}

func (r *repository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id, email, password, full_name, created_at
		 FROM users
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

func (r *repository) Create(ctx context.Context, u *models.User) (*models.User, error) {
	row := r.db.QueryRow(
		ctx,
		`INSERT INTO users (id, email, password, full_name)
		 VALUES (gen_random_uuid(), $1, $2, $3)
		 RETURNING id, created_at`,
		u.Email, u.PasswordHash, u.FullName,
	)
	if err := row.Scan(&u.ID, &u.CreatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) Update(ctx context.Context, u *models.User) (*models.User, error) {
	row := r.db.QueryRow(
		ctx,
		`UPDATE users
		 SET email = $1,
		     password = $2,
		     full_name = $3
		 WHERE id = $4
		 RETURNING created_at`,
		u.Email, u.PasswordHash, u.FullName, u.ID,
	)
	if err := row.Scan(&u.CreatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(
		ctx,
		`DELETE FROM users WHERE id = $1`,
		id,
	)
	return err
}

func scanUser(row pgx.Row) (*models.User, error) {
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
