package dtos

type CreateCategory struct {
	Name  string `json:"name" validate:"required,min=3,max=255"`
	Theme string `json:"theme" validate:"required,min=3,max=255"`
	Icon  string `json:"icon" validate:"required,min=3,max=255"`
}

type UpdateCategory struct {
	Name  string `json:"name" validate:"required,min=3,max=255"`
	Theme string `json:"theme" validate:"required,min=3,max=255"`
	Icon  string `json:"icon" validate:"required,min=3,max=255"`
}
