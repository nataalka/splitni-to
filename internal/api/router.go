package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/nataalka/splitni-to/internal/api/handlers"
)

type RouterConfig struct {
	UserHandler   *handlers.UserHandler
	GroupHandler  *handlers.GroupHandler
	FriendHandler *handlers.FriendHandler
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
			r.Mount("/friends", friendRoutes(cfg))
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
	r.Get("/", cfg.GroupHandler.ListUsersGroups)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", cfg.GroupHandler.GetByID)
		r.Route("/members", func(r chi.Router) {
			r.Get("/", cfg.GroupHandler.ListMembers)
			r.Post("/", cfg.GroupHandler.AddMember)
			r.Delete("/", cfg.GroupHandler.RemoveMember)
		})
	})
	return r
}

func friendRoutes(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", cfg.FriendHandler.ListFriends)
	r.Get("/pending", cfg.FriendHandler.ListPendingRequests)
	r.Post("/", cfg.FriendHandler.AddFriend)
	r.Put("/{id}/accept", cfg.FriendHandler.AcceptFriend)
	r.Delete("/{id}", cfg.FriendHandler.DeleteFriendship)

	return r
}
