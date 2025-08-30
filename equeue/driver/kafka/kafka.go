package kafka

import (
	"github.com/IBM/sarama"
	"goframe-ex/equeue/inter"
	"time"
)

type KafkaMq struct {
	Partitions  int32
	producerIns sarama.AsyncProducer
	consumerIns sarama.ConsumerGroup
	consumerMap map[string]*KaConsumer
}

func RegisterKafkaMqConsumer(config inter.MqConfig) (client inter.MqConsumer, err error) {
	mqIns := &KafkaMq{}
	kfkVersion, err := sarama.ParseKafkaVersion(config.Version)
	if err != nil {
		return
	}
	if !validateVersion(kfkVersion) {
		kfkVersion = sarama.V2_4_0_0
	}

	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Version = kfkVersion

	if config.UserName != "" {
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = config.UserName
		conf.Net.SASL.Password = config.Password
	}
	// 默认按随机方式消费
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	conf.Consumer.Offsets.AutoCommit.Interval = 10 * time.Millisecond
	conf.ClientID = "equeue-consumer-" + config.GroupName
	mqIns.consumerIns, err = sarama.NewConsumerGroup(config.Address, config.GroupName, conf)

	return mqIns, err
}

// RegisterKafkaProducer 注册并启动生产者接口实现
func RegisterKafkaProducer(config inter.MqConfig) (client inter.MqProducer, err error) {
	mqIns := &KafkaMq{}
	var clientId = "equeue-producer-" + config.GroupName

	// 这里如果使用go程需要处理chan同步问题
	kfkVersion, err := sarama.ParseKafkaVersion(config.Version)
	if err != nil {
		return
	}
	if !validateVersion(kfkVersion) {
		kfkVersion = sarama.V2_4_0_0
	}

	brokers := config.Address
	conf := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	conf.Producer.RequiredAcks = sarama.WaitForAll
	// 随机向partition发送消息
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	conf.Producer.Return.Successes = true

	conf.Producer.Return.Errors = true
	conf.Producer.Compression = sarama.CompressionNone
	conf.ClientID = clientId

	conf.Version = kfkVersion
	if config.UserName != "" {
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = config.UserName
		conf.Net.SASL.Password = config.Password
	}

	mqIns.producerIns, err = sarama.NewAsyncProducer(brokers, conf)
	if err != nil {
		return
	}

	return mqIns, nil
}

// validateVersion 验证版本是否有效
func validateVersion(version sarama.KafkaVersion) bool {
	for _, item := range sarama.SupportedVersions {
		if version.String() == item.String() {
			return true
		}
	}
	return false
}
