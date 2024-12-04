package account

import (
	"context"
	"errors"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, entities.Account) (*entities.Account, error)
	FindByID(context.Context, uuid.UUID) (*entities.Account, error)
	Update(context.Context, entities.Account) (*entities.Account, error)
}

type memoryRepo struct {
	accounts []entities.Account
}

func NewMemoryRepo() Repository {
	return &memoryRepo{}
}

func (r *memoryRepo) Create(_ context.Context, account entities.Account) (*entities.Account, error) {
	account.ID = uuid.New()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	r.accounts = append(r.accounts, account)
	return &account, nil
}

func (r *memoryRepo) FindByID(_ context.Context, id uuid.UUID) (*entities.Account, error) {
	for _, account := range r.accounts {
		if account.ID == id {
			return &account, nil
		}
	}
	return nil, errors.New("account not found")
}

func (r *memoryRepo) Update(_ context.Context, account entities.Account) (*entities.Account, error) {
	for i, a := range r.accounts {
		if a.ID == account.ID {
			account.UpdatedAt = time.Now()
			r.accounts[i] = account
			return &account, nil
		}
	}
	return nil, errors.New("account not found")
}
