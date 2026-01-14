package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	// token dummy; nanti ganti JWT
	return "dummy-token", nil
}
