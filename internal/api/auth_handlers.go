package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/edulustosa/verdin/internal/auth"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/edulustosa/verdin/internal/factories"
	"github.com/golang-jwt/jwt/v5"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	req, problems, err := Decode[dtos.Register](r)
	if err != nil {
		api.InvalidRequest(w, problems)
		return
	}

	userService := factories.MakeUserService(api.Database)
	authService := auth.New(userService)
	categoryService := factories.MakeCategoriesService(api.Database)

	userID, err := authService.Register(r.Context(), &req)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			api.Error(w, http.StatusConflict, Error{
				StatusCode: http.StatusConflict,
				Message:    "user already exists",
			})
			return
		}

		api.InternalServerError(w, "failed to register user", "error", err)
		return
	}

	err = categoryService.CreateDefaultCategories(r.Context(), userID)
	if err != nil {
		api.InternalServerError(w, "failed to create default categories", "error", err)
		return
	}

	Encode(w, http.StatusCreated, JSON{"userId": userID})
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	req, problems, err := Decode[dtos.Login](r)
	if err != nil {
		api.InvalidRequest(w, problems)
		return
	}

	userService := factories.MakeUserService(api.Database)
	authService := auth.New(userService)

	userID, err := authService.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			api.Error(w, http.StatusUnauthorized, Error{
				StatusCode: http.StatusUnauthorized,
				Message:    "invalid credentials",
			})
			return
		}

		api.InternalServerError(w, "failed to login", "error", err)
		return
	}

	token, err := createToken(userID.String(), api.JWTKey)
	if err != nil {
		api.InternalServerError(w, "failed to create token", "error", err)
		return
	}

	Encode(w, http.StatusOK, JSON{"userId": userID, "token": token})
}

func createToken(userID, jwtKey string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	)

	return token.SignedString([]byte(jwtKey))
}
