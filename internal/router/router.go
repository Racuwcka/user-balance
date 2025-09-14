package router

import (
	"log/slog"
	"net/http"

	"github.com/Racuwcka/shorter-url/pkg/client/postgresql"
	"github.com/Racuwcka/user-balance.git/internal/config"
	depositHandler "github.com/Racuwcka/user-balance.git/internal/http-server/handlers/balance/deposit"
	mwLogger "github.com/Racuwcka/user-balance.git/internal/http-server/middleware/logger"
	depositService "github.com/Racuwcka/user-balance.git/internal/service/deposit"
	"github.com/Racuwcka/user-balance.git/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	client, err := postgresql.NewClient()
	if err != nil {
		log.Error("Postgresql is not running, err: %v", err)
	}

	depositHandle := depositHandler.New(depositService.New(postgres.New(client)))

	// ручки API
	router.Route("/api/v1/balance", func(r chi.Router) {
		r.Post("/deposit", depositHandle.Handle)
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	return router
}
