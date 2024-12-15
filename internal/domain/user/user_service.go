package user

import (
	"context"
	"fmt"

	"github.com/edulustosa/verdin/internal/domain/account"
	"github.com/edulustosa/verdin/internal/domain/balance"
	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Service interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user entities.User) (uuid.UUID, error)
	FindByID(context.Context, uuid.UUID) (*entities.User, error)
}

type service struct {
	repo    Repository
	balance balance.Service
	account account.Service
}

func NewService(
	repo Repository,
	balance balance.Service,
	account account.Service,
) Service {
	return &service{
		repo,
		balance,
		account,
	}
}

func (s *service) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) Create(ctx context.Context, user entities.User) (uuid.UUID, error) {
	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	_, err = s.balance.Create(ctx, userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create balance: %w", err)
	}

	_, err = s.account.Create(ctx, entities.Account{
		UserID: userID,
		Title:  "Carteira",
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create account: %w", err)
	}

	return userID, nil
}
