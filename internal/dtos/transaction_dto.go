package dtos

import (
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type CreateTransaction struct {
	CategoryID  int                      `json:"categoryId" validate:"required"`
	AccountID   uuid.UUID                `json:"accountId" validate:"required"`
	Title       string                   `json:"title" validate:"required,min=3,max=255"`
	Description *string                  `json:"description,omitempty"`
	Amount      float64                  `json:"amount" validate:"required"`
	Type        entities.TransactionType `json:"type" validate:"required,oneof=INCOME EXPENSE"`
}

type GetMonthlyTransactionsQuery struct {
	Month      time.Time
	CategoryID int
}
