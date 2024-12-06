package account

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(context.Context, entities.Account) (uuid.UUID, error)
	FindByID(context.Context, uuid.UUID) (*entities.Account, error)
	Update(context.Context, entities.Account) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repository {
	return &repo{db}
}

const create = `
	INSERT INTO accounts (
		user_id,
		title,
		balance	
	) VALUES ($1, $2, $3) RETURNING id
`

func (r *repo) Create(ctx context.Context, account entities.Account) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(
		ctx,
		create,
		account.UserID,
		account.Title,
		account.Balance,
	).Scan(&id)

	return id, err
}

const findByID = "SELECT * FROM accounts WHERE id = $1"

func (r *repo) FindByID(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	var account entities.Account
	err := r.db.QueryRow(ctx, findByID, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Title,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	return &account, err
}

const update = `
	UPDATE accounts
	SET title = $1, balance = $2
	WHERE id = $3
`

func (r *repo) Update(ctx context.Context, account entities.Account) error {
	_, err := r.db.Exec(
		ctx,
		update,
		account.Title,
		account.Balance,
		account,
	)

	return err
}
