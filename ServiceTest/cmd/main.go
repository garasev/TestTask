package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"WebTest/config"
	"WebTest/internal/handler"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.GetConfig()

	log := setupLogger(*cfg)
	log.Info("App was started")

	handler := handler.NewHandler(*log)

	go loggingRequests(handler, log)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: Router(handler),
	}
	err := srv.ListenAndServe()
	log.Error(err.Error())

}

func loggingRequests(h *handler.Handler, log *slog.Logger) {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			log.Info(fmt.Sprintf("Requests: %d", h.Cnt))
		}
	}
}

func Router(handler *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", handler.Info)
		r.Get("/ping", handler.Ping)
	})

	return r
}

func setupLogger(cfg config.Config) *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.Level(config.LogLevels[cfg.Logger.Level]),
		}),
	)
}
