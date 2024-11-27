package user

import (
	"context"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
)

type MemoryRepo struct {
	Users []entities.User
}

func (r *MemoryRepo) Create(_ context.Context, user entities.User) (*entities.User, error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.Users = append(r.Users, user)
	return &user, nil
}
