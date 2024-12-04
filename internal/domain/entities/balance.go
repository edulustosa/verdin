package entities

import (
	"errors"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/edulustosa/verdin/pkg/utils"
	"github.com/google/uuid"
)

type Balance struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Current   float64   `json:"current"`
	Income    float64   `json:"income"`
	Expenses  float64   `json:"expenses"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Balance) Update(transaction *Transaction) error {
	currentBalance := money.NewFromFloat(b.Current, money.BRL)
	transactionAmount := money.NewFromFloat(transaction.Amount, money.BRL)

	if transaction.Type == Expense {
		hasInsufficientFunds, _ := currentBalance.LessThan(transactionAmount)
		if hasInsufficientFunds {
			return errors.New("insufficient funds")
		}

		b.Expenses = utils.SumMoney(b.Expenses, transactionAmount)
		transactionAmount = transactionAmount.Negative()
	} else {
		b.Income = utils.SumMoney(b.Income, transactionAmount)
	}

	currentBalance, _ = currentBalance.Add(transactionAmount)
	b.Current = currentBalance.AsMajorUnits()

	return nil
}
