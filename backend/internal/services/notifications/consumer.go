package notifications

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"harmony/internal/interfaces"
	"log/slog"
	"strings"
)

type ConsumerHandlerDeps struct {
	Logger       *slog.Logger
	Service      interfaces.NotificationsService
	ConfigSarama *sarama.Config
	KafkaAddr    []string
	GroupId      string
	Topics       []string
}

type ConsumerHandler struct {
	Logger    *slog.Logger
	Service   interfaces.NotificationsService
	Consumer  sarama.ConsumerGroup
	KafkaAddr []string
	Topics    []string
}

func NewConsumerHandler(deps *ConsumerHandlerDeps) (*ConsumerHandler, error) {
	cons, err := sarama.NewConsumerGroup(deps.KafkaAddr, deps.GroupId, deps.ConfigSarama)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "sarama.NewConsumerGroup"))
		return nil, err
	}
	return &ConsumerHandler{
		Logger:    deps.Logger,
		Service:   deps.Service,
		Consumer:  cons,
		KafkaAddr: deps.KafkaAddr,
		Topics:    deps.Topics,
	}, nil
}

func (handler *ConsumerHandler) Listen() {
	for {
		err := handler.Consumer.Consume(context.Background(), handler.Topics, handler)
		if err != nil {
			handler.Logger.Error(err.Error(),
				slog.String("Error location", "Kafka.Consume"),
				slog.String("Topics", strings.Join(handler.Topics, ", ")),
				slog.String("Kafka address", strings.Join(handler.KafkaAddr, ", ")),
			)
			return
		}
	}
}
func (handler *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	handler.Logger.Info("Consumer is starting",
		slog.String("Kafka address", strings.Join(handler.KafkaAddr, ",")),
		slog.String("Topics", strings.Join(handler.Topics, ", ")),
	)
	return nil
}
func (handler *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	if handler.Consumer != nil {
		handler.Logger.Info("Consumer is stopping",
			slog.String("Topics", strings.Join(handler.Topics, ", ")),
			slog.String("Kafka address", strings.Join(handler.KafkaAddr, ",")),
		)
	}
	return nil
}

func (handler *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Println(string(msg.Value))
		session.MarkMessage(msg, string(msg.Value))
	}
	session.Commit()
	return nil
}
