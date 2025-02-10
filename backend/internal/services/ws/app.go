package ws

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"harmony/internal/config"
	conf "harmony/internal/services/ws/config"
	"log"
	"log/slog"
	"net/http"
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
	router := chi.NewRouter()
	server := http.Server{
		Addr:    app.Config.WebsocketAddress,
		Handler: router,
	}
	defer server.Close()

	config := conf.NewKafkaConfig()
	producer, err := sarama.NewSyncProducer([]string{app.Config.KafkaAddr}, config)
	if err != nil {
		log.Fatalf("Ошибка создания producer: %v", err)
	}
	chanKafka := make(chan []byte)
	consumer, err := NewConsumer(&ConsumerGroupHandlerDeps{
		KafkaAddr: []string{app.Config.KafkaAddr},
		GroupId:   "ws",
		Topics:    []string{"add_friend"},
		Logger:    app.Logger,
		KafkaChan: chanKafka,
	})
	defer consumer.Kafka.Close()
	if err != nil {
		app.Logger.Error(err.Error(),
			slog.String("Error location", "NewConsumer"),
		)
		return
	}

	defer producer.Close()
	NewHandlers(router, &HandlersDeps{
		Logger:    app.Logger,
		Config:    app.Config,
		Producer:  producer,
		KafkaChan: chanKafka,
	})

	app.Logger.Info("Server start",
		slog.String("Address", app.Config.WebsocketAddress),
		slog.String("Name", "Websocket"),
	)

	go consumer.Listen()
	defer consumer.Kafka.Close()

	err = server.ListenAndServe()
	defer server.Close()
	defer server.Shutdown(context.Background())
	if err != nil {
		log.Fatalln("The server is not starting, the port may be busy: ", app.Config.ApiAddress)
	}

}
