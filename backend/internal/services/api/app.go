package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"harmony/internal/config"
	"harmony/internal/services/api/middleware"
	"log"
	"log/slog"
	"net/http"
)

type App struct {
	Config *config.Config
	Logger *slog.Logger
}

type AppDeps struct {
	Config *config.Config
	Logger *slog.Logger
}

func NewApp(deps *AppDeps) *App {
	return &App{
		Config: deps.Config,
		Logger: deps.Logger,
	}
}

func (app *App) Run() {
	router := chi.NewRouter()
	mw := chi.Chain(middleware.CORS)
	router.Use(mw...)
	service := NewService(app.Config.JWTSecret)
	NewHandlers(router, &HandlersDeps{
		Logger:  app.Logger,
		Service: service,
		Config:  app.Config,
	})
	server := http.Server{
		Addr:    app.Config.ApiAddress,
		Handler: router,
	}
	app.Logger.Info("Server start",
		slog.String("Address", app.Config.ApiAddress),
		slog.String("Name", "Api"),
	)

	err := server.ListenAndServe()
	defer server.Shutdown(context.Background())
	defer server.Close()

	if err != nil {
		log.Fatalln("The server is not starting, the port may be busy: ", app.Config.ApiAddress)
	}
}
