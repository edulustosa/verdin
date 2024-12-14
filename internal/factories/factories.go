package factories

import (
	"github.com/edulustosa/verdin/internal/domain/account"
	"github.com/edulustosa/verdin/internal/domain/balance"
	"github.com/edulustosa/verdin/internal/domain/category"
	"github.com/edulustosa/verdin/internal/domain/transaction"
	"github.com/edulustosa/verdin/internal/domain/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MakeUserService(db *pgxpool.Pool) user.Service {
	repo := user.NewRepo(db)
	balanceService := MakeBalanceService(db)
	accountService := MakeAccountService(db)

	return user.NewService(repo, balanceService, accountService)
}

func MakeBalanceService(db *pgxpool.Pool) balance.Service {
	repo := balance.NewRepo(db)
	return balance.NewService(repo)
}

func MakeAccountService(db *pgxpool.Pool) account.Service {
	repo := account.NewRepo(db)
	return account.NewService(repo)
}

func MakeCategoriesService(db *pgxpool.Pool) category.Service {
	repo := category.NewRepo(db)
	userRepo := user.NewRepo(db)

	return category.NewService(repo, userRepo)
}

func MakeTransactionService(db *pgxpool.Pool) transaction.Service {
	repo := transaction.NewRepo(db)
	userService := MakeUserService(db)
	categoryService := MakeCategoriesService(db)
	accountService := MakeAccountService(db)
	balanceService := MakeBalanceService(db)

	return transaction.NewService(
		repo,
		userService,
		categoryService,
		accountService,
		balanceService,
	)
}
