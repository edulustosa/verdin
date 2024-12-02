package balance

import (
	"context"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/pkg/utils"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID) (*entities.Balance, error)
	FindByID(context.Context, uuid.UUID) (*entities.Balance, error)
	FindByMonth(
		ctx context.Context,
		userID uuid.UUID,
		month time.Month,
	) (*entities.Balance, error)
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

func (s *service) FindByMonth(
	ctx context.Context,
	userID uuid.UUID,
	month time.Month,
) (*entities.Balance, error) {
	return s.repo.FindByMonth(ctx, userID, month)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*entities.Balance, error) {
	return s.repo.FindByID(ctx, id)
}
