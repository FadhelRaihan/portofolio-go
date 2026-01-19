package skills

import (
	"context"
	"fmt"

	"backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.Skill, error)
	Create(ctx context.Context, s *models.Skill) (*models.Skill, error)
	Update(ctx context.Context, s *models.Skill) (*models.Skill, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*models.Skill, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*models.Skill, error) {
	const q = `SELECT id, name, experience, image, created_at, updated_at 
	           FROM skills WHERE id = $1`
	
	row := r.db.QueryRow(ctx, q, id)
	var s models.Skill
	err := row.Scan(&s.ID, &s.Name, &s.Experience, &s.Image, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("skill with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}
	return &s, nil
}

func (r *repository) Create(ctx context.Context, s *models.Skill) (*models.Skill, error) {
	const q = `INSERT INTO skills (name, experience, image) 
	           VALUES ($1, $2, $3) 
	           RETURNING id, created_at, updated_at`
	
	err := r.db.QueryRow(ctx, q, s.Name, s.Experience, s.Image).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create skill: %w", err)
	}
	return s, nil
}

func (r *repository) Update(ctx context.Context, s *models.Skill) (*models.Skill, error) {
	const q = `UPDATE skills 
	           SET name = $1, experience = $2, image = $3, updated_at = NOW() 
	           WHERE id = $4 
	           RETURNING updated_at`
	
	err := r.db.QueryRow(ctx, q, s.Name, s.Experience, s.Image, s.ID).
		Scan(&s.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("skill with id %s not found", s.ID)
		}
		return nil, fmt.Errorf("failed to update skill: %w", err)
	}
	return s, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM skills WHERE id = $1`
	
	cmd, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("skill with id %s not found", id)
	}
	return nil
}

func (r *repository) List(ctx context.Context) ([]*models.Skill, error) {
	const q = `SELECT id, name, experience, image, created_at, updated_at 
	           FROM skills 
	           ORDER BY created_at DESC`
	
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}
	defer rows.Close()
	
	var skills []*models.Skill
	for rows.Next() {
		var s models.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Experience, &s.Image, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		skills = append(skills, &s)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating skills: %w", err)
	}
	
	return skills, nil
}