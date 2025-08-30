package disk

import "goframe-ex/equeue/inter"

func init() {
	inter.RegisterConsumerFunc("disk", func(config inter.MqConfig) (inter.MqConsumer, error) {
		return RegisterDiskMqConsumer(config)
	})

	inter.RegisterProducerFunc("disk", func(config inter.MqConfig) (inter.MqProducer, error) {
		return RegisterDiskMqProducer(config)
	})
}
