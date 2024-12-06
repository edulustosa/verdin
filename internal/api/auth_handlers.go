package api

import (
	"errors"
	"net/http"

	"github.com/edulustosa/verdin/internal/auth"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/edulustosa/verdin/internal/factories"
	"golang.org/x/exp/slog"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	req, problems, err := Decode[dtos.Register](r)
	if err != nil {
		api.InvalidRequest(w, problems)
		return
	}

	userService := factories.MakeUserService(api.Database)
	authService := auth.New(userService)

	userID, err := authService.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			api.Error(w, http.StatusConflict, Error{
				StatusCode: http.StatusConflict,
				Message:    "user already exists",
			})
			return
		}

		slog.Error("failed to register user", "error", err)
		api.InternalServerError(w)
		return
	}

	Encode(w, http.StatusCreated, JSON{"userId": userID})
}
