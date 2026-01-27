package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nataalka/splitni-to/backend/internal/api"
	handlers2 "github.com/nataalka/splitni-to/backend/internal/api/handlers"
	"github.com/nataalka/splitni-to/backend/internal/config"
	services2 "github.com/nataalka/splitni-to/backend/internal/domain/services"
	postgres2 "github.com/nataalka/splitni-to/backend/internal/infrastructure/postgres"
)

// @title           Splitni-to API
// @version         1.0
// @description     Backend for a group expense sharing application.
// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <your_token>" to authenticate.
func main() {
	configPath := flag.String("config", "config/splitni-to.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("Fatal: %v", err)
		os.Exit(1)
	}

	err = runApp(cfg)
	if err != nil {
		log.Printf("Server failed: %v", err)
		os.Exit(1)
	}
}

func runApp(cfg *config.Config) error {
	db, err := sqlx.Connect("postgres", cfg.Storage.DSN)
	if err != nil {
		return err
	}

	jwtSecret := cfg.Auth.JWTSecret

	userRepo := postgres2.NewPostgresUserRepository(db)
	userService := services2.NewUserService(userRepo)
	userHandler := handlers2.NewUserHandler(userService, jwtSecret)

	groupRepo := postgres2.NewPostgresGroupRepository(db)
	groupService := services2.NewGroupService(groupRepo)
	groupHandler := handlers2.NewGroupHandler(groupService)

	friendRepo := postgres2.NewPostgresFriendRepository(db)
	friendService := services2.NewFriendService(friendRepo, userRepo)
	friendHandler := handlers2.NewFriendHandler(friendService)

	expenseRepo := postgres2.NewExpenseRepository(db)
	expenseService := services2.NewExpenseService(expenseRepo, userRepo)
	expenseHandler := handlers2.NewExpenseHandler(expenseService)

	routerConfig := api.RouterConfig{
		UserHandler:    userHandler,
		GroupHandler:   groupHandler,
		FriendHandler:  friendHandler,
		ExpenseHandler: expenseHandler,
		AllowedOrigins: cfg.CORS.AllowedOrigins,
	}
	router := api.NewRouter(routerConfig, jwtSecret)

	log.Printf("Starting %s server on %s", cfg.Server.Environment, cfg.Server.ListenAddress)
	return http.ListenAndServe(cfg.Server.ListenAddress, router)
}
