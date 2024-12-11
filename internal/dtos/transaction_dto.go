package dtos

import (
	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type CreateTransaction struct {
	UserID     uuid.UUID                `json:"userId"`
	CategoryID int                      `json:"categoryId"`
	AccountID  uuid.UUID                `json:"accountId"`
	Title      string                   `json:"title"`
	Amount     float64                  `json:"amount"`
	Type       entities.TransactionType `json:"type"`
}
