package user

import (
	"context"
	"errors"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user entities.User) (uuid.UUID, error)
	FindByID(context.Context, uuid.UUID) (*entities.User, error)
}

type memoryRepository struct {
	users []entities.User
}

func NewMemoryRepo() Repository {
	return &memoryRepository{}
}

func (r *memoryRepository) FindByEmail(_ context.Context, email string) (*entities.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *memoryRepository) Create(_ context.Context, user entities.User) (uuid.UUID, error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.users = append(r.users, user)
	return user.ID, nil
}

func (r *memoryRepository) FindByID(_ context.Context, id uuid.UUID) (*entities.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}
