package user_test

import (
	"context"
	"testing"

	"github.com/edulustosa/verdin/internal/domain/account"
	"github.com/edulustosa/verdin/internal/domain/balance"
	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/internal/domain/user"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("it should create a new user", func(t *testing.T) {
		sut := build()

		_, err := sut.Create(ctx, entities.User{
			Username:     "edulustosa",
			Email:        "test@test.com",
			PasswordHash: "123456",
		})
		if err != nil {
			t.Errorf("failed to create user: %v", err)
		}
	})
}

func build() user.Service {
	return user.NewService(
		user.NewMemoryRepo(),
		balance.NewService(balance.NewMemoryRepo()),
		account.NewService(account.NewMemoryRepo()),
	)
}
