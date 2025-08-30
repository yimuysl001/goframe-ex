package rabbit

import (
	"context"
	"errors"
	"fmt"
	"goframe-ex/equeue/driver/logger"
	"goframe-ex/equeue/inter"
)

func (c *RabbitMQClient) Unsubscribe(ctx context.Context, topic string) (err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.cancel()

	c.connected = false
	if c.channel != nil {
		err = c.channel.Close()
	}
	if c.connection != nil {
		err = c.connection.Close()
	}

	return err
}

// ListenReceiveMsgDo 消费数据
func (c *RabbitMQClient) ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg inter.MqMsg) error) (err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return errors.New("not connected to RabbitMQ")
	}
	if c.config.BatchSize < 1 {
		c.config.BatchSize = 10
	}

	err = c.channel.Qos(
		int(c.config.BatchSize), // prefetch count
		0,                       // prefetch size
		true,                    // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	// 开始消费
	msgs, err := c.channel.Consume(
		c.config.Queue, // queue
		topic,          // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}
	c.ctx, c.cancel = context.WithCancel(ctx)

	for {
		select {
		case <-c.ctx.Done():
			return errors.New("context canceled")
		case msg := <-msgs:
			err = receiveDo(c.ctx, inter.MqMsg{MsgId: getRandMsgId(), Topic: topic, Body: msg.Body})
			if err != nil {
				logger.Logger().Error(c.ctx, err)
			}

		}

	}

}
