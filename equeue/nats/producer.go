package nats

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"goframe-ex/equeue/inter"
	"math/rand"
	"time"
)

// SendByteMsg 生产数据
func (r *NatsMq) SendByteMsg(ctx context.Context, topic string, body []byte) (mqMsg inter.MqMsg, err error) {

	//mqMsg = inter.MqMsg{
	//	RunType:   inter.SendMsg,
	//	Topic:     topic,
	//	MsgId:     getRandMsgId(),
	//	Body:      body,
	//	Timestamp: time.Now(),
	//}
	//data, err := json.Marshal(mqMsg)
	//if err != nil {
	//	return
	//}
	_, err = r.js.Publish(r.subjects+"."+topic, body, nats.Context(ctx))
	if err != nil {
		return mqMsg, err
	}

	return mqMsg, nil

}

// SendMsg 按字符串类型生产数据
func (r *NatsMq) SendMsg(ctx context.Context, topic string, body string) (mqMsg inter.MqMsg, err error) {
	return r.SendByteMsg(ctx, topic, []byte(body))
}

func (r *NatsMq) SendDelayMsg(ctx context.Context, topic string, body string, delaySecond int64) (mqMsg inter.MqMsg, err error) {
	//mqMsg = inter.MqMsg{
	//	RunType:   inter.SendMsg,
	//	Topic:     topic,
	//	MsgId:     getRandMsgId(),
	//	Body:      []byte(body),
	//	Timestamp: time.Now(),
	//}
	//data, err := json.Marshal(mqMsg)
	//if err != nil {
	//	return
	//}
	_, err = r.js.PublishAsync(r.subjects+"."+topic, []byte(body), nats.Context(ctx))
	if err != nil {
		return mqMsg, err
	}

	return mqMsg, nil
}

func getRandMsgId() string {
	rand.NewSource(time.Now().UnixNano())
	radium := rand.Intn(999) + 1
	timeCode := time.Now().UnixNano()
	return fmt.Sprintf("%d%.4d", timeCode, radium)
}
