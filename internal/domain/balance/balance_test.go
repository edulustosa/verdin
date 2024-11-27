package balance_test

import (
	"context"
	"testing"
	"time"

	"github.com/edulustosa/verdin/internal/domain/balance"
	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/internal/domain/user"
	"github.com/google/uuid"
)

func TestCreateBalance(t *testing.T) {
	ctx := context.Background()
	t.Run("it should create a new monthly balance", func(t *testing.T) {
		userID := createUser(ctx, t)
		sut := buildSUT(t)

		_, err := sut.Create(ctx, userID)
		if err != nil {
			t.Errorf("failed to create a balance: %v", err)
		}
	})

	t.Run("it should create a balance with a previous month balance if it exists", func(t *testing.T) {
		userID := createUser(ctx, t)
		repo := new(balance.MemoryRepo)
		sut := balance.NewService(repo)

		lastBalance := entities.Balance{
			UserID:    userID,
			ID:        uuid.New(),
			Current:   1000,
			CreatedAt: time.Now().AddDate(0, -1, 0),
		}
		repo.Balances = append(repo.Balances, lastBalance)

		balance, err := sut.Create(ctx, userID)
		if err != nil {
			t.Errorf("failed to create a balance: %v", err)
		}

		if balance.Current != lastBalance.Current {
			t.Errorf("expected balance to be equal to last balance, got %v", lastBalance)
		}
	})
}

func createUser(ctx context.Context, t testing.TB) uuid.UUID {
	t.Helper()
	userRepo := new(user.MemoryRepo)

	createdUser, _ := userRepo.Create(ctx, entities.User{
		Username:     "John Doe",
		Email:        "johndoe@email.com",
		PasswordHash: "$2a$10$3Q",
	})

	return createdUser.ID
}

func buildSUT(t testing.TB) balance.Service {
	t.Helper()
	repo := new(balance.MemoryRepo)
	return balance.NewService(repo)
}
