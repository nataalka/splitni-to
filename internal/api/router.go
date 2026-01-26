package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/nataalka/splitni-to/internal/api/handlers"
)

type RouterConfig struct {
	UserHandler  *handlers.UserHandler
	GroupHandler *handlers.GroupHandler
}

func NewRouter(cfg RouterConfig, jwtSecret string) *chi.Mux {
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		// public routes
		r.Group(func(r chi.Router) {
			r.Post("/register", cfg.UserHandler.Register)
			r.Post("/login", cfg.UserHandler.Login)
		})

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Mount("/users", userRoutes(cfg))
			r.Mount("/groups", groupRoutes(cfg))
		})
	})

	return r
}

func userRoutes(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}", cfg.UserHandler.GetByID)
	return r
}

func groupRoutes(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", cfg.GroupHandler.Create)
	return r
}
