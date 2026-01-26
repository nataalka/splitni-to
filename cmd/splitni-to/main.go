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

	routerConfig := api.RouterConfig{
		UserHandler: userHandler,
	}
	router := api.NewRouter(routerConfig, jwtSecret)

	log.Printf("Starting %s server on %s", cfg.Server.Environment, cfg.Server.ListenAddress)
	return http.ListenAndServe(cfg.Server.ListenAddress, router)
}
