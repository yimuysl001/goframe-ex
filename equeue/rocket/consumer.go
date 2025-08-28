package rocket

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gogf/gf/v2/errors/gerror"
	"goframe-ex/equeue/inter"
	"goframe-ex/equeue/logger"
)

// ListenReceiveMsgDo 消费数据
func (r *RocketMq) ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg inter.MqMsg) error) (err error) {
	if r.consumerIns == nil {
		return gerror.New("rocketMq consumer not register")
	}

	err = r.consumerIns.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, item := range msgs {
			err = receiveDo(ctx, inter.MqMsg{
				RunType: inter.ReceiveMsg,
				Topic:   item.Topic,
				MsgId:   item.MsgId,
				Body:    item.Body,
			})
			if err != nil {
				logger.Logger().Error(ctx, "MQ执行报错：", err)
				return consumer.Rollback, err
			}
		}
		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		return
	}

	if err = r.consumerIns.Start(); err != nil {
		_ = r.consumerIns.Unsubscribe(topic)
		return
	}
	return
}
func (r *RocketMq) Unsubscribe(ctx context.Context, topic string) (err error) {
	if r.consumerIns == nil {
		return r.consumerIns.Unsubscribe(topic)
	}
	return err

}
