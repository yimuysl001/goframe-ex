package main

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"goframe-ex/econf"
	"goframe-ex/ecron"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	econf.InitConf()
	var ctx = gctx.New()

	paths := g.Cfg().MustGet(ctx, "clearLog.paths", "").Strings()

	cron := g.Cfg().MustGet(ctx, "clearLog.cron", "0 0 2 * * *").String()

	saveDay := g.Cfg().MustGet(ctx, "clearLog.day", "30").Int()

	g.Log().Info(ctx, "paths:", paths)
	g.Log().Info(ctx, "cron:", cron)
	g.Log().Info(ctx, "saveDay:", saveDay)
	clearLog(ctx, paths, saveDay)
	_, err := ecron.AddSingleton(ctx, cron, func(ctx context.Context) {
		clearLog(ctx, paths, saveDay)

	})
	g.Log().Info(ctx, "AddSingleton:", err)
	select {}

}

func clearLog(ctx context.Context, paths []string, saveDay int) {
	g.Log().Info(ctx, "开始清理日志")
	err := g.Try(ctx, func(ctx context.Context) {
		date := time.Now().AddDate(0, 0, -1*saveDay)
		g.Log().Info(ctx, "clear date:", date)
		for _, p := range paths {
			g.Log().Info(ctx, "path:", p)
			dirs, err := os.ReadDir(p)
			if err != nil {
				g.Log().Info(ctx, err)
				continue
			}
			for _, dir := range dirs {
				info, err := dir.Info()
				if err != nil {
					g.Log().Info(ctx, err)
					continue
				}

				g.Log().Info(ctx, "Name:", info.Name())
				g.Log().Info(ctx, "ModTime:", info.ModTime())
				if (strings.HasSuffix(info.Name(), ".log") || strings.HasSuffix(info.Name(), ".gz")) && info.ModTime().Unix() < date.Unix() {
					g.Log().Info(ctx, "需要清理的数据", path.Join(p, info.Name()))
					if info.IsDir() {
						err = gfile.RemoveAll(path.Join(p, info.Name()))
					} else {
						err = gfile.RemoveFile(path.Join(p, info.Name()))
					}
					g.Log().Info(ctx, "删除情况：", err)
				}

			}

		}

	})
	if err != nil {
		g.Log().Info(ctx, "开始清理出错：", err)
	}
	g.Log().Info(ctx, "开始清理完成")
}
