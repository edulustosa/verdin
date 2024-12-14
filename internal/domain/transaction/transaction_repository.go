package transaction

import (
	"context"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(context.Context, entities.Transaction) (uuid.UUID, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repository {
	return &repo{db}
}

const create = `
	INSERT INTO transactions (
		user_id,
		category_id,
		account_id,
		balance_id,	
		title,
		description,
		amount,
		type,
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id;
`

func (r *repo) Create(ctx context.Context, t entities.Transaction) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(
		ctx,
		create,
		t.UserID,
		t.CategoryID,
		t.AccountID,
		t.BalanceID,
		t.Title,
		t.Description,
		t.Amount,
		t.Type,
	).Scan(&id)

	return id, err
}
