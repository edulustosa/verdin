package balance

import (
	"context"
	"errors"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, *entities.Balance) (*entities.Balance, error)
	FindByID(context.Context, uuid.UUID) (*entities.Balance, error)
	FindByMonth(
		ctx context.Context,
		userID uuid.UUID,
		month time.Month,
	) (*entities.Balance, error)
	Update(context.Context, entities.Balance) (*entities.Balance, error)
}

type MemoryRepo struct {
	Balances []entities.Balance
}

func NewMemoryRepo() Repository {
	return &MemoryRepo{}
}

func (r *MemoryRepo) Create(_ context.Context, balance *entities.Balance) (*entities.Balance, error) {
	balance.ID = uuid.New()
	balance.CreatedAt = time.Now()
	balance.UpdatedAt = time.Now()

	r.Balances = append(r.Balances, *balance)
	return balance, nil
}

func (r *MemoryRepo) FindByMonth(
	_ context.Context,
	userID uuid.UUID,
	month time.Month,
) (*entities.Balance, error) {
	for _, b := range r.Balances {
		if b.UserID == userID {
			if b.CreatedAt.Month() == month && b.CreatedAt.Year() == time.Now().Year() {
				return &b, nil
			}
		}
	}
	return nil, errors.New("balance not found")
}

func (r *MemoryRepo) FindByID(_ context.Context, id uuid.UUID) (*entities.Balance, error) {
	for _, b := range r.Balances {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, errors.New("balance not found")
}

func (r *MemoryRepo) Update(_ context.Context, balance entities.Balance) (*entities.Balance, error) {
	for i, b := range r.Balances {
		if b.ID == balance.ID {
			balance.UpdatedAt = time.Now()
			r.Balances[i] = balance
			return &balance, nil
		}
	}
	return nil, errors.New("balance not found")
}
