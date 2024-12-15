package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/edulustosa/verdin/internal/domain/balance"
	"github.com/edulustosa/verdin/internal/factories"
	"github.com/google/uuid"
)

func (api *API) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)

	month, err := time.Parse(time.DateOnly, r.URL.Query().Get("month"))
	if err != nil {
		api.Error(w, http.StatusBadRequest, Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid month",
			Details:    "failed to parse query, must be in the format YYYY-MM-DD",
		})
		return
	}

	balanceService := factories.MakeBalanceService(api.Database)
	monthlyBalance, err := balanceService.FindByMonth(r.Context(), userID, month)
	if err != nil {
		if errors.Is(err, balance.ErrMonthNotFound) {
			api.Error(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
			return
		}

		api.InternalServerError(w, "failed to get balance", "error", err)
		return
	}

	Encode(w, http.StatusOK, monthlyBalance)
}
