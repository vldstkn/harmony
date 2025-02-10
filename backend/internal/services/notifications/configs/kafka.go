package notifsconf

import "github.com/IBM/sarama"

func NewConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	return config
}
func NewProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	return config
}
