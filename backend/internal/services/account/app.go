package account

import (
	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"harmony/internal/config"
	accconf "harmony/internal/services/account/config"
	pb "harmony/pkg/api/account"
	"harmony/pkg/db"
	"log"
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
	defer lis.Close()
	config := accconf.NewProducerConfig()
	producer, err := sarama.NewSyncProducer([]string{app.Config.KafkaAddr}, config)
	if err != nil {
		log.Fatalf("Ошибка создания producer: %v", err)
	}

	repository := NewRepository(app.DB)

	service := NewService(&ServiceDeps{
		Repository: repository,
		Logger:     app.Logger,
	})

	handler := NewHandler(&HandlerDeps{
		Config:   app.Config,
		Logger:   app.Logger,
		Service:  service,
		Producer: producer,
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
