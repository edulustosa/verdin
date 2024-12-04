package category

import (
	"context"
	"errors"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, entities.Category) (*entities.Category, error)
	FindByName(ctx context.Context, userID uuid.UUID, name string) (*entities.Category, error)
	FindByID(context.Context, uuid.UUID) (*entities.Category, error)
}

type memoryRepo struct {
	categories []entities.Category
}

func NewMemoryRepo() Repository {
	return &memoryRepo{}
}

func (r *memoryRepo) Create(_ context.Context, category entities.Category) (*entities.Category, error) {
	category.ID = uuid.New()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	r.categories = append(r.categories, category)
	return &category, nil
}

func (r *memoryRepo) FindByName(_ context.Context, userID uuid.UUID, name string) (*entities.Category, error) {
	for _, c := range r.categories {
		if c.UserID == userID && c.Name == name {
			return &c, nil
		}
	}
	return nil, errors.New("category not found")
}

func (r *memoryRepo) FindByID(_ context.Context, id uuid.UUID) (*entities.Category, error) {
	for _, c := range r.categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, errors.New("category not found")
}
