package api

import (
	"errors"
	"net/http"

	"github.com/edulustosa/verdin/internal/domain/account"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/edulustosa/verdin/internal/factories"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)

	accountService := factories.MakeAccountService(api.Database)
	accounts, err := accountService.GetAll(r.Context(), userID)
	if err != nil {
		api.InternalServerError(w, "failed to get accounts", "error", err)
		return
	}

	Encode(w, http.StatusOK, JSON{"accounts": accounts})
}

func (api *API) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "accountId"))
	if err != nil {
		api.Error(w, http.StatusBadRequest, Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid account id",
			Details:    "account id must be a valid uuid",
		})
		return
	}

	accountService := factories.MakeAccountService(api.Database)
	account, err := accountService.FindByID(r.Context(), accountID)
	if err != nil {
		api.Error(w, http.StatusNotFound, Error{
			StatusCode: http.StatusNotFound,
			Message:    "account not found",
		})
		return
	}

	Encode(w, http.StatusOK, account)
}

func (api *API) EditAccount(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "accountId"))
	if err != nil {
		api.Error(w, http.StatusBadRequest, Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid account id",
			Details:    "account id must be a valid uuid",
		})
		return
	}

	req, problems, err := Decode[dtos.EditAccount](r)
	if err != nil {
		api.InvalidRequest(w, problems)
		return
	}

	accountService := factories.MakeAccountService(api.Database)
	err = accountService.Edit(r.Context(), accountID, &req)
	if err != nil {
		if errors.Is(err, account.ErrAccountNotFound) {
			api.NotFound(w, err)
			return
		}

		if errors.Is(err, account.ErrInvalidAmount) {
			api.Error(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid amount",
				Details:    "amount must be greater than or equal to 0",
			})
			return
		}

		api.InternalServerError(w, "failed to edit account", "error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *API) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	req, problems, err := Decode[dtos.EditAccount](r)
	if err != nil {
		api.InvalidRequest(w, problems)
		return
	}

	accountService := factories.MakeAccountService(api.Database)
	accountID, err := accountService.NewAccount(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, account.ErrInvalidAmount) {
			api.Error(w, http.StatusBadRequest, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid amount",
				Details:    "amount must be greater than or equal to 0",
			})
			return
		}

		if errors.Is(err, account.ErrAccountAlreadyExists) {
			api.Error(w, http.StatusConflict, Error{
				StatusCode: http.StatusConflict,
				Message:    "account already exists",
				Details:    "an account with the same title already exists",
			})
			return
		}

		api.InternalServerError(w, "failed to create account", "error", err)
		return
	}

	Encode(w, http.StatusCreated, JSON{"accountId": accountID})
}
