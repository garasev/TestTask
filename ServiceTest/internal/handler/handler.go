package handler

import (
	"net/http"
	"os"

	"log/slog"
)

type Handler struct {
	Cnt    int
	logger slog.Logger
}

func NewHandler(logger slog.Logger) *Handler {
	return &Handler{
		logger: logger}
}

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	h.Cnt += 1
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	//h.logger.Debug("Pong!")
}
