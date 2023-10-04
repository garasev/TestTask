package handler

import (
	"log/slog"
	"net/http"
	ser "test/internal/service"
)

type Handler struct {
	pool   *ser.ServicePool
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger, pool *ser.ServicePool) *Handler {
	return &Handler{
		logger: logger,
		pool:   pool,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	peer := h.pool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
