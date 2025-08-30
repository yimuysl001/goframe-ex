package rabbit

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"goframe-ex/equeue/inter"
	"math/rand"
	"time"
)

// SendByteMsg 生产数据
func (c *RabbitMQClient) SendByteMsg(ctx context.Context, topic string, body []byte) (mqMsg inter.MqMsg, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	mqMsg = inter.MqMsg{}

	if !c.connected {
		return mqMsg, errors.New("not connected to RabbitMQ")
	}
	mqMsg.MsgId = getRandMsgId()
	mqMsg.Timestamp = time.Now()
	mqMsg.Body = body

	err = c.channel.PublishWithContext(ctx, c.config.Exchange, c.config.Key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
		MessageId:   mqMsg.MsgId,
		Timestamp:   mqMsg.Timestamp,
	})

	return mqMsg, err

}

// SendMsg 按字符串类型生产数据
func (c *RabbitMQClient) SendMsg(ctx context.Context, topic string, body string) (mqMsg inter.MqMsg, err error) {
	return c.SendByteMsg(ctx, topic, []byte(body))
}

func (c *RabbitMQClient) SendDelayMsg(ctx context.Context, topic string, body string, delaySecond int64) (mqMsg inter.MqMsg, err error) {

	panic("implement me")

}

func getRandMsgId() string {
	rand.NewSource(time.Now().UnixNano())
	radium := rand.Intn(999) + 1
	timeCode := time.Now().UnixNano()
	return fmt.Sprintf("%d%.4d", timeCode, radium)
}
