package bus

import "goframe-ex/equeue/inter"

func init() {
	inter.RegisterConsumerFunc("bus", func(config inter.MqConfig) (inter.MqConsumer, error) {
		return localEventBus, nil
	})

	inter.RegisterProducerFunc("bus", func(config inter.MqConfig) (inter.MqProducer, error) {
		return localEventBus, nil
	})
}
