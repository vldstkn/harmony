package accconf

import "github.com/IBM/sarama"

func NewProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.Return.Successes = true
	//config.Producer.Idempotent = true // Включаем идемпотентность
	//config.Producer.RequiredAcks = sarama.WaitForAll
	return config
}
