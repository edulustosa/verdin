package entities

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         int             `json:"id"`
	UserID     uuid.UUID       `json:"userId"`
	CategoryID uuid.UUID       `json:"categoryId"`
	AccountID  uuid.UUID       `json:"accountId"`
	Title      string          `json:"title"`
	Amount     float64         `json:"amount"`
	Type       TransactionType `json:"type"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

type TransactionType string

const (
	Income  TransactionType = "INCOME"
	Expense TransactionType = "EXPENSE"
)
