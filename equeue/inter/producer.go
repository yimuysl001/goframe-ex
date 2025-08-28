package inter

import "context"

type MqProducer interface {
	SendMsg(ctx context.Context, topic string, body string) (mqMsg MqMsg, err error)
	SendByteMsg(ctx context.Context, topic string, body []byte) (mqMsg MqMsg, err error)
	SendDelayMsg(ctx context.Context, topic string, body string, delaySecond int64) (mqMsg MqMsg, err error)
}
