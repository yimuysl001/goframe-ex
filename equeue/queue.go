package equeue

import (
	_ "goframe-ex/equeue/driver/bus"
	_ "goframe-ex/equeue/driver/disk"
	_ "goframe-ex/equeue/driver/kafka"
	_ "goframe-ex/equeue/driver/nats"
	_ "goframe-ex/equeue/driver/rabbit"
	_ "goframe-ex/equeue/driver/redis"
	_ "goframe-ex/equeue/driver/rocket"
	"goframe-ex/equeue/inter"
)

func Listen(name ...string) inter.MqConsumer {
	df := defName
	if len(name) > 0 && name[0] != "" {
		df = name[0]
	}
	config, ok := mqConfig[df]
	if !ok {
		panic("mqConfig[" + df + "] is nil")
	}

	consumerFunc := inter.GetConsumerFunc(config.Driver)
	if consumerFunc == nil {
		panic("no driver name " + config.Driver)
	}
	ins, err := consumerFunc(config)
	if err != nil {
		panic(err)
	}
	return ins

}

func Mq(name ...string) inter.MqProducer {
	df := defName
	if len(name) > 0 && name[0] != "" {
		df = name[0]
	}
	return producerMap.GetOrSetFunc(df, func() interface{} {
		config, ok := mqConfig[df]
		if !ok {
			panic("mqConfig[" + df + "] is nil")
		}
		producerFunc := inter.GetProducerFunc(config.Driver)
		if producerFunc == nil {
			panic("no driver name " + config.Driver)
		}
		ins, err := producerFunc(config)
		if err != nil {
			panic(err)
		}
		return ins

	}).(inter.MqProducer)
}

func SetConfig(name string, config inter.MqConfig) {

	_, found := producerMap.Search(name)
	if found {
		producerMap.Remove(name)
	}

	mqConfig[name] = config

}
