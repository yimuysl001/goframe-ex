package nats

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"goframe-ex/equeue/inter"
	"time"
)

func (r *NatsMq) Unsubscribe(ctx context.Context, topic string) error {
	cancelFunc, ok := r.cancelMap[topic]
	if ok {
		cancelFunc()
	}

	return nil
}

// ListenReceiveMsgDo 消费数据
func (r *NatsMq) ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg inter.MqMsg) error) (err error) {
	if topic == "" {
		return gerror.New("RedisMq topic is empty")
	}

	_, ok := r.cancelMap[topic]
	if ok {
		return errors.New(topic + " is started ")
	}

	newCtx, cancel := context.WithCancel(ctx)

	r.cancelMap[topic] = cancel

	sub, err := r.js.PullSubscribe(r.subjects+"."+topic, r.subjects, nats.Context(newCtx), nats.PullMaxWaiting(r.maxWaiting))

	for {
		select {
		case <-newCtx.Done():
			break
		default:
		}

		msgs, _ := sub.Fetch(r.batch, nats.Context(newCtx))
		for _, msg := range msgs {
			_ = msg.Ack()

			mqMsg := inter.MqMsg{
				Body:      msg.Data,
				Topic:     topic,
				Timestamp: time.Now(),
				RunType:   inter.ReceiveMsg,
			}

			_ = receiveDo(newCtx, mqMsg)

		}
	}

}
