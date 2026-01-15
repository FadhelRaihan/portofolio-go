package user

import (
	"backend/internal/auth"
	"backend/internal/models"
	"context"
	"errors"
	"time"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service interface {
	Login(ctx context.Context, email, password string) (string, *models.User, error)

	ListUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, email, fullName, password string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, email, fullName, password *string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type service struct {
	repo        Repository
	jwtProvider auth.JWTProvider
}

func NewService(r Repository, jwtProvider auth.JWTProvider) Service {
	return &service{repo: r, jwtProvider: jwtProvider}
}

func (s *service) Login(ctx context.Context, email, password string) (string, *models.User, error) {
    log.Printf("Login attempt email=%s", email)

    u, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        log.Printf("Login: user not found for email=%s, err=%v", email, err)
        return "", nil, ErrInvalidCredentials
    }

    log.Printf("Login: found user id=%s email=%s", u.ID, u.Email)
    log.Printf("Login: incoming password='%s'", password)
    log.Printf("Login: stored hash='%s'", u.PasswordHash)

    if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
        log.Printf("Login: bcrypt mismatch: %v", err)
        return "", nil, ErrInvalidCredentials
    }

    token, err := s.jwtProvider.GenerateToken(u.ID, time.Hour)
    if err != nil {
        return "", nil, err
    }

    return token, u, nil
}


func (s *service) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.List(ctx)
}

func (s *service) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) CreateUser(ctx context.Context, email, fullName, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
	}
	return s.repo.Create(ctx, u)
}

func (s *service) UpdateUser(ctx context.Context, id string, email, fullName, password *string) (*models.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if email != nil {
		u.Email = *email
	}
	if fullName != nil {
		u.FullName = *fullName
	}
	if password != nil && *password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = string(hash)
	}

	return s.repo.Update(ctx, u)
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
