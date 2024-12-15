package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/edulustosa/verdin/internal/domain/balance"
	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/google/uuid"
)

type Service interface {
	Create(context.Context, entities.Account) (uuid.UUID, error)
	FindByID(context.Context, uuid.UUID) (*entities.Account, error)
	Update(context.Context, entities.Account) error
	GetAll(ctx context.Context, userID uuid.UUID) ([]entities.Account, error)
	Edit(context.Context, uuid.UUID, *dtos.EditAccount) error
	NewAccount(ctx context.Context, userID uuid.UUID, req *dtos.EditAccount) (uuid.UUID, error)
}

type service struct {
	repo    Repository
	balance balance.Service
}

func NewService(repo Repository, balance balance.Service) Service {
	return &service{
		repo,
		balance,
	}
}

func (s *service) Create(ctx context.Context, account entities.Account) (uuid.UUID, error) {
	return s.repo.Create(ctx, account)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Update(
	ctx context.Context,
	account entities.Account,
) error {
	return s.repo.Update(ctx, account)
}

func (s *service) GetAll(
	ctx context.Context,
	userID uuid.UUID,
) ([]entities.Account, error) {
	return s.repo.FindMany(ctx, userID)
}

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrInvalidAmount        = errors.New("invalid amount")
	ErrAccountAlreadyExists = errors.New("account already exists")
)

func (s *service) Edit(ctx context.Context, id uuid.UUID, editAccount *dtos.EditAccount) error {
	account, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrAccountNotFound
	}

	account.Title = editAccount.Title
	if editAccount.Amount < 0 {
		return ErrInvalidAmount
	}
	account.Balance = editAccount.Amount

	return s.repo.Update(ctx, *account)
}

func (s *service) NewAccount(
	ctx context.Context,
	userID uuid.UUID,
	req *dtos.EditAccount,
) (uuid.UUID, error) {
	if req.Amount < 0 {
		return uuid.Nil, ErrInvalidAmount
	}

	_, err := s.repo.FindByTitle(ctx, userID, req.Title)
	if err == nil {
		return uuid.Nil, ErrAccountAlreadyExists
	}

	if err := s.updateBalance(ctx, userID, req.Amount); err != nil {
		return uuid.Nil, fmt.Errorf("failed to update balance: %w", err)
	}

	account := entities.Account{
		UserID:  userID,
		Title:   req.Title,
		Balance: req.Amount,
	}

	return s.repo.Create(ctx, account)
}

func (s *service) updateBalance(
	ctx context.Context,
	userID uuid.UUID,
	accountAmount float64,
) error {
	balance, err := s.balance.FindByMonth(ctx, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to find balance: %w", err)
	}

	currentBalance := money.NewFromFloat(balance.Current, money.BRL)
	accountAmountMoney := money.NewFromFloat(accountAmount, money.BRL)

	currentBalance, _ = currentBalance.Add(accountAmountMoney)
	balance.Current = currentBalance.AsMajorUnits()

	return s.balance.Update(ctx, *balance)
}
