package entities

import (
	"time"

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
