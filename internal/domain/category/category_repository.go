package category

import (
	"context"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(context.Context, entities.Category) (int, error)
	FindByName(ctx context.Context, userID uuid.UUID, name string) (*entities.Category, error)
	FindByID(context.Context, int) (*entities.Category, error)
	Update(context.Context, entities.Category) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repository {
	return &repo{
		db,
	}
}

const create = `
	INSERT INTO categories (
		user_id,
		name,
		theme,
		icon	
	) VALUES ($1, $2, $3, $4) RETURNING id;
`

func (r *repo) Create(ctx context.Context, category entities.Category) (int, error) {
	var id int
	err := r.db.QueryRow(
		ctx,
		create,
		category.UserID,
		category.Name,
		category.Theme,
		category.Icon,
	).Scan(&id)

	return id, err
}

const findByName = "SELECT * FROM categories WHERE user_id = $1 AND name = $2;"

func (r *repo) FindByName(
	ctx context.Context,
	userID uuid.UUID,
	name string,
) (*entities.Category, error) {
	row := r.db.QueryRow(
		ctx,
		findByName,
		userID,
		name,
	)

	return scanCategory(row)
}

func scanCategory(row pgx.Row) (*entities.Category, error) {
	var category entities.Category
	err := row.Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Theme,
		&category.Icon,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	return &category, err
}

const findByID = "SELECT * FROM categories WHERE id = $1;"

func (r *repo) FindByID(ctx context.Context, id int) (*entities.Category, error) {
	row := r.db.QueryRow(ctx, findByID, id)

	return scanCategory(row)
}

const update = `
	UPDATE categories 
	SET name = $1, 
		theme = $2,
		icon = $3,
		updated_at = $4
	WHERE id = $5;
`

func (r *repo) Update(ctx context.Context, category entities.Category) error {
	category.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		ctx,
		update,
		category.Name,
		category.Theme,
		category.Icon,
		category.UpdatedAt,
		category.ID,
	)

	return err
}
