package einit

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
)

const cronName = "libInitCron"

func StartCron(p Pri) {
	g.Log().Info(gctx.GetInitCtx(), "注册时间：", p.BootTime)
	g.Log().Info(gctx.GetInitCtx(), "程序有效期：", p.LastTime)
	check(gctx.GetInitCtx(), p)
	_, err := gcron.AddSingleton(gctx.GetInitCtx(), "@every 1m", func(ctx context.Context) {
		check(ctx, p)
	}, cronName)

	if err != nil {
		panic(err)
	}

}
