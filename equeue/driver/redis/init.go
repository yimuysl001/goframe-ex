package redis

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"goframe-ex/equeue/inter"
)

func init() {
	inter.RegisterConsumerFunc("redis", func(config inter.MqConfig) (inter.MqConsumer, error) {
		return RegisterRedis(config)
	})

	inter.RegisterProducerFunc("redis", func(config inter.MqConfig) (inter.MqProducer, error) {
		return RegisterRedis(config)
	})
}
