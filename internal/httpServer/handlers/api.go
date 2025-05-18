package handlers

import "log/slog"

type API struct {
	Log *slog.Logger
}

func New(log *slog.Logger) *API {
	return &API{
		Log: log,
	}
}
