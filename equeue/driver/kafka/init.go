package kafka

import "goframe-ex/equeue/inter"

func init() {
	inter.RegisterConsumerFunc("kafka", func(config inter.MqConfig) (inter.MqConsumer, error) {
		return RegisterKafkaMqConsumer(config)
	})

	inter.RegisterProducerFunc("kafka", func(config inter.MqConfig) (inter.MqProducer, error) {
		return RegisterKafkaProducer(config)
	})
}
