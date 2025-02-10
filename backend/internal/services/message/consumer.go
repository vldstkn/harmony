package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"harmony/internal/interfaces"
	kafkaconf "harmony/internal/services/message/configs"
	"log/slog"
	"strings"
)

type ConsumerGroupHandlerDeps struct {
	KafkaAddr []string
	GroupId   string
	Topics    []string
	Service   interfaces.MessageService
	Logger    *slog.Logger
}

type ConsumerGroupHandler struct {
	Service   interfaces.MessageService
	Kafka     sarama.ConsumerGroup
	Logger    *slog.Logger
	KafkaAddr []string
	Topics    []string
}

func NewConsumer(deps *ConsumerGroupHandlerDeps) (*ConsumerGroupHandler, error) {
	config := kafkaconf.NewConsumerConfig()
	consumerGroup, err := sarama.NewConsumerGroup(deps.KafkaAddr, deps.GroupId, config)
	if err != nil {
		deps.Logger.Error(err.Error(),
			slog.String("Error location", "sarama.NewConsumerGroup"))
		return nil, err
	}
	consumer := &ConsumerGroupHandler{
		Service:   deps.Service,
		Kafka:     consumerGroup,
		Logger:    deps.Logger,
		Topics:    deps.Topics,
		KafkaAddr: deps.KafkaAddr,
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
		case "message_create":
			var message MessageReq
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				h.Logger.Error(err.Error(), slog.String("error location", "json.Unmarshal"))
			} else {
				h.Service.Create(message.SenderId, message.RoomId, message.Message)
			}
		}
		fmt.Println("Mark message!!", string(msg.Value))
		session.MarkMessage(msg, string(msg.Value))
	}
	session.Commit()
	return nil
}
