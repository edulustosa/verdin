package balance

import (
	"context"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	FindByID(context.Context, uuid.UUID) (*entities.Balance, error)
	FindByMonth(
		ctx context.Context,
		userID uuid.UUID,
		month time.Time,
	) (*entities.Balance, error)
	Update(ctx context.Context, balance entities.Balance) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	balance := entities.Balance{
		UserID: userID,
	}

	lastMonth := time.Now().AddDate(0, -1, 0)
	lastBalance, err := s.repo.FindByMonth(ctx, userID, lastMonth)
	if err == nil {
		balance.Current = lastBalance.Current
	}

	return s.repo.Create(ctx, balance)
}

func (s *service) FindByMonth(
	ctx context.Context,
	userID uuid.UUID,
	month time.Time,
) (*entities.Balance, error) {
	return s.repo.FindByMonth(ctx, userID, month)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*entities.Balance, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, balance entities.Balance) error {
	return s.repo.Update(ctx, balance)
}
