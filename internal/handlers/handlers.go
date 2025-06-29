package handlers

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamBelugaa/iam/internal/config"
	users_handler "github.com/iamBelugaa/iam/internal/handlers/users"
	"github.com/iamBelugaa/iam/internal/okta"
	"go.uber.org/zap"
)

const (
	APIVersion1URL = "/api/v1"
)

type Config struct {
	Router *chi.Mux
	Okta   *okta.Service
	Config *config.AppConfig
	Log    *zap.SugaredLogger
}

func InitRoutes(cfg *Config) {
	cfg.Router.Use(middleware.RequestID)
	cfg.Router.Use(middleware.RealIP)
	cfg.Router.Use(middleware.Logger)
	cfg.Router.Use(middleware.Recoverer)
	cfg.Router.Use(middleware.Timeout(60 * time.Second))

	usersHandler := users_handler.New(cfg.Okta, cfg.Log)

	cfg.Router.Route(APIVersion1URL+"/users", func(r chi.Router) {
		r.Get("/", usersHandler.ListUsers)
		r.Post("/", usersHandler.CreateUser)
		r.Get("/{id}", usersHandler.GetUser)
		r.Delete("/{id}", usersHandler.DeleteUser)
	})
}
