package entities

import "github.com/google/uuid"

type Pot struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	CategoryID uuid.UUID `json:"categoryId"`
	Name       string    `json:"name"`
	Total      float64   `json:"total"`
	Target     float64   `json:"target"`
	Theme      string    `json:"theme"`
}
