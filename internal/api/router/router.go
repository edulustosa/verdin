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

		r.Post("/category", api.CreateCategory)
		r.Put("/category/{id}", api.UpdateCategory)
		r.Get("/category", api.GetCategories)
		r.Get("/category/{id}", api.GetCategory)

		r.Get("/accounts", api.GetAccounts)
		r.Get("/accounts/{accountId}", api.GetAccount)
		r.Post("/accounts", api.CreateAccount)
		r.Put("/accounts/{accountId}", api.EditAccount)

		r.Post("/transactions", api.AddTransaction)
		r.Get("/transactions", api.GetTransactions)

		r.Get("/balance", api.GetBalance)
	})

	return r
}
