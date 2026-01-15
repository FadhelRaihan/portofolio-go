package profile

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"backend/internal/models"
)

type Repository interface {
	GetLatest(ctx context.Context) (*models.Profile, error)
	GetByID(ctx context.Context, id string) (*models.Profile, error)
	Create(ctx context.Context, p *models.Profile) (*models.Profile, error)
	Update(ctx context.Context, p *models.Profile) (*models.Profile, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetLatest(ctx context.Context) (*models.Profile, error) {
	const q = `
		select id, full_name, headline, bio, location,
		       github_url, linkedin_url,
		       created_at, updated_at
		from profiles
		order by created_at desc
		limit 1;
	`
	row := r.db.QueryRow(ctx, q)

	var p models.Profile
	if err := row.Scan(
		&p.ID, &p.FullName, &p.Headline, &p.Bio, &p.Location,
		&p.GithubURL, &p.LinkedinURL,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*models.Profile, error) {
	const q = `
		select id, full_name, headline, bio, location,
		       github_url, linkedin_url,
		       created_at, updated_at
		from profiles
		where id = $1;
	`
	row := r.db.QueryRow(ctx, q, id)

	var p models.Profile
	if err := row.Scan(
		&p.ID, &p.FullName, &p.Headline, &p.Bio, &p.Location,
		&p.GithubURL, &p.LinkedinURL,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) Create(ctx context.Context, p *models.Profile) (*models.Profile, error) {
	const q = `
		insert into profiles (
			full_name, headline, bio, location,
			github_url, linkedin_url
		)
		values ($1,$2,$3,$4,$5,$6)
		returning id, created_at, updated_at;
	`
	err := r.db.QueryRow(ctx, q,
		p.FullName, p.Headline, p.Bio, p.Location,
		p.GithubURL, p.LinkedinURL,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *repository) Update(ctx context.Context, p *models.Profile) (*models.Profile, error) {
	const q = `
		update profiles
		set full_name = $1,
		    headline  = $2,
		    bio       = $3,
		    location  = $4,
		    github_url   = $5,
		    linkedin_url = $6,

		    updated_at   = now()
		where id = $7
		returning created_at, updated_at;
	`
	err := r.db.QueryRow(ctx, q,
		p.FullName, p.Headline, p.Bio, p.Location,
		p.GithubURL, p.LinkedinURL,
		p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	const q = `delete from profiles where id = $1;`
	cmd, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("profile not found")
	}
	return nil
}
