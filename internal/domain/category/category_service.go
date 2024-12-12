package category

import (
	"context"
	"errors"
	"sync"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/edulustosa/verdin/internal/domain/user"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req *dtos.CreateCategory) (int, error)
	FindByID(context.Context, int) (*entities.Category, error)
	CreateDefaultCategories(ctx context.Context, userID uuid.UUID) error
	Update(
		ctx context.Context,
		id int,
		userID uuid.UUID,
		req *dtos.UpdateCategory,
	) error
	GetAll(ctx context.Context, userID uuid.UUID) ([]entities.Category, error)
}

type service struct {
	repo Repository
	user user.Repository
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
	ErrCategoryNotFound      = errors.New("category not found")
)

func (s *service) Create(
	ctx context.Context,
	userID uuid.UUID,
	createCategoryReq *dtos.CreateCategory,
) (int, error) {
	_, err := s.user.FindByID(ctx, userID)
	if err != nil {
		return 0, ErrUserNotFound
	}

	_, err = s.repo.FindByName(
		ctx,
		userID,
		createCategoryReq.Name,
	)
	if err == nil {
		return 0, ErrCategoryAlreadyExists
	}

	category := entities.Category{
		UserID: userID,
		Name:   createCategoryReq.Name,
		Theme:  createCategoryReq.Theme,
		Icon:   createCategoryReq.Icon,
	}

	return s.repo.Create(ctx, category)
}

func (s *service) FindByID(ctx context.Context, id int) (*entities.Category, error) {
	return s.repo.FindByID(ctx, id)
}

var defaultCategories = []dtos.CreateCategory{
	{
		Name:  "alimentação",
		Theme: "green",
		Icon:  "food",
	},
	{
		Name:  "transporte",
		Theme: "blue",
		Icon:  "car",
	},
	{
		Name:  "moradia",
		Theme: "purple",
		Icon:  "home",
	},
}

func (s *service) CreateDefaultCategories(
	ctx context.Context,
	userID uuid.UUID,
) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(defaultCategories))

	for _, category := range defaultCategories {
		wg.Add(1)
		go func(category dtos.CreateCategory) {
			defer wg.Done()

			_, err := s.Create(ctx, userID, &category)
			if err != nil {
				errChan <- err
			}
		}(category)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Update(
	ctx context.Context,
	id int,
	userID uuid.UUID,
	req *dtos.UpdateCategory,
) error {
	_, err := s.user.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	category, err := s.repo.FindByID(ctx, id)
	if err != nil || category.UserID != userID {
		return ErrCategoryNotFound
	}

	category.Name = req.Name
	category.Theme = req.Theme
	category.Icon = req.Icon

	return s.repo.Update(ctx, *category)
}

func (s *service) GetAll(
	ctx context.Context,
	userID uuid.UUID,
) ([]entities.Category, error) {
	return s.repo.FindMany(ctx, userID)
}
