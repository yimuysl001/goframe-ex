package redis

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/equeue/inter"
)

type RedisMq struct {
	poolName  string
	groupName string
	timeout   int64
	flag      map[string]bool
}

func RegisterRedis(config inter.MqConfig) (client *RedisMq, err error) {
	err = g.Try(gctx.GetInitCtx(), func(ctx context.Context) {
		g.Redis(config.Name)
	})
	if err != nil {
		return nil, err
	}

	return &RedisMq{
		poolName:  config.Name,
		groupName: config.GroupName,
		timeout:   config.Timeout,
		flag:      make(map[string]bool),
	}, nil
}
