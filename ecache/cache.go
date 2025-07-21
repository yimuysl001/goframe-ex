package ecache

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/ecache/adapter"
	"path"
	"strings"
)

func checkRedis(name string) bool {
	var err2 error = nil

	err := g.Try(gctx.GetInitCtx(), func(ctx context.Context) {
		name = strings.TrimPrefix(name, "db_") // 数据库使用缓存会添加一个 db_
		_, err2 = g.Redis(name).Expire(ctx, "initCheckTest", 1)

	})
	return err == nil && err2 == nil
}

func Cache(pri string, name ...string) *gcache.Cache {
	n := cachedb
	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	return localCacheMap.GetOrSetFunc(n+pri, func() interface{} {
		if checkRedis(n) {
			n = strings.TrimPrefix(n, "db_") // 数据库使用缓存会添加一个 db_
			g.Log().Info(gctx.GetInitCtx(), "使用redis缓存")
			return gcache.NewWithAdapter(adapter.NewAdapterRedis(g.Redis(n), pri))
		}

		switch cacheType {
		case "file":
			return gcache.NewWithAdapter(adapter.NewAdapterFile(path.Join(cachePath, pri)))
		case "leveldb":
			return gcache.NewWithAdapter(adapter.NewLevelDbAdapterFile(path.Join(cachePath, "leveldb", pri)))

		}

		//var localCache *gcache.Cache
		// 有内存泄漏风险
		//err := g.Try(gctx.GetInitCtx(), func(ctx context.Context) {
		//	cpath := fmt.Sprintf(cachepath, n)
		//	localCache = gcache.NewWithAdapter(redka_cache.NewAdapterRedka(cpath))
		//	g.Log().Info(ctx, "使用redka缓存")
		//})
		//if err == nil {
		//	return localCache
		//}

		g.Log().Info(gctx.GetInitCtx(), "使用内存缓存")

		return gcache.New()
	}).(*gcache.Cache)
}
