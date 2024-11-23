package entities

import (
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"userId"`
	CategoryID       uuid.UUID `json:"categoryId"`
	MaximumExpending float64   `json:"maximumExpending"`
	Spending         float64   `json:"spending"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
