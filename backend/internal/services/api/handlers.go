package api

import (
	"github.com/go-chi/chi/v5"
	"harmony/internal/config"
	"harmony/internal/interfaces"
	"harmony/internal/services/api/handlers"
	"log/slog"
)

type HandlersDeps struct {
	Config  *config.Config
	Logger  *slog.Logger
	Service interfaces.ApiService
}

func NewHandlers(router *chi.Mux, deps *HandlersDeps) {
	handlers.NewAccountHandler(router, &handlers.AccountHandlerDeps{
		Config:  deps.Config,
		Service: deps.Service,
		Logger:  deps.Logger,
	})
	handlers.NewRoomHandler(router, &handlers.RoomHandlerDeps{
		Logger: deps.Logger,
		Config: deps.Config,
	})
}
