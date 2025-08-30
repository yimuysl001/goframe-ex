package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/gogf/gf/v2/errors/gerror"
	"goframe-ex/equeue/inter"
	"time"
)

// SendByteMsg 生产数据
func (r *KafkaMq) SendByteMsg(ctx context.Context, topic string, body []byte) (mqMsg inter.MqMsg, err error) {

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(body),
		Timestamp: time.Now(),
	}
	if r.producerIns == nil {
		err = gerror.New("queue kafka producerIns is nil")
		return
	}

	r.producerIns.Input() <- msg
	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case info := <-r.producerIns.Successes():
		return inter.MqMsg{
			RunType:   inter.SendMsg,
			Topic:     info.Topic,
			Offset:    info.Offset,
			Partition: info.Partition,
			Timestamp: info.Timestamp,
		}, nil
	case fail := <-r.producerIns.Errors():
		if nil != fail {
			return mqMsg, fail.Err
		}
	case <-sendCtx.Done():
		return mqMsg, gerror.New("send mqMst timeout")
	}
	return mqMsg, nil

}

// SendMsg 按字符串类型生产数据
func (r *KafkaMq) SendMsg(ctx context.Context, topic string, body string) (mqMsg inter.MqMsg, err error) {
	return r.SendByteMsg(ctx, topic, []byte(body))
}

func (r *KafkaMq) SendDelayMsg(ctx context.Context, topic string, body string, delaySecond int64) (mqMsg inter.MqMsg, err error) {

	err = gerror.New("implement me")
	return
}
