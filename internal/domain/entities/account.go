package entities

import (
	"errors"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Title     string    `json:"title"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Account) Update(transaction *Transaction) error {
	balanceMoney := money.NewFromFloat(a.Balance, money.BRL)
	transactionMoney := money.NewFromFloat(transaction.Amount, money.BRL)

	if transaction.Type == Expense {
		hasInsufficientFunds, _ := balanceMoney.LessThan(transactionMoney)
		if hasInsufficientFunds {
			return errors.New("insufficient funds")
		}

		transactionMoney = transactionMoney.Negative()
	}

	balanceMoney, _ = balanceMoney.Add(transactionMoney)
	a.Balance = balanceMoney.AsMajorUnits()

	return nil
}
