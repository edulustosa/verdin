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
	FindMany(ctx context.Context, userID uuid.UUID) ([]entities.Account, error)
	FindByTitle(ctx context.Context, userID uuid.UUID, title string) (*entities.Account, error)
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
	SET title = $1, balance = $2, updated_at = NOW()
	WHERE id = $3;
`

func (r *repo) Update(ctx context.Context, account entities.Account) error {
	_, err := r.db.Exec(
		ctx,
		update,
		account.Title,
		account.Balance,
		account.ID,
	)

	return err
}

const findMany = "SELECT * FROM accounts WHERE user_id = $1"

func (r *repo) FindMany(
	ctx context.Context,
	userID uuid.UUID,
) ([]entities.Account, error) {
	rows, err := r.db.Query(ctx, findMany, userID)
	if err != nil {
		return nil, err
	}

	var accounts []entities.Account
	for rows.Next() {
		var account entities.Account
		err = rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Title,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

const findByTitle = "SELECT * FROM accounts WHERE user_id = $1 AND title = $2;"

func (r *repo) FindByTitle(
	ctx context.Context,
	userID uuid.UUID,
	title string,
) (*entities.Account, error) {
	row := r.db.QueryRow(
		ctx,
		findByTitle,
		userID,
		title,
	)

	var account entities.Account
	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.Title,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	return &account, err
}
