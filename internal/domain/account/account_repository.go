package account

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, entities.Account) (*entities.Account, error)
	FindByID(context.Context, uuid.UUID) (*entities.Account, error)
}
