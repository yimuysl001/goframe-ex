package equeue

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/equeue/inter"
	"testing"
)

func TestConsumer(t *testing.T) {

	err := Listen("rabbit").ListenReceiveMsgDo(gctx.New(), "test", func(ctx context.Context, mqMsg inter.MqMsg) error {
		g.Log().Info(ctx, "mqMsg：", string(mqMsg.Body))
		return nil
	})

	fmt.Println(err)
}

func TestProducer(t *testing.T) {
	ctx := gctx.New()
	for i := 0; i < 100; i++ {
		msg, err := Mq("rabbit").SendMsg(ctx, "test", `{"a":"b"}`)
		if err != nil {
			g.Log().Error(ctx, "mqMsg err：", err)
		}
		g.Log().Info(ctx, "mqMsg：", msg)
	}

	//time.Sleep(2 * time.Hour)
}
