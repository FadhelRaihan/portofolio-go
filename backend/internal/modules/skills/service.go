package skills

import (
	"context"
	"fmt"

	"backend/internal/models"
)

type Service interface {
	Create(ctx context.Context, input *models.CreateSkillInput, imageURL string) (*models.SkillDTO, error)
	Update(ctx context.Context, id string, input *models.UpdateSkillInput, imageURL string) (*models.SkillDTO, error)
	GetByID(ctx context.Context, id string) (*models.SkillDTO, error)
	Delete(ctx context.Context, id string) (string, error) // return old image URL for deletion
	List(ctx context.Context) ([]*models.SkillDTO, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func toDTO(s *models.Skill) *models.SkillDTO {
	return &models.SkillDTO{
		ID:         s.ID,
		Name:       s.Name,
		Experience: s.Experience,
		Image:      s.Image,
	}
}

func (s *service) Create(ctx context.Context, input *models.CreateSkillInput, imageURL string) (*models.SkillDTO, error) {
	skill := &models.Skill{
		Name:       input.Name,
		Experience: input.Experience,
		Image:      imageURL,
	}
	
	created, err := s.repo.Create(ctx, skill)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create skill: %w", err)
	}
	
	return toDTO(created), nil
}

func (s *service) Update(ctx context.Context, id string, input *models.UpdateSkillInput, imageURL string) (*models.SkillDTO, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get existing skill: %w", err)
	}
	
	// Update only provided fields
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Experience != nil {
		existing.Experience = *input.Experience
	}
	if imageURL != "" {
		existing.Image = imageURL
	}
	
	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update skill: %w", err)
	}
	
	return toDTO(updated), nil
}

func (s *service) GetByID(ctx context.Context, id string) (*models.SkillDTO, error) {
	skill, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get skill: %w", err)
	}
	return toDTO(skill), nil
}

func (s *service) Delete(ctx context.Context, id string) (string, error) {
	// Get skill first to return image URL for deletion
	skill, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("service: failed to get skill for deletion: %w", err)
	}
	
	if err := s.repo.Delete(ctx, id); err != nil {
		return "", fmt.Errorf("service: failed to delete skill: %w", err)
	}
	
	return skill.Image, nil
}

func (s *service) List(ctx context.Context) ([]*models.SkillDTO, error) {
	skills, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list skills: %w", err)
	}
	
	dtos := make([]*models.SkillDTO, len(skills))
	for i, skill := range skills {
		dtos[i] = toDTO(skill)
	}
	
	return dtos, nil
}