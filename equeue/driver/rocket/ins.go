package rocket

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"goframe-ex/equeue/inter"
)

type RocketMq struct {
	endPoints   []string
	producerIns rocketmq.Producer
	consumerIns rocketmq.PushConsumer
}

func rewriteLog(config inter.MqConfig) {
	if config.LogLevel == "" {
		config.LogLevel = "debug"
	}

	rlog.SetLogger(&RocketMqLogger{Flag: "[rocket_mq]", LevelLog: config.LogLevel})
}

func RegisterRocketMqConsumer(config inter.MqConfig) (mqIns *RocketMq, err error) {
	rewriteLog(config)
	addr, err := primitive.NewNamesrvAddr(config.Address...)
	if err != nil {
		return nil, err
	}
	mqIns = &RocketMq{
		endPoints: config.Address,
	}

	var opts = make([]consumer.Option, 0)
	opts = append(opts, consumer.WithNameServer(addr))
	opts = append(opts, consumer.WithConsumerModel(consumer.Clustering))
	opts = append(opts, consumer.WithGroupName(config.GroupName))

	if config.UserName != "" || config.Password != "" || config.Token != "" {
		opts = append(opts, consumer.WithCredentials(primitive.Credentials{AccessKey: config.UserName, SecretKey: config.Password, SecurityToken: config.Token}))
	}

	mqIns.consumerIns, err = rocketmq.NewPushConsumer(opts...)
	return mqIns, err

}

func RegisterRocketMqProducer(config inter.MqConfig) (mqIns *RocketMq, err error) {
	addr, err := primitive.NewNamesrvAddr(config.Address...)
	if err != nil {
		return nil, err
	}
	mqIns = &RocketMq{
		endPoints: config.Address,
	}

	if config.Retry <= 0 {
		config.Retry = 0
	}

	var opts = make([]producer.Option, 0)
	opts = append(opts, producer.WithNameServer(addr))
	opts = append(opts, producer.WithRetry(config.Retry))
	opts = append(opts, producer.WithGroupName(config.GroupName))

	if config.UserName != "" || config.Password != "" || config.Token != "" {
		opts = append(opts, producer.WithCredentials(primitive.Credentials{AccessKey: config.UserName, SecretKey: config.Password, SecurityToken: config.Token}))
	}
	mqIns.producerIns, err = rocketmq.NewProducer(opts...)

	if err != nil {
		return nil, err
	}

	err = mqIns.producerIns.Start()
	if err != nil {
		return nil, err
	}
	return mqIns, nil
}
