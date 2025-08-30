package nats

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"goframe-ex/equeue/inter"
	"strings"
)

type NatsMq struct {
	js         nats.JetStreamContext
	subjects   string
	durable    string
	maxWaiting int
	batch      int

	cancelMap map[string]context.CancelFunc
}

func RegisterNats(config inter.MqConfig) (client *NatsMq, err error) {
	if len(config.Address) == 0 {
		return nil, errors.New("address 地址不能为空")
	}
	if len(config.GroupName) == 0 {
		return nil, errors.New("groupName 不能为空")
	}

	client = &NatsMq{
		subjects:   config.GroupName,
		durable:    config.Name,
		cancelMap:  make(map[string]context.CancelFunc),
		maxWaiting: int(config.SegmentLimit),
		batch:      int(config.BatchSize),
	}
	if client.maxWaiting == 0 {
		client.maxWaiting = 128
	}

	if client.batch == 0 {
		client.batch = 10
	}

	var urls = make([]string, len(config.Address))
	for _, addr := range config.Address {
		urls = append(urls, "nats://"+addr)
	}

	nc, err := nats.Connect(strings.Join(urls, ","), nats.UserInfo(config.UserName, config.Password))
	if err != nil {
		return nil, err
	}
	client.js, err = nc.JetStream()
	if err != nil {
		return nil, err
	}
	_, err = client.js.AddStream(&nats.StreamConfig{
		Name:     client.subjects,
		Subjects: []string{client.subjects + ".*"},
	})

	return client, err
}
