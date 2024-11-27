package user

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user entities.User) (uuid.UUID, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *service) Create(ctx context.Context, user entities.User) (uuid.UUID, error) {
	return s.repo.Create(ctx, user)
}
