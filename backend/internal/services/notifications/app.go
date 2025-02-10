package notifications

import (
	"harmony/internal/config"
	notifsconf "harmony/internal/services/notifications/configs"
	"harmony/pkg/db"
	"log/slog"
)

type AppDeps struct {
	Logger *slog.Logger
	DB     *db.DB
	Config *config.Config
}

type App struct {
	Logger *slog.Logger
	DB     *db.DB
	Config *config.Config
}

func NewApp(deps *AppDeps) *App {
	return &App{
		Config: deps.Config,
		Logger: deps.Logger,
		DB:     deps.DB,
	}
}

func (app *App) Run() {
	repository := NewRepository(app.DB)
	service := NewService(&ServiceDeps{
		Repository: repository,
		Logger:     app.Logger,
	})
	config := notifsconf.NewConsumerConfig()
	consumer, err := NewConsumerHandler(&ConsumerHandlerDeps{
		Service:      service,
		Logger:       app.Logger,
		KafkaAddr:    []string{app.Config.KafkaAddr},
		ConfigSarama: config,
		Topics:       []string{"create_notification"},
		GroupId:      "notifications",
	})
	if err != nil {
		return
	}
	defer consumer.Consumer.Close()

	app.Logger.Info("Server start",
		slog.String("Address", app.Config.NotificationsAddress),
		slog.String("Name", "Messages"),
	)

	consumer.Listen()

}
