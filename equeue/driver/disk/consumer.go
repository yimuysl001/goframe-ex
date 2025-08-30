package disk

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/errors/gerror"
	"goframe-ex/equeue/driver/logger"
	"goframe-ex/equeue/inter"
	"time"
)

func RegisterDiskMqConsumer(config inter.MqConfig) (client inter.MqConsumer, err error) {
	return &DiskConsumerMq{
		config: config,
		flag:   make(map[string]bool),
		queue:  make(map[string]*Queue),
	}, nil
}

// ListenReceiveMsgDo 消费数据
func (q *DiskConsumerMq) ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg inter.MqMsg) error) (err error) {
	if topic == "" {
		return gerror.New("disk.ListenReceiveMsgDo topic is empty")
	}
	if q.flag[topic] == true {
		return gerror.New("disk.ListenReceiveMsgDo topic already in use")
	}
	var queue = NewDiskQueue(topic, q.config)
	q.queue[topic] = queue
	q.flag[topic] = true
	var (
		sleep = time.Second
	)

	for {
		if q.flag[topic] == false {
			break
		}

		if index, offset, data, err := queue.Read(); err == nil {
			var mqMsg inter.MqMsg
			if err = json.Unmarshal(data, &mqMsg); err != nil {
				logger.Logger().Warningf(ctx, "disk.ListenReceiveMsgDo Unmarshal err:%+v, topic：%v, data:%+v .", err, topic, string(data))
				continue
			}
			if mqMsg.MsgId != "" {
				if err = receiveDo(ctx, mqMsg); err != nil {
					logger.Logger().Error(ctx, "disk.ListenReceiveMsgDo receiveDo err:%+v", err)
				}
				queue.Commit(index, offset)
				sleep = time.Millisecond * 10
			}
		} else {
			sleep = time.Second
		}

		time.Sleep(sleep)
	}
	return
}

func (q *DiskConsumerMq) Unsubscribe(ctx context.Context, topic string) (err error) {
	q.flag[topic] = false

	if q.queue[topic] == nil {
		q.queue[topic].Close()
	}
	return err

}
