package auth

import (
	"context"
	"errors"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/internal/domain/user"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	user user.Service
}

func New(user user.Service) *Auth {
	return &Auth{
		user,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

func (a *Auth) Login(ctx context.Context, login *dtos.Login) (uuid.UUID, error) {
	user, err := a.user.FindByEmail(ctx, login.Email)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(login.Password),
	)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	return user.ID, nil
}

func (a *Auth) Register(
	ctx context.Context,
	register *dtos.Register,
) (uuid.UUID, error) {
	_, err := a.user.FindByEmail(ctx, register.Email)
	if err == nil {
		return uuid.Nil, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(register.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return uuid.Nil, err
	}

	user := entities.User{
		Username:     register.Username,
		Email:        register.Email,
		PasswordHash: string(passwordHash),
	}

	return a.user.Create(ctx, user)
}
