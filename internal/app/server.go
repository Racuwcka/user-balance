package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Racuwcka/shorter-url/pkg/closer"

	"github.com/Racuwcka/user-balance.git/internal/config"
	"github.com/Racuwcka/user-balance.git/internal/router"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run(ctx context.Context) error {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router.New(log, cfg),
	}

	shutdowner := &closer.Closer{}
	shutdowner.Add(srv.Shutdown)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("listen and serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return shutdowner.Close(shutdownCtx)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
