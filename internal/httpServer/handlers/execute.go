package handlers

import (
	"executor/internal/lib/api/flushWriter"
	"executor/internal/lib/api/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
)

func (a *API) Execute(command []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Execute"

		log := a.Log.With(
			slog.String("operation", op),
			slog.String("command", strings.Join(command, " ")),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("request started")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		fw := &flushWriter.FlushWriter{
			W:       w,
			Flusher: w.(http.Flusher),
		}

		cmd := exec.CommandContext(r.Context(), command[0], command[1:]...)
		cmd.Stdout = fw
		cmd.Stderr = fw

		err := cmd.Run()
		if err != nil {
			if r.Context().Err() != nil {
				log.Info("context closed", sl.Err(r.Context().Err()))
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Warn("failed to execute command", sl.Err(err))
			return
		}
	}
}
