package transaction

import (
	"context"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
)

type Repository interface {
	Create(context.Context, entities.Transaction) (*entities.Transaction, error)
}

type memory struct {
	transactions []entities.Transaction
}

func NewMemoryRepo() Repository {
	return &memory{}
}

func (r *memory) Create(
	_ context.Context,
	t entities.Transaction,
) (*entities.Transaction, error) {
	t.ID = len(r.transactions) + 1
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	r.transactions = append(r.transactions, t)
	return &t, nil
}
