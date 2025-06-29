package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iamBelugaa/iam/internal/config"
	"github.com/iamBelugaa/iam/internal/handlers"
	"github.com/iamBelugaa/iam/internal/okta"
	"github.com/iamBelugaa/iam/pkg/logger"
	oktaPkg "github.com/iamBelugaa/iam/pkg/okta"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	log := logger.New("flexera-iam")
	defer func() {
		if err := log.Sync(); err != nil {
			log.Errorw("sync error", "error", err)
		}
	}()

	if err := godotenv.Load(); err != nil {
		log.Fatalw("error loading envs", "error", err)
	}

	if err := run(log); err != nil {
		log.Errorw("startup error", "error", err)
		if err := log.Sync(); err != nil {
			log.Errorw("sync error", "error", err)
		}
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	cfg := config.Load()
	oktaClient, err := oktaPkg.NewClient(cfg.Okta)
	if err != nil {
		return err
	}

	router := chi.NewRouter()
	oktaService := okta.NewOktaService(oktaClient.GetSDK())

	handlers.InitRoutes(&handlers.Config{
		Log:    log,
		Config: cfg,
		Router: router,
		Okta:   oktaService,
	})

	server := http.Server{
		Handler:      router,
		Addr:         ":" + cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	serverErrors := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Infow("server starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutting down server", "signal", sig)
		defer log.Infow("shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
