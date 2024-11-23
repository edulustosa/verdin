package entities

import "github.com/google/uuid"

type Category struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Theme string    `json:"theme"`
	Icon  string    `json:"icon"`
}