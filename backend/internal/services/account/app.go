package account

import (
	"google.golang.org/grpc"
	"harmony/internal/config"
	pb "harmony/pkg/api/account"
	"harmony/pkg/db"
	"log/slog"
	"net"
)

type AppDeps struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *db.DB
}

type App struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *db.DB
}

func NewApp(deps *AppDeps) *App {
	return &App{
		Config: deps.Config,
		Logger: deps.Logger,
		DB:     deps.DB,
	}
}

func (app *App) Run() {
	var opts []grpc.ServerOption
	lis, err := net.Listen("tcp", app.Config.AccountAddress)
	if err != nil {
		panic(err)
	}

	repository := NewRepository(app.DB)

	service := NewService(&ServiceDeps{
		Repository: repository,
		Logger:     app.Logger,
	})

	handler := NewHandler(&HandlerDeps{
		Config:  app.Config,
		Logger:  app.Logger,
		Service: service,
	})

	server := grpc.NewServer(opts...)
	pb.RegisterAccountServer(server, handler)
	app.Logger.Info("Server start",
		slog.String("Address", app.Config.AccountAddress),
		slog.String("Name", "Account"),
	)
	err = server.Serve(lis)
	server.Stop()
	if err != nil {
		panic(err)
	}
}
