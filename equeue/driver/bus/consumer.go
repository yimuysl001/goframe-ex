package bus

import (
	"context"
	"errors"
	"goframe-ex/equeue/inter"
)

func (c *LocalEventBus) Unsubscribe(ctx context.Context, topic string) (err error) {
	cancelFunc, ok := c.cancelMap[topic]
	if ok {
		cancelFunc()
	}

	if c.handler == nil {
		return
	}

	return c.bus.Unsubscribe(topic, c.handler)
}

// ListenReceiveMsgDo 消费数据
func (c *LocalEventBus) ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg inter.MqMsg) error) (err error) {
	if c.handler == nil {
		c.handler = func(mqMsg inter.MqMsg) {
			_ = receiveDo(ctx, mqMsg)
		}
	}

	err = c.bus.Subscribe(topic, c.handler)
	if err != nil {
		return err
	}
	newctx, cancel := context.WithCancel(ctx)
	c.cancelMap[topic] = cancel
	select {
	case <-newctx.Done():
		return errors.New("receive msg cancel")
	}

}
