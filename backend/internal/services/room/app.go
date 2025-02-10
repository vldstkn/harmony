package room

import (
	"google.golang.org/grpc"
	"harmony/internal/config"
	pb "harmony/pkg/api/room"
	"harmony/pkg/db"
	"log/slog"
	"net"
)

type App struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *db.DB
}

type AppDeps struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *db.DB
}

func NewApp(deps *AppDeps) *App {
	return &App{
		DB:     deps.DB,
		Logger: deps.Logger,
		Config: deps.Config,
	}
}

func (app *App) Run() {
	var opts []grpc.ServerOption
	lis, err := net.Listen("tcp", app.Config.RoomAddress)
	if err != nil {
		panic(err)
	}
	repository := NewRepository(app.DB)
	service := NewService(&ServiceDeps{
		Repository: repository,
		Logger:     app.Logger,
	})
	handler := NewHandler(&HandlerDeps{
		Logger:  app.Logger,
		Service: service,
		Config:  app.Config,
	})
	server := grpc.NewServer(opts...)
	pb.RegisterRoomServer(server, handler)
	app.Logger.Info("Server start",
		slog.String("Address", app.Config.RoomAddress),
		slog.String("Name", "Room"),
	)
	err = server.Serve(lis)
	server.Stop()
	if err != nil {
		panic(err)
	}
}
