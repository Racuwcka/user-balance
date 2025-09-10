package router

import (
	"log"
	"net/http"

	"github.com/Racuwcka/shorter-url/pkg/client/postgresql"
	depositHandler "github.com/Racuwcka/user-balance.git/internal/handlers/balance/deposit"
	depositService "github.com/Racuwcka/user-balance.git/internal/service/deposit"
	"github.com/Racuwcka/user-balance.git/storage/postgres"
)

func New() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./../../swagger-ui"))))
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/swagger.html", http.StatusFound)
	})

	client, err := postgresql.NewClient()
	if err != nil {
		log.Fatalf("Postgresql is not running, err: %v", err)
	}
	depositHandle := depositHandler.New(depositService.New(postgres.New(client)))

	mux.HandleFunc("POST /api/v1/balance/deposit", depositHandle.Handle)

	return mux
}
