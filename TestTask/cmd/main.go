package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"test/config"
	"test/internal/handler"
	ser "test/internal/service"
)

const (
	Attempts int = iota
	Retry
)

func main() {
	cfg := config.GetConfigRun()

	log := setupLogger(*cfg)
	log.Info("App was started")

	pool := ser.NewServicePool()
	handler := handler.NewHandler(log, pool)

	for i := 0; i < cfg.ServiceCnt; i++ {
		//u, err := url.Parse(fmt.Sprintf("http://web%d:8080", i+1))
		u, err := url.Parse(fmt.Sprintf("http://localhost:808%d", i+1))
		if err != nil {
			fmt.Println("URL problem:", err)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(u)

		pool.Add(ser.NewService(u, proxy))
		log.Info(fmt.Sprintf("New service was added: %s", u))
	}

	go healthCheck(*pool, *log)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Info(err.Error())
	}
}

func healthCheck(pool ser.ServicePool, log slog.Logger) {
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t.C:
			pool.HealthCheck()
			log.Info(fmt.Sprintf("Status: %s", pool.String()))
		}
	}
}

func setupLogger(cfg config.Config) *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.Level(config.LogLevels[cfg.Logger.Level]),
		}),
	)
}
