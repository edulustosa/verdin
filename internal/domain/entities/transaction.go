package entities

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Income  TransactionType = "INCOME"
	Expense TransactionType = "EXPENSE"
)

type Transaction struct {
	ID          int             `json:"id"`
	UserID      uuid.UUID       `json:"userId"`
	CategoryID  int             `json:"categoryId"`
	AccountID   uuid.UUID       `json:"accountId"`
	BalanceID   uuid.UUID       `json:"balanceId"`
	Title       string          `json:"title"`
	Description *string         `json:"description,omitempty"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
