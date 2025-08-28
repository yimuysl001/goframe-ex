package disk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"goframe-ex/equeue/inter"
	"math/rand"
	"time"
)

func RegisterDiskMqProducer(config inter.MqConfig) (client inter.MqProducer, err error) {
	return &DiskProducerMq{
		config:    config,
		producers: make(map[string]*Queue),
	}, nil
}

// SendMsg 按字符串类型生产数据
func (d *DiskProducerMq) SendMsg(ctx context.Context, topic string, body string) (mqMsg inter.MqMsg, err error) {
	return d.SendByteMsg(ctx, topic, []byte(body))
}

// SendByteMsg 生产数据
func (d *DiskProducerMq) SendByteMsg(ctx context.Context, topic string, body []byte) (mqMsg inter.MqMsg, err error) {
	if topic == "" {
		return mqMsg, gerror.New("DiskMq topic is empty")
	}

	mqMsg = inter.MqMsg{
		RunType:   inter.SendMsg,
		Topic:     topic,
		MsgId:     getRandMsgId(),
		Body:      body,
		Timestamp: time.Now(),
	}

	mqMsgJson, err := json.Marshal(mqMsg)
	if err != nil {
		return mqMsg, gerror.New(fmt.Sprint("queue redis 生产者解析json消息失败:", err))
	}

	queue := d.getProducer(topic)
	if err = queue.Write(mqMsgJson); err != nil {
		return mqMsg, gerror.New(fmt.Sprint("queue disk 生产者添加消息失败:", err))
	}
	return
}

func (d *DiskProducerMq) SendDelayMsg(ctx context.Context, topic string, body string, delaySecond int64) (mqMsg inter.MqMsg, err error) {
	err = gerror.New("implement me")
	return
}

func (d *DiskProducerMq) getProducer(topic string) *Queue {
	queue, ok := d.producers[topic]
	if ok {
		return queue
	}
	queue = NewDiskQueue(topic, d.config)
	d.Lock()
	defer d.Unlock()
	d.producers[topic] = queue
	return queue
}

func getRandMsgId() string {
	rand.NewSource(time.Now().UnixNano())
	radium := rand.Intn(999) + 1
	timeCode := time.Now().UnixNano()
	return fmt.Sprintf("%d%.4d", timeCode, radium)
}
