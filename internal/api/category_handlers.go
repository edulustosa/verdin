package api

import (
	"errors"
	"net/http"

	"github.com/edulustosa/verdin/internal/domain/category"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/edulustosa/verdin/internal/factories"
	"github.com/google/uuid"
)

func (api *API) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	req, problems, err := Decode[dtos.CreateCategory](r)
	if err != nil {
		api.InvalidRequest(w, problems)
		return
	}

	categoryService := factories.MakeCategoriesService(api.Database)
	id, err := categoryService.Create(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, category.ErrUserNotFound) {
			api.Error(w, http.StatusNotFound, Error{
				StatusCode: http.StatusNotFound,
				Message:    "user not found",
			})
			return
		}

		if errors.Is(err, category.ErrCategoryAlreadyExists) {
			api.Error(w, http.StatusConflict, Error{
				StatusCode: http.StatusConflict,
				Message:    "category already exists",
				Details:    "a category with the same name already exists for this user",
			})
			return
		}

		api.InternalServerError(w, "failed to create category", "error", err)
		return
	}

	Encode(w, http.StatusCreated, JSON{"categoryId": id})
}
