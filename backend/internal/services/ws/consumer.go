package ws

import (
	"context"
	"github.com/IBM/sarama"
	wsconf "harmony/internal/services/ws/config"
	"log/slog"
	"strings"
)

type ConsumerGroupHandlerDeps struct {
	KafkaAddr []string
	GroupId   string
	Topics    []string
	Logger    *slog.Logger
	KafkaChan chan []byte
}

type ConsumerGroupHandler struct {
	Kafka     sarama.ConsumerGroup
	Logger    *slog.Logger
	KafkaAddr []string
	Topics    []string
	KafkaChan chan []byte
}

func NewConsumer(deps *ConsumerGroupHandlerDeps) (*ConsumerGroupHandler, error) {
	config := wsconf.NewConsumerKafkaConfig()
	consumerGroup, err := sarama.NewConsumerGroup(deps.KafkaAddr, deps.GroupId, config)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "sarama.NewConsumerGroup"))
		return nil, err
	}
	consumer := &ConsumerGroupHandler{
		Kafka:     consumerGroup,
		Logger:    deps.Logger,
		Topics:    deps.Topics,
		KafkaAddr: deps.KafkaAddr,
		KafkaChan: deps.KafkaChan,
	}
	return consumer, nil
}

func (cons *ConsumerGroupHandler) Listen() {
	for {
		err := cons.Kafka.Consume(context.Background(), cons.Topics, cons)
		if err != nil {
			cons.Logger.Error(err.Error(),
				slog.String("Error location", "Kafka.Consume"),
				slog.String("Topics", strings.Join(cons.Topics, ", ")),
				slog.String("Kafka address", strings.Join(cons.KafkaAddr, ", ")),
			)
			return
		}
	}
}
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	h.Logger.Info("Consumer is starting",
		slog.String("Kafka address", strings.Join(h.KafkaAddr, ",")),
		slog.String("Topics", strings.Join(h.Topics, ", ")),
	)
	return nil
}
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	if h.Kafka != nil {
		h.Logger.Info("Consumer is stopping",
			slog.String("Topics", strings.Join(h.Topics, ", ")),
			slog.String("Kafka address", strings.Join(h.KafkaAddr, ",")),
		)
	}
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		switch msg.Topic {
		case "add_friend":
			h.KafkaChan <- msg.Value
		}
		session.MarkMessage(msg, string(msg.Value))
	}
	session.Commit()
	return nil
}
