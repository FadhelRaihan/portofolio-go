package profile

import (
	"backend/internal/models"
	"context"
)

type Service interface {
	GetCurrentProfile(ctx context.Context) (*models.ProfileDTO, error)
	GetByID(ctx context.Context, id string) (*models.ProfileDTO, error)
	Create(ctx context.Context, input CreateProfileInput) (*models.ProfileDTO, error)
	Update(ctx context.Context, id string, input UpdateProfileInput) (*models.ProfileDTO, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

type CreateProfileInput struct {
	FullName    *string `json:"full_name"`
	Headline    *string `json:"headline"`
	Bio         *string `json:"bio"`
	Location    *string `json:"location"`
	GithubURL   *string `json:"github_url"`
	LinkedinURL *string `json:"linkedin_url"`
}

type UpdateProfileInput = CreateProfileInput

func toDTO(p *models.Profile) *models.ProfileDTO {
	return &models.ProfileDTO{
		ID:          p.ID,
		FullName:    p.FullName,
		Headline:    p.Headline,
		Bio:         p.Bio,
		Location:    p.Location,
		GithubURL:   p.GithubURL,
		LinkedinURL: p.LinkedinURL,
	}
}

func (s *service) GetCurrentProfile(ctx context.Context) (*models.ProfileDTO, error) {
	p, err := s.repo.GetLatest(ctx)
	if err != nil {
		return nil, err
	}
	return toDTO(p), nil
}

func (s *service) GetByID(ctx context.Context, id string) (*models.ProfileDTO, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDTO(p), nil
}

func (s *service) Create(ctx context.Context, input CreateProfileInput) (*models.ProfileDTO, error) {
	p := &models.Profile{
		FullName:    *input.FullName,
		Headline:    *input.Headline,
		Bio:         *input.Bio,
		Location:    *input.Location,
		GithubURL:   *input.GithubURL,
		LinkedinURL: *input.LinkedinURL,
	}
	p, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}
	return toDTO(p), nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateProfileInput) (*models.ProfileDTO, error) {
    // 1. ambil data lama
    existing, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 2. merge field
    if input.FullName != nil {
        existing.FullName = *input.FullName
    }
    if input.Headline != nil {
        existing.Headline = *input.Headline
    }
    if input.Bio != nil {
        existing.Bio = *input.Bio
    }
    if input.Location != nil {
        existing.Location = *input.Location
    }
    if input.GithubURL != nil {
        existing.GithubURL = *input.GithubURL
    }
    if input.LinkedinURL != nil {
        existing.LinkedinURL = *input.LinkedinURL
    }

    // 3. simpan
    updated, err := s.repo.Update(ctx, existing)
    if err != nil {
        return nil, err
    }
    return toDTO(updated), nil
}


func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
