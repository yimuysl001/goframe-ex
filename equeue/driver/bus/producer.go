package bus

import (
	"context"
	"fmt"
	"goframe-ex/equeue/inter"
	"math/rand"
	"time"
)

// SendByteMsg 生产数据
func (c *LocalEventBus) SendByteMsg(ctx context.Context, topic string, body []byte) (mqMsg inter.MqMsg, err error) {

	mqMsg = inter.MqMsg{
		RunType:   inter.SendMsg,
		Topic:     topic,
		MsgId:     getRandMsgId(),
		Body:      body,
		Timestamp: time.Now(),
	}
	//data, err := json.Marshal(mqMsg)
	//if err != nil {
	//	return
	//}
	//marshal, err := json.Marshal(mqMsg)
	//if err != nil {
	//	return mqMsg, err
	//}

	c.bus.Publish(topic, mqMsg)

	return mqMsg, nil

}

// SendMsg 按字符串类型生产数据
func (c *LocalEventBus) SendMsg(ctx context.Context, topic string, body string) (mqMsg inter.MqMsg, err error) {
	return c.SendByteMsg(ctx, topic, []byte(body))
}

func (c *LocalEventBus) SendDelayMsg(ctx context.Context, topic string, body string, delaySecond int64) (mqMsg inter.MqMsg, err error) {
	mqMsg = inter.MqMsg{
		RunType:   inter.SendMsg,
		Topic:     topic,
		MsgId:     getRandMsgId(),
		Body:      []byte(body),
		Timestamp: time.Now(),
	}
	//data, err := json.Marshal(mqMsg)
	//if err != nil {
	//	return
	//}
	c.bus.Publish(topic, mqMsg)

	return mqMsg, nil
}

func getRandMsgId() string {
	rand.NewSource(time.Now().UnixNano())
	radium := rand.Intn(999) + 1
	timeCode := time.Now().UnixNano()
	return fmt.Sprintf("%d%.4d", timeCode, radium)
}
