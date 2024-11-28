package account

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
)

type Repository interface {
	Create(context.Context, entities.Account) (*entities.Account, error)
}
