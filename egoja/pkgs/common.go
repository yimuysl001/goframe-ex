package pkgs

import (
	"context"
	timeconv "github.com/Andrew-M-C/go.timeconv"
	"github.com/gogf/gf/v2/container/gmap"
	"goframe-ex/ecache"
	"goframe-ex/egoja/gojaapi"
	"goframe-ex/eid"
	"time"
)

var localData = gmap.NewStrAnyMap(true)

func init() {
	gojaapi.RegisterImport("common", map[string]any{
		"AddDate": timeconv.AddDate,
		"GetLocalData": func() *gmap.StrAnyMap {
			return localData
		},
		"SetCommon": func(key string, value any) {
			localData.Set(key, value)
		},

		"SetCommonCache": func(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
			return ecache.Cache("gojaCommon:").Set(ctx, key, value, duration)
		},

		"SetCommonCacheByFunc": func(ctx context.Context, key interface{}, f func() interface{}, duration time.Duration) error {
			_, err := ecache.Cache("gojaCommon:").SetIfNotExist(ctx, key, f, duration)
			return err
		},
		"SetCommonByFunc": func(key string, f func() interface{}) {
			localData.GetOrSetFuncLock(key, f)
		},
		"GetCommon": func(key string) (any, bool) {
			return localData.Search(key)
		},

		"GetCommonCache": func(ctx context.Context, key string) (any, bool) {
			data, err := ecache.Cache("gojaCommon:").Get(ctx, key)
			if err != nil || data == nil {
				return nil, false
			}
			return data, true
		},

		// 雪花id
		"GenSnowID":    eid.GenSnowID,
		"GenSnowIDStr": eid.GenSnowIDStr,
		"GetUlid":      eid.GetUlid,
	})
}
