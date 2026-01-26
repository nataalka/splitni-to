package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nataalka/splitni-to/internal/api"
	"github.com/nataalka/splitni-to/internal/api/handlers"
	"github.com/nataalka/splitni-to/internal/config"
	"github.com/nataalka/splitni-to/internal/domain/services"
	"github.com/nataalka/splitni-to/internal/infrastructure/postgres"
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

	userRepo := postgres.NewPostgresUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, jwtSecret)

	groupRepo := postgres.NewPostgresGroupRepository(db)
	groupService := services.NewGroupService(groupRepo)
	groupHandler := handlers.NewGroupHandler(groupService)

	friendRepo := postgres.NewPostgresFriendRepository(db)
	friendService := services.NewFriendService(friendRepo, userRepo)
	friendHandler := handlers.NewFriendHandler(friendService)

	routerConfig := api.RouterConfig{
		UserHandler:   userHandler,
		GroupHandler:  groupHandler,
		FriendHandler: friendHandler,
	}
	router := api.NewRouter(routerConfig, jwtSecret)

	log.Printf("Starting %s server on %s", cfg.Server.Environment, cfg.Server.ListenAddress)
	return http.ListenAndServe(cfg.Server.ListenAddress, router)
}
