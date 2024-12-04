package transaction_test

import "testing"

func TestTransaction(t *testing.T) {
	// ctx := context.Background()
	t.Run("it should create a new transaction", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("it should update a balance when a new transaction is created", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("it should return an error when the user does not have enough balance", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("it should update a account balance when a new transaction is created", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("it should return an error when the account does not have enough balance", func(t *testing.T) {
		t.Skip("not implemented")
	})
}

/*
func build(t testing.TB) transaction.Service {
	t.Helper()

	balanceService := balance.NewService(balance.NewMemoryRepo())
	accountService := account.NewService(account.NewMemoryRepo())
	userService := user.NewService(
		user.NewMemoryRepo(),
		balanceService,
		accountService,
	)
	categoryService := category.NewService(category.NewMemoryRepo(), userService)

	return transaction.NewService(
		transaction.NewMemoryRepo(),
		userService,
		categoryService,
		accountService,
		balanceService,
	)
}
*/
