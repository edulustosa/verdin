package transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/edulustosa/verdin/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(context.Context, entities.Transaction) (int, error)
	FindManyByMonth(
		ctx context.Context,
		userID uuid.UUID,
		month time.Time,
	) ([]entities.Transaction, error)
	FindManyByMonthAndCategory(
		ctx context.Context,
		userID uuid.UUID,
		month time.Time,
		categoryID int,
	) ([]entities.Transaction, error)
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
		type
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id;
`

func (r *repo) Create(ctx context.Context, t entities.Transaction) (int, error) {
	var id int
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

const findManyByMonth = `
	SELECT * FROM transactions
	WHERE user_id = $1
	AND EXTRACT(MONTH FROM created_at) = $2
	AND EXTRACT(YEAR FROM created_at) = $3
	ORDER BY created_at DESC;
`

func (r *repo) FindManyByMonth(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
) ([]entities.Transaction, error) {
	year, month, _ := date.Date()

	rows, err := r.db.Query(ctx, findManyByMonth, userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []entities.Transaction
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func scanTransaction(row pgx.Row) (entities.Transaction, error) {
	var t entities.Transaction
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.CategoryID,
		&t.AccountID,
		&t.BalanceID,
		&t.Title,
		&t.Description,
		&t.Amount,
		&t.Type,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	return t, err
}

const findManyByMonthAndCategory = `
	SELECT * FROM transactions
	WHERE user_id = $1
	AND EXTRACT(MONTH FROM created_at) = $2
	AND EXTRACT(YEAR FROM created_at) = $3
	AND category_id = $4	
	ORDER BY created_at DESC;
`

func (r *repo) FindManyByMonthAndCategory(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
	categoryID int,
) ([]entities.Transaction, error) {
	year, month, _ := date.Date()

	rows, err := r.db.Query(
		ctx,
		findManyByMonthAndCategory,
		userID,
		month,
		year,
		categoryID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []entities.Transaction
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}
