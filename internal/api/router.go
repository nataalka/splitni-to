package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nataalka/splitni-to/internal/api/handlers"
)

type RouterConfig struct {
	UserHandler *handlers.UserHandler
}

func NewRouter(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", userRoutes(cfg))
	})

	return r
}

func userRoutes(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/register", cfg.UserHandler.Register)
	r.Get("/{id}", cfg.UserHandler.GetByID)
	return r
}
