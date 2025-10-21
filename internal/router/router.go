package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Racuwcka/user-balance.git/internal/config"
	balanceHandler "github.com/Racuwcka/user-balance.git/internal/http-server/handlers/user-balance/balance"
	depositHandler "github.com/Racuwcka/user-balance.git/internal/http-server/handlers/user-balance/deposit"
	reserveHandler "github.com/Racuwcka/user-balance.git/internal/http-server/handlers/user-balance/reserve"
	revenueHandler "github.com/Racuwcka/user-balance.git/internal/http-server/handlers/user-balance/revenue"
	mwLogger "github.com/Racuwcka/user-balance.git/internal/http-server/middleware/logger"
	"github.com/Racuwcka/user-balance.git/internal/repo/db"
	"github.com/Racuwcka/user-balance.git/pkg/client/postgresclient"
)

func New(log *slog.Logger, cfg *config.Config) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	fs := http.FileServer(http.Dir("./../../swagger-ui"))
	router.Handle("/swagger/*", http.StripPrefix("/swagger/", fs))
	router.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/swagger.html", http.StatusFound)
	})

	client, err := postgresclient.New()
	if err != nil {
		log.Error("Postgresql is not running", slog.Any("err", err))
		return nil
	}
	repo := db.New(client)

	// ручки API
	router.Route("/api/v1/balance", func(r chi.Router) {
		r.Get("/{user_id}", balanceHandler.New(log, repo).Handle)

		r.Post("/deposit", depositHandler.New(log, repo).Handle)
		r.Post("/reserve", reserveHandler.New(log, repo).Handle)
		r.Post("/revenue", revenueHandler.New(log, repo).Handle)
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	return router
}
