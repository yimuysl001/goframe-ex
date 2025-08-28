package equeue

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/equeue/inter"
	"testing"
)

func TestConsumer(t *testing.T) {

	Listen("nats").ListenReceiveMsgDo(gctx.New(), "test", func(ctx context.Context, mqMsg inter.MqMsg) error {
		g.Log().Info(ctx, "mqMsg：", string(mqMsg.Body))
		return nil
	})

}

func TestProducer(t *testing.T) {
	ctx := gctx.New()
	msg, err := Mq("nats").SendMsg(ctx, "test", `{"a":"b"}`)
	if err != nil {
		g.Log().Error(ctx, "mqMsg err：", err)
	}
	g.Log().Info(ctx, "mqMsg：", msg)
	//time.Sleep(2 * time.Hour)
}
