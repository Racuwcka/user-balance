package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Racuwcka/shorter-url/pkg/closer"
	"github.com/Racuwcka/user-balance.git/internal/router"
)

func Run(ctx context.Context) error {
	//cfg := config.MustLoadConfig()
	shutdowner := &closer.Closer{}

	addr := ":" + "8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: router.New(),
	}

	shutdowner.Add(srv.Shutdown)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return shutdowner.Close(shutdownCtx)
}
