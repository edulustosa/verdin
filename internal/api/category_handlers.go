package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/edulustosa/verdin/internal/domain/category"
	"github.com/edulustosa/verdin/internal/dtos"
	"github.com/edulustosa/verdin/internal/factories"
	"github.com/go-chi/chi/v5"
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

func (api *API) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	categoryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		api.Error(w, http.StatusBadRequest, Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid category id",
		})
		return
	}

	req, problems, err := Decode[dtos.UpdateCategory](r)
	if err != nil {
		api.InvalidRequest(w, problems)
		return
	}

	categoryService := factories.MakeCategoriesService(api.Database)
	err = categoryService.Update(r.Context(), categoryID, userID, &req)
	if err != nil {
		if errors.Is(err, category.ErrCategoryNotFound) {
			api.Error(w, http.StatusNotFound, Error{
				StatusCode: http.StatusNotFound,
				Message:    "category not found",
			})
			return
		}

		api.InternalServerError(w, "failed to update category", "error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *API) GetCategories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)

	categoryService := factories.MakeCategoriesService(api.Database)
	categories, err := categoryService.GetAll(r.Context(), userID)
	if err != nil {
		api.InternalServerError(w, "failed to get categories", "error", err)
		return
	}

	Encode(w, http.StatusOK, JSON{"categories": categories})
}

func (api *API) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		api.Error(w, http.StatusBadRequest, Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid category id",
		})
		return
	}

	categoryService := factories.MakeCategoriesService(api.Database)
	category, err := categoryService.FindByID(r.Context(), id)
	if err != nil {
		api.Error(w, http.StatusNotFound, Error{
			StatusCode: http.StatusNotFound,
			Message:    "category not found",
		})
		return
	}

	Encode(w, http.StatusOK, category)
}
