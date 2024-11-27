package balance

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/pkg/utils"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID) (*entities.Balance, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) Create(
	ctx context.Context,
	userID uuid.UUID,
) (*entities.Balance, error) {
	balance := entities.Balance{
		UserID: userID,
	}

	lastBalance, err := s.repo.FindByMonth(ctx, userID, utils.GetLastMonth())
	if err == nil {
		balance.Current = lastBalance.Current
	}

	return s.repo.Create(ctx, &balance)
}
