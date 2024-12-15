package balance

import (
	"context"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(context.Context, entities.Balance) (uuid.UUID, error)
	FindByID(context.Context, uuid.UUID) (*entities.Balance, error)
	FindByMonth(
		ctx context.Context,
		userID uuid.UUID,
		month time.Time,
	) (*entities.Balance, error)
	Update(context.Context, entities.Balance) error
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
	INSERT INTO balances (
		user_id,
		current,
		income,
		expenses	
	) VALUES ($1, $2, $3, $4) RETURNING id;
`

func (r *repo) Create(
	ctx context.Context,
	balance entities.Balance,
) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(
		ctx,
		create,
		balance.UserID,
		balance.Current,
		balance.Income,
		balance.Expenses,
	).Scan(&id)

	return id, err
}

func scanBalance(row pgx.Row) (*entities.Balance, error) {
	var balance entities.Balance
	err := row.Scan(
		&balance.ID,
		&balance.UserID,
		&balance.Current,
		&balance.Income,
		&balance.Expenses,
		&balance.CreatedAt,
		&balance.UpdatedAt,
	)

	return &balance, err
}

const findByID = "SELECT * FROM balances WHERE id = $1;"

func (r *repo) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Balance, error) {
	row := r.db.QueryRow(ctx, findByID, id)
	return scanBalance(row)
}

const findByMonth = `
	SELECT * FROM balances
	WHERE user_id = $1
	AND EXTRACT(MONTH FROM created_at) = $2
	AND EXTRACT(YEAR FROM created_at) = $3;
`

func (r *repo) FindByMonth(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
) (*entities.Balance, error) {
	year, month, _ := date.Date()

	row := r.db.QueryRow(
		ctx,
		findByMonth,
		userID,
		month,
		year,
	)

	return scanBalance(row)
}

const update = `
	UPDATE balances
	SET current = $1,
		income = $2,
		expenses = $3,
		updated_at = $4
	WHERE id = $5;
`

func (r *repo) Update(ctx context.Context, balance entities.Balance) error {
	balance.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		ctx,
		update,
		balance.Current,
		balance.Income,
		balance.Expenses,
		balance.UpdatedAt,
		balance.ID,
	)

	return err
}
