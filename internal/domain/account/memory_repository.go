package account

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
)

type MemoryRepo struct {
	Accounts []entities.Account
}

func (r *MemoryRepo) Create(_ context.Context, account entities.Account) (*entities.Account, error) {
	r.Accounts = append(r.Accounts, account)
	return &account, nil
}
