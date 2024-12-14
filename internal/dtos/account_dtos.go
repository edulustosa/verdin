package dtos

type EditAccount struct {
	Title  string  `json:"title" validate:"required,min=3,max=50"`
	Amount float64 `json:"amount" validate:"gte=0"`
}
