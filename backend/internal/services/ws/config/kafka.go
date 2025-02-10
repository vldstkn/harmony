package wsconf

import "github.com/IBM/sarama"

func NewKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.Return.Successes = true
	//config.Producer.Idempotent = true // Включаем идемпотентность
	//config.Producer.RequiredAcks = sarama.WaitForAll
	return config
}

func NewConsumerKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = false // Отключаем авто-коммит
	return config
}
