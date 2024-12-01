package entities

import "github.com/google/uuid"

type Category struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userId"`
	Name   string    `json:"name"`
	Theme  string    `json:"theme"`
	Icon   string    `json:"icon"`
}
