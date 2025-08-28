package nats

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestNet(t *testing.T) {
	nc, err := nats.Connect("nats://192.168.200.26:4222", nats.UserInfo("nats", "PTJK123qwe,.26"))
	if err != nil {
		panic(err)
	}
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "TEST",
		Subjects: []string{"TEST.DATA"},
	})
	if err != nil {
		panic(err)
	}

	for i := 0; i < 1000; i++ {
		publish, err := js.Publish("TEST.DATA", []byte("hello world "+strconv.Itoa(i)))
		if err != nil {
			panic(err)
		}

		fmt.Println(publish)
	}

}

func TestNetC(t *testing.T) {
	nc, err := nats.Connect("nats://192.168.200.26:4222", nats.UserInfo("nats", "PTJK123qwe,.26"))
	if err != nil {
		t.Fatal(err)
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	ctx := gctx.New()
	sub, err := js.PullSubscribe("TEST.DATA", "aaa", nats.Context(ctx), nats.PullMaxWaiting(10))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			//	select {
			//	case <-ctx.Done():
			//		break
			//	default:
			//	}

			msgs, _ := sub.Fetch(10, nats.Context(ctx))
			for _, msg := range msgs {
				msg.Ack()
				fmt.Println(string(msg.Data))
			}
		}

	}()
	//time.Sleep(30 * time.Second)
	//canel()

	time.Sleep(1 * time.Hour)
}
