package balance

import (
	"context"
	"errors"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type MemoryRepo struct {
	Balances []entities.Balance
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
