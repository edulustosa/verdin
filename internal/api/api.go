package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	Database *pgxpool.Pool
	JWTKey   string
}

var validate = validator.New(validator.WithRequiredStructEnabled())

type JSON map[string]any

type contextKey string

const UserIDKey contextKey = "userId"

func Encode[T any](w http.ResponseWriter, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, fmt.Sprintf("encode json: %v", err), http.StatusInternalServerError)
	}
}

func Decode[T any](r *http.Request) (v T, problems map[string]string, err error) {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}

	if err := validate.Struct(v); err != nil {
		problems := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			problems[err.Field()] = err.Error()
		}

		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}

type Error struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Details    string `json:"details"`
}

type Errors struct {
	Errors []Error `json:"errors"`
}

func (api *API) Error(w http.ResponseWriter, statusCode int, errors ...Error) {
	Encode(w, statusCode, Errors{errors})
}

func (api *API) InvalidRequest(w http.ResponseWriter, problems map[string]string) {
	errors := make([]Error, 0, len(problems))

	if problems != nil {
		for field, problem := range problems {
			errors = append(errors, Error{
				StatusCode: http.StatusBadRequest,
				Message:    fmt.Sprintf("invalid field %s", field),
				Details:    problem,
			})
		}
	} else {
		errors = append(errors, Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid request",
			Details:    "failed to decode request body",
		})
	}

	api.Error(w, http.StatusBadRequest, errors...)
}

func (api *API) InternalServerError(w http.ResponseWriter, logMsg string, slogArgs ...any) {
	slog.Error(logMsg, slogArgs...)

	api.Error(w, http.StatusInternalServerError, Error{
		StatusCode: http.StatusInternalServerError,
		Message:    "something went wrong, please try again later",
	})
}

func (api *API) NotFound(w http.ResponseWriter, err error) {
	api.Error(w, http.StatusNotFound, Error{
		StatusCode: http.StatusNotFound,
		Message:    err.Error(),
	})
}
