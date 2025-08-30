package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/equeue/driver/logger"

	"github.com/gogf/gf/v2/errors/gerror"
	"goframe-ex/equeue/inter"
)

type KaConsumer struct {
	ready        chan bool
	receiveDoFun func(ctx context.Context, mqMsg inter.MqMsg) error
}

func (r *KafkaMq) Unsubscribe(ctx context.Context, topic string) error {
	consumer := r.consumerMap[topic]
	if consumer == nil {
		return nil
	}

	consumer.ready <- true

	return nil
}

// ListenReceiveMsgDo 消费数据
func (r *KafkaMq) ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg inter.MqMsg) error) (err error) {
	if topic == "" {
		return gerror.New("RedisMq topic is empty")
	}
	_, ok := r.consumerMap[topic]
	if ok {
		return errors.New(topic + " is started ")
	}

	var consumer = &KaConsumer{
		ready:        make(chan bool),
		receiveDoFun: receiveDo,
	}

	r.consumerMap[topic] = consumer

	consumerCtx, cancel := context.WithCancel(ctx)
	go func(consumerCtx context.Context) {
		for {
			select {
			case <-consumerCtx.Done():
				break
			default:

			}

			if err = r.consumerIns.Consume(consumerCtx, []string{topic}, consumer); err != nil {
				logger.Logger().Fatalf(ctx, "kafka Error from consumer, err%+v", err)
			}

			if consumerCtx.Err() != nil {
				logger.Logger().Debugf(ctx, fmt.Sprintf("kafka consoumer stop : %v", consumerCtx.Err()))
				return
			}
			consumer.ready = make(chan bool)

		}
	}(consumerCtx)
	logger.Logger().Debug(ctx, "kafka consumer up and running!...")
	<-consumer.ready
	defer close(consumer.ready)

	logger.Logger().Debug(ctx, "kafka consumer close...")
	cancel()
	if err = r.consumerIns.Close(); err != nil {
		logger.Logger().Fatalf(ctx, "kafka Error closing client, err:%+v", err)
	}
	return

}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *KaConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *KaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *KaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	// `ConsumeClaim` 方法已经是 goroutine 调用 不要在该方法内进行 goroutine
	for message := range claim.Messages() {
		_ = c.receiveDoFun(gctx.New(), inter.MqMsg{
			RunType:   inter.ReceiveMsg,
			Topic:     message.Topic,
			Body:      message.Value,
			Offset:    message.Offset,
			Timestamp: message.Timestamp,
			Partition: message.Partition,
		})
		session.MarkMessage(message, "")
	}
	return nil
}
