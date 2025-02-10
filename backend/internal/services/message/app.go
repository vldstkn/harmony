package messages

import (
	"google.golang.org/grpc"
	"harmony/internal/config"
	pbmes "harmony/pkg/api/messages"
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
	lis, err := net.Listen("tcp", app.Config.MessagesAddress)
	if err != nil {
		app.Logger.Error(err.Error(), slog.String("Error location", "net.Listen")) //slog.String("Message address", app.Config.MessagesAddress)
		return
	}

	defer lis.Close()
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
	consumer, err := NewConsumer(&ConsumerGroupHandlerDeps{
		KafkaAddr: []string{app.Config.KafkaAddr},
		GroupId:   "messages",
		Service:   service,
		Topics:    []string{"message_create"},
		Logger:    app.Logger,
	})
	defer consumer.Kafka.Close()
	if err != nil {
		app.Logger.Error(err.Error(),
			slog.String("Error location", "NewConsumer"),
		)
		return
	}

	server := grpc.NewServer(opts...)
	pbmes.RegisterMessagesServer(server, handler)
	app.Logger.Info("Server start",
		slog.String("Address", app.Config.MessagesAddress),
		slog.String("Name", "Messages"),
	)

	go consumer.Listen()
	defer consumer.Kafka.Close()

	err = server.Serve(lis)
	defer server.Stop()

	if err != nil {
		panic(err)
	}

}
