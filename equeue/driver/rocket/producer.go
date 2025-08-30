package rocket

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gogf/gf/v2/errors/gerror"
	"goframe-ex/equeue/inter"
)

// SendMsg 按字符串类型生产数据
func (r *RocketMq) SendMsg(ctx context.Context, topic string, body string) (mqMsg inter.MqMsg, err error) {
	return r.SendByteMsg(ctx, topic, []byte(body))
}

// SendByteMsg 生产数据
func (r *RocketMq) SendByteMsg(ctx context.Context, topic string, body []byte) (mqMsg inter.MqMsg, err error) {
	if r.producerIns == nil {
		return mqMsg, gerror.New("rocketMq producer not register")
	}

	result, err := r.producerIns.SendSync(ctx, &primitive.Message{
		Topic: topic,
		Body:  body,
	})

	if err != nil {
		return
	}
	if result.Status != primitive.SendOK {
		return mqMsg, gerror.Newf("rocketMq producer send msg error status:%v", result.Status)
	}

	mqMsg = inter.MqMsg{
		RunType: inter.SendMsg,
		Topic:   topic,
		MsgId:   result.MsgID,
		Body:    body,
	}
	return mqMsg, nil
}

func (r *RocketMq) SendDelayMsg(ctx context.Context, topic string, body string, delaySecond int64) (mqMsg inter.MqMsg, err error) {
	err = gerror.New("implement me")
	return
}
