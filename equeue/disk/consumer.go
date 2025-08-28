package disk

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/errors/gerror"
	"goframe-ex/equeue/inter"
	"goframe-ex/equeue/logger"
	"time"
)

func RegisterDiskMqConsumer(config inter.MqConfig) (client inter.MqConsumer, err error) {
	return &DiskConsumerMq{
		config: config,
	}, nil
}

// ListenReceiveMsgDo 消费数据
func (q *DiskConsumerMq) ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg inter.MqMsg) error) (err error) {
	if topic == "" {
		return gerror.New("disk.ListenReceiveMsgDo topic is empty")
	}
	q.queue = NewDiskQueue(topic, q.config)
	q.flag = true
	var (
		sleep = time.Second
	)

	for {
		if q.flag == false {
			break
		}

		if index, offset, data, err := q.queue.Read(); err == nil {
			var mqMsg inter.MqMsg
			if err = json.Unmarshal(data, &mqMsg); err != nil {
				logger.Logger().Warningf(ctx, "disk.ListenReceiveMsgDo Unmarshal err:%+v, topic：%v, data:%+v .", err, topic, string(data))
				continue
			}
			if mqMsg.MsgId != "" {
				if err = receiveDo(ctx, mqMsg); err != nil {
					logger.Logger().Error(ctx, "disk.ListenReceiveMsgDo receiveDo err:%+v", err)
				}
				q.queue.Commit(index, offset)
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
	q.flag = false

	if q.queue == nil {
		q.queue.Close()
	}
	return err

}
