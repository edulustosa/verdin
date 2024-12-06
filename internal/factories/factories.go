package factories

import (
	"github.com/edulustosa/verdin/internal/domain/account"
	"github.com/edulustosa/verdin/internal/domain/balance"
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
