package account

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	Create(context.Context, entities.Account) (uuid.UUID, error)
	FindByID(context.Context, uuid.UUID) (*entities.Account, error)
	Update(context.Context, entities.Account) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) Create(ctx context.Context, account entities.Account) (uuid.UUID, error) {
	return s.repo.Create(ctx, account)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(
	ctx context.Context,
	account entities.Account,
) error {
	return s.repo.Update(ctx, account)
}
