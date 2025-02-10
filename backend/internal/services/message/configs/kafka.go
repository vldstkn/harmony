package kafkaconf

import "github.com/IBM/sarama"

func NewConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = false // Отключаем авто-коммит
	return config
}
