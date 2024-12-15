package balance

import (
	"context"
	"errors"
	"fmt"
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

var ErrMonthNotFound = errors.New("month not found")

func (s *service) FindByMonth(
	ctx context.Context,
	userID uuid.UUID,
	month time.Time,
) (*entities.Balance, error) {
	balance, err := s.repo.FindByMonth(ctx, userID, month)
	if err == nil {
		return balance, nil
	}

	if isSameMonth(month, time.Now()) {
		createdBalanceID, err := s.repo.Create(ctx, entities.Balance{
			UserID:   userID,
			Income:   0,
			Expenses: 0,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create balance: %w", err)
		}

		return s.repo.FindByID(ctx, createdBalanceID)
	}

	return nil, ErrMonthNotFound
}

func isSameMonth(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month()
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*entities.Balance, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, balance entities.Balance) error {
	return s.repo.Update(ctx, balance)
}
