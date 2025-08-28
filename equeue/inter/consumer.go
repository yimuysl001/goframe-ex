package inter

import "context"

type MqConsumer interface {
	ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg MqMsg) error) (err error)
	Unsubscribe(ctx context.Context, topic string) error
}
