package inter

var (
	consumerFuncMap = make(map[string]func(config MqConfig) (MqConsumer, error))
	producerFuncMap = make(map[string]func(config MqConfig) (MqProducer, error))
)

func RegisterConsumerFunc(driver string, f func(config MqConfig) (MqConsumer, error)) {
	consumerFuncMap[driver] = f
}

func RegisterProducerFunc(driver string, f func(config MqConfig) (MqProducer, error)) {
	producerFuncMap[driver] = f
}

func GetConsumerFunc(driver string) func(config MqConfig) (MqConsumer, error) {
	return consumerFuncMap[driver]
}

func GetProducerFunc(driver string) func(config MqConfig) (MqProducer, error) {
	return producerFuncMap[driver]
}
