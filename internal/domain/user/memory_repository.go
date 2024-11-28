package user

import (
	"context"
	"errors"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type MemoryRepo struct {
	Users []entities.User
}

func (r *MemoryRepo) Create(_ context.Context, user entities.User) (uuid.UUID, error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.Users = append(r.Users, user)
	return user.ID, nil
}

func (r *MemoryRepo) FindByEmail(_ context.Context, email string) (*entities.User, error) {
	for _, u := range r.Users {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, errors.New("user not found")
}
