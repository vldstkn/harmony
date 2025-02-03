package email

import (
	"google.golang.org/grpc"
	"harmony/internal/config"
	pb "harmony/pkg/api/email"
	"log/slog"
	"net"
)

type AppDeps struct {
	Config *config.Config
	Logger *slog.Logger
}

type App struct {
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
	var opts []grpc.ServerOption
	lis, err := net.Listen("tcp", app.Config.EmailAddress)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(opts...)
	service := NewService(&ServiceDeps{
		Config: app.Config,
	})
	handler := NewHandler(HandlerDeps{
		Config:  app.Config,
		Logger:  app.Logger,
		Service: service,
	})
	pb.RegisterEmailServer(server, handler)
	app.Logger.Info("Server start",
		slog.String("Address", app.Config.EmailAddress),
		slog.String("Name", "Email"),
	)
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
