package user

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user entities.User) (uuid.UUID, error)
	FindByID(context.Context, uuid.UUID) (*entities.User, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repository {
	return &repo{
		db,
	}
}

const findByEmail = "SELECT * FROM users WHERE email = $1;"

func (r *repo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	row := r.db.QueryRow(
		ctx,
		findByEmail,
		email,
	)

	return scanUser(row)
}

func scanUser(row pgx.Row) (*entities.User, error) {
	var u entities.User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	return &u, err
}

const create = `
	INSERT INTO users (
		username,
		email,
		password_hash
	) VALUES ($1, $2, $3) RETURNING id;
`

func (r *repo) Create(ctx context.Context, user entities.User) (uuid.UUID, error) {
	row := r.db.QueryRow(
		ctx,
		create,
		user.Username,
		user.Email,
		user.PasswordHash,
	)

	var id uuid.UUID
	err := row.Scan(&id)

	return id, err
}

const findByID = "SELECT * FROM users WHERE id = $1;"

func (r *repo) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	row := r.db.QueryRow(ctx, findByID, id)
	return scanUser(row)
}
