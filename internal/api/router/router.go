package router

import (
	"net/http"

	"github.com/edulustosa/verdin/internal/api"
	"github.com/edulustosa/verdin/internal/api/middlewares"
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
	r.Post("/login", api.Login)

	middlewares := &middlewares.Middlewares{JWTKey: api.JWTKey}
	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("pong"))
		})
	})

	return r
}
