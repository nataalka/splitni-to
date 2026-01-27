package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/nataalka/splitni-to/backend/docs"
	handlers2 "github.com/nataalka/splitni-to/backend/internal/api/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	defaultAllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	defaultAllowedHeaders = []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}
	defaultExposedHeaders = []string{"Link"}
	defaultMaxAgeSeconds  = 300
)

type RouterConfig struct {
	UserHandler    *handlers2.UserHandler
	GroupHandler   *handlers2.GroupHandler
	FriendHandler  *handlers2.FriendHandler
	ExpenseHandler *handlers2.ExpenseHandler
	AllowedOrigins []string
}

func NewRouter(cfg RouterConfig, jwtSecret string) *chi.Mux {
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   defaultAllowedMethods,
		AllowedHeaders:   defaultAllowedHeaders,
		ExposedHeaders:   defaultExposedHeaders,
		AllowCredentials: true,
		MaxAge:           defaultMaxAgeSeconds,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

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
	r.Get("/me", cfg.UserHandler.GetMe)
	r.Get("/{user_id}", cfg.UserHandler.GetByID)
	return r
}

func groupRoutes(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", cfg.GroupHandler.Create)
	r.Get("/", cfg.GroupHandler.ListUsersGroups)
	r.Route("/{group_id}", func(r chi.Router) {
		r.Get("/", cfg.GroupHandler.GetByID)

		r.Route("/expenses", func(r chi.Router) {
			r.Post("/", cfg.ExpenseHandler.Create)
			r.Get("/", cfg.ExpenseHandler.ListByGroup)
			r.Get("/total", cfg.ExpenseHandler.GetTotal)
		})

		r.Route("/balances", func(r chi.Router) {
			r.Get("/", cfg.ExpenseHandler.GetGroupBalances)
			r.Get("/{user_id}", cfg.ExpenseHandler.GetUserBalance)
		})

		r.Route("/members", func(r chi.Router) {
			r.Get("/", cfg.GroupHandler.ListMembers)
			r.Post("/", cfg.GroupHandler.AddMember)
			r.Delete("/{user_id}", cfg.GroupHandler.RemoveMember)
		})
	})
	return r
}

func friendRoutes(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", cfg.FriendHandler.ListFriends)
	r.Get("/pending", cfg.FriendHandler.ListPendingRequests)
	r.Post("/", cfg.FriendHandler.AddFriend)
	r.Put("/{requester_id}/accept", cfg.FriendHandler.AcceptFriend)
	r.Delete("/{friend_id}", cfg.FriendHandler.DeleteFriendship)

	return r
}
