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
	CreateTransaction(context.Context, *dtos.CreateTransaction) (*entities.Transaction, error)
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
	ErrUserNotFound      = errors.New("user not found")
	ErrAccountNotFound   = errors.New("account not found")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

func (s *service) CreateTransaction(
	ctx context.Context,
	transactionReq *dtos.CreateTransaction,
) (*entities.Transaction, error) {
	user, err := s.user.FindByID(ctx, transactionReq.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	category, err := s.category.FindByID(ctx, transactionReq.CategoryID)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	account, err := s.account.FindByID(ctx, transactionReq.AccountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	lastMonth := time.Now().AddDate(0, -1, 0)
	balance, err := s.balance.FindByMonth(
		ctx,
		user.ID,
		lastMonth,
	)
	if err != nil {
		return nil, err
	}

	transaction := entities.Transaction{
		UserID:     user.ID,
		CategoryID: category.ID,
		AccountID:  account.ID,
		BalanceID:  balance.ID,
		Title:      transactionReq.Title,
		Amount:     transactionReq.Amount,
		Type:       transactionReq.Type,
	}

	if err := s.updateAccount(ctx, account, &transaction); err != nil {
		return nil, err
	}

	if err := s.updateBalance(ctx, balance, &transaction); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, transaction)
}

func (s *service) updateBalance(
	ctx context.Context,
	balance *entities.Balance,
	transaction *entities.Transaction,
) error {
	if err := balance.Update(transaction); err != nil {
		return ErrInsufficientFunds
	}

	return s.balance.Update(ctx, *balance)
}

func (s *service) updateAccount(
	ctx context.Context,
	account *entities.Account,
	transaction *entities.Transaction,
) error {
	if err := account.Update(transaction); err != nil {
		return ErrInsufficientFunds
	}

	return s.account.Update(ctx, *account)
}
