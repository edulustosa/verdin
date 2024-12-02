package category

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, entities.Category) (*entities.Category, error)
	FindByName(ctx context.Context, userID uuid.UUID, name string) (*entities.Category, error)
	FindByID(context.Context, uuid.UUID) (*entities.Category, error)
}
