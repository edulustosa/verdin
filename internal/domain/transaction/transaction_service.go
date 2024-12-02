package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/edulustosa/verdin/internal/domain/account"
	"github.com/edulustosa/verdin/internal/domain/balance"
	"github.com/edulustosa/verdin/internal/domain/category"
	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/internal/domain/user"
	"github.com/edulustosa/verdin/internal/dtos"
)

type Service interface {
}

type service struct {
	repo     Repository
	user     user.Service
	category category.Service
	account  account.Service
	balance  balance.Service
}

func NewService(
	repo Repository,
	user user.Service,
	category category.Service,
	account account.Service,
	balance balance.Service,
) Service {
	return &service{
		repo,
		user,
		category,
		account,
		balance,
	}
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrAccountNotFound  = errors.New("account not found")
	ErrCategoryNotFound = errors.New("category not found")
)

func (s *service) CreateTransaction(
	ctx context.Context,
	transaction *dtos.CreateTransaction,
) (*entities.Transaction, error) {
	user, err := s.user.FindByID(ctx, transaction.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	_, err = s.category.FindByID(ctx, transaction.CategoryID)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	account, err := s.account.FindByID(ctx, transaction.AccountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	balance, err := s.balance.FindByMonth(
		ctx,
		user.ID,
		time.Now().Month(),
	)
	if err != nil {
		return nil, err
	}
}
