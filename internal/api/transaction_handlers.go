package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/edulustosa/verdin/internal/domain/transaction"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/edulustosa/verdin/internal/factories"
	"github.com/google/uuid"
)

func (api *API) AddTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	req, problems, err := Decode[dtos.CreateTransaction](r)
	if err != nil {
		api.InvalidRequest(w, problems)
		return
	}

	transactionService := factories.MakeTransactionService(api.Database)
	transactionID, err := transactionService.CreateTransaction(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, transaction.ErrUserNotFound) {
			api.NotFound(w, err)
			return
		}
		if errors.Is(err, transaction.ErrCategoryNotFound) {
			api.NotFound(w, err)
			return
		}
		if errors.Is(err, transaction.ErrAccountNotFound) {
			api.NotFound(w, err)
			return
		}

		if errors.Is(err, transaction.ErrInsufficientFunds) {
			api.Error(w, http.StatusConflict, Error{
				StatusCode: http.StatusConflict,
				Message:    err.Error(),
			})
			return
		}

		api.InternalServerError(w, "failed to create transaction", "error", err)
		return
	}

	Encode(w, http.StatusCreated, JSON{"transactionId": transactionID})
}

func (api *API) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	month, err := time.Parse(time.DateOnly, r.URL.Query().Get("month"))
	if err != nil {
		month = time.Now()
	}

	transactionService := factories.MakeTransactionService(api.Database)
	transactions, err := transactionService.GetMonthlyTransactions(r.Context(), userID, month)
	if err != nil {
		api.InternalServerError(w, "failed to get transactions", "error", err)
		return
	}

	Encode(w, http.StatusOK, JSON{"transactions": transactions})
}
