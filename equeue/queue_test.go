package equeue

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/equeue/inter"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {

	err := Listen("nats").ListenReceiveMsgDo(gctx.New(), "test", func(ctx context.Context, mqMsg inter.MqMsg) error {
		g.Log().Info(ctx, "mqMsg：", string(mqMsg.Body))
		return nil
	})

	fmt.Println(err)

}

func TestProducer(t *testing.T) {

	ctx := gctx.New()
	for i := 0; i < 100; i++ {
		msg, err := Mq("nats").SendMsg(ctx, "test", `{"a":"b"}`)
		if err != nil {
			g.Log().Error(ctx, "mqMsg err：", err)
		}
		g.Log().Info(ctx, "mqMsg：", msg)
	}

	//time.Sleep(2 * time.Hour)
}

func TestBus(t *testing.T) {
	SetConfig("bus", inter.MqConfig{Driver: "bus"})
	go func() {
		err := Listen("bus").ListenReceiveMsgDo(gctx.New(), "test", func(ctx context.Context, mqMsg inter.MqMsg) error {
			time.Sleep(1 * time.Second)
			g.Log().Info(ctx, "mqMsg：", string(mqMsg.Body))
			return nil
		})

		fmt.Println(err)
	}()
	time.Sleep(1 * time.Second)
	ctx := gctx.New()
	for i := 0; i < 100; i++ {
		msg, err := Mq("bus").SendMsg(ctx, "test", `{"a":"b"}`)
		if err != nil {
			g.Log().Error(ctx, "mqMsg err：", err)
		}
		g.Log().Info(ctx, "mqMsg：", msg)
	}

	time.Sleep(2 * time.Hour)
}
