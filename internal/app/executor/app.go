package executor

import (
	"context"
	"executor/internal/config"
	"executor/internal/httpServer/handlers"
	"executor/internal/httpServer/middleware/logger"
	"executor/internal/lib/api/logger/sl"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type App struct {
	Log         *slog.Logger
	HTTPSServer *http.Server
	Cfg         *config.HTTPServerConfig
}

func New(cfg *config.HTTPServerConfig, paths map[string][]string, api *handlers.API) *App {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(logger.New(api.Log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	for k, v := range paths {
		r.Get("/"+k, api.Execute(v))
	}

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{
		Log:         api.Log,
		HTTPSServer: srv,
		Cfg:         cfg,
	}
}

func (a *App) MustRun() {
	const op = "app.club.MustRun"

	log := a.Log.With(
		slog.String("operation", op),
	)

	if err := a.RunExecutor(); err != nil {
		log.Error("failed to start server", sl.Err(err))

		panic(err)
	}
}

func (a *App) RunExecutor() error {
	const op = "app.club.Run"

	if err := a.HTTPSServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: failed to start club server: %w", op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.cars.Run"

	err := a.HTTPSServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to stop club server: %w", op, err)
	}

	return nil
}
