package ecache

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	//blackCache local_cache.Cache
	localCacheMap = gmap.NewStrAnyMap(true) //  *gcache.Cache
)

var (
	cachedb   = "cache"
	cacheType = "file"
	cachePath = "storage/cache"
	//cachepath = "manifest/data/cache/%v.db"
)

func init() {
	ctx := gctx.New()
	get, err := g.Cfg().Get(ctx, "local.cache.name", "")
	if err == nil && !get.IsNil() && !get.IsEmpty() {
		cachedb = get.String()
	}

	get, err = g.Cfg().Get(ctx, "local.cache.type", "")
	if err == nil && !get.IsNil() && !get.IsEmpty() {
		cacheType = get.String()
	}

	get, err = g.Cfg().Get(ctx, "local.cache.filePath", "")
	if err == nil && !get.IsNil() && !get.IsEmpty() {
		cachePath = get.String()
	}
	exists := gfile.Exists(cachePath)
	if !exists {
		gfile.Mkdir(cachePath)
	}

}
