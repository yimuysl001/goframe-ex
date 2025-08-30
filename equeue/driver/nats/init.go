package nats

import "goframe-ex/equeue/inter"

func init() {
	inter.RegisterConsumerFunc("nats", func(config inter.MqConfig) (inter.MqConsumer, error) {
		return RegisterNats(config)
	})

	inter.RegisterProducerFunc("nats", func(config inter.MqConfig) (inter.MqProducer, error) {
		return RegisterNats(config)
	})
}
