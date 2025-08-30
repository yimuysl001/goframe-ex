package rabbit

import "goframe-ex/equeue/inter"

func init() {
	inter.RegisterConsumerFunc("rabbit", func(config inter.MqConfig) (inter.MqConsumer, error) {
		return NewRabbitMQClient(config)
	})

	inter.RegisterProducerFunc("rabbit", func(config inter.MqConfig) (inter.MqProducer, error) {
		return NewRabbitMQClient(config)
	})
}
