package equeue

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/equeue/inter"
)

const (
	defName = "default"
)

var (
	consumerMap = gmap.NewStrAnyMap(true)
	producerMap = gmap.NewStrAnyMap(true)

	mqConfig = make(map[string]inter.MqConfig)
)

func init() {
	get, err := g.Cfg().Get(gctx.GetInitCtx(), "equeue")
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), "get equeue config error:"+err.Error())
	}
	err = get.Scan(&mqConfig)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), "read equeue config error:"+err.Error())
	}

	if len(mqConfig) < 1 {
		g.Log().Debug(gctx.GetInitCtx(), "mq config is null")
	}
}
