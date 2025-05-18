package app

import (
	"executor/internal/app/executor"
	"executor/internal/config"
	"executor/internal/httpServer/handlers"
	"log/slog"
)

type App struct {
	Executor *executor.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	api := handlers.New(log)

	app := executor.New(cfg.HttpServer, cfg.CommandPaths, api)

	return &App{
		Executor: app,
	}
}
