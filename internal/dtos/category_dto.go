package dtos

import "github.com/google/uuid"

type CreateCategory struct {
	UserID uuid.UUID `json:"userId"`
	Name   string    `json:"name"`
	Theme  string    `json:"theme"`
	Icon   string    `json:"icon"`
}
