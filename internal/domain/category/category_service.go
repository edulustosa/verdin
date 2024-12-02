package category

import (
	"context"
	"errors"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/internal/domain/user"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/google/uuid"
)

type Service interface {
	Create(context.Context, *dtos.CreateCategory) (*entities.Category, error)
	FindByID(context.Context, uuid.UUID) (*entities.Category, error)
}

type service struct {
	repo Repository
	user user.Service
}

func NewService(repo Repository, user user.Service) Service {
	return &service{
		repo,
		user,
	}
}

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
)

func (s *service) Create(
	ctx context.Context,
	createCategoryReq *dtos.CreateCategory,
) (*entities.Category, error) {
	_, err := s.user.FindByID(ctx, createCategoryReq.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	_, err = s.repo.FindByName(
		ctx,
		createCategoryReq.UserID,
		createCategoryReq.Name,
	)
	if err == nil {
		return nil, ErrCategoryAlreadyExists
	}

	category := entities.Category{
		UserID: createCategoryReq.UserID,
		Name:   createCategoryReq.Name,
		Theme:  createCategoryReq.Theme,
		Icon:   createCategoryReq.Icon,
	}

	return s.repo.Create(ctx, category)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	return s.repo.FindByID(ctx, id)
}
