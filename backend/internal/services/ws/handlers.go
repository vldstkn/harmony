package ws

import (
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"harmony/internal/config"
	"harmony/internal/services/ws/handlers"
	"log/slog"
)

type HandlersDeps struct {
	Config    *config.Config
	Logger    *slog.Logger
	Producer  sarama.SyncProducer
	KafkaChan chan []byte
}

func NewHandlers(router *chi.Mux, deps *HandlersDeps) {
	handlers.NewAppHandler(router, &handlers.AppHandlerDeps{
		Logger:    deps.Logger,
		Config:    deps.Config,
		Producer:  deps.Producer,
		ChanKafka: deps.KafkaChan,
	})
	handlers.NewChatHandler(router, &handlers.ChatHandlerDeps{
		Logger:   deps.Logger,
		Config:   deps.Config,
		Producer: deps.Producer,
	})

}
