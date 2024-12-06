package router

import (
	"net/http"

	"github.com/edulustosa/verdin/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(api *api.API) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Post("/register", api.Register)

	return r
}
