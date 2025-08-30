package rocket

import "goframe-ex/equeue/inter"

func init() {
	inter.RegisterConsumerFunc("rocketmq", func(config inter.MqConfig) (inter.MqConsumer, error) {
		return RegisterRocketMqConsumer(config)
	})

	inter.RegisterProducerFunc("rocketmq", func(config inter.MqConfig) (inter.MqProducer, error) {
		return RegisterRocketMqProducer(config)
	})
}
