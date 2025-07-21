package edb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dromara/carbon/v2"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/ghash"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/sync/singleflight"
	"strings"
	"time"
)

const (
	preKey       = "myLocalKey"
	preCache     = "myLocalCache"
	separatorStr = ":"
)

var (
	keyCache *gcache.Cache = nil
	single   singleflight.Group
)

func getCacheKey(ctx context.Context, name string, in *gdb.HookSelectInput) string {
	var tableName = in.Table
	if strings.Contains(tableName, " AS ") {
		tableName = strings.SplitN(tableName, " AS ", 2)[0]
	}

	var localCache = preCache + separatorStr + name + separatorStr + in.Schema + separatorStr + tableName + separatorStr + gmd5.MustEncrypt(in.Sql+"@"+gconv.String(in.Args))

	err := grpool.Add(gctx.New(), func(ctx context.Context) {
		var localKey = preKey + separatorStr + name + separatorStr + in.Schema + separatorStr + tableName + separatorStr
		get, err := keyCache.Get(ctx, localKey)
		if err != nil || get == nil {
			err = keyCache.Set(ctx, localKey, []string{localCache}, 0)
			if err != nil {
				g.Log().Error(ctx, "keyCache.Set ", err)
				return
			}
			return
		}

		interfaces := get.Interfaces()

		var b = false
		for _, i := range interfaces {
			if i == localKey {
				b = true
				break
			}
		}
		if b {
			return
		}

		interfaces = append(interfaces, localCache)
		err = keyCache.Set(ctx, localKey, interfaces, 0)
		if err != nil {
			g.Log().Error(ctx, "keyCache.Set ", err)
			return
		}
		return
	})

	if err != nil {
		g.Log().Error(ctx, err)
	}
	g.Log().Info(ctx, "数据key:", localCache)
	return localCache

}

func genSelectCacheKey(table, group, schema, name, sql string, args ...interface{}) string {

	name = fmt.Sprintf(
		`%s@%s#%s:%d`,
		guessPrimaryTableName(group, table),
		group,
		schema,
		ghash.BKDR64([]byte(sql+", @PARAMS:"+gconv.String(args))),
	)

	return fmt.Sprintf(`%s%s`, `SelectCache:`, name)
}

func guessPrimaryTableName(name, tableStr string) string {
	if tableStr == "" {
		return ""
	}
	var (
		guessedTableName string
		array1           = gstr.SplitAndTrim(tableStr, ",")
		array2           = gstr.SplitAndTrim(array1[0], " ")
		array3           = gstr.SplitAndTrim(array2[0], ".")
	)
	if len(array3) >= 2 {
		guessedTableName = array3[1]
	} else {
		guessedTableName = array3[0]
	}
	charL, charR := DB(name).GDB().GetChars()
	if charL != "" || charR != "" {
		guessedTableName = gstr.Trim(guessedTableName, charL+charR)
	}
	if !gregex.IsMatchString(`^[\w\.\-]+$`, guessedTableName) {
		return ""
	}
	return guessedTableName
}
func DbHook(name string, nullCache ...bool) gdb.HookHandler {
	fprce := false
	if len(nullCache) > 0 {
		fprce = nullCache[0]
	}
	cachetime := DB(name).cachetime

	return gdb.HookHandler{
		Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
			if cachetime < 0 {
				return in.Next(ctx)
			}

			var ckey = genSelectCacheKey(in.Table, name, in.Schema, "", in.Sql, in.Args...)
			get, err := DB(name).GDB().GetCache().Get(ctx, ckey)
			if err == nil && get != nil {
				err = get.Scan(&result)
				return
			}
			SetArgs(in)

			// 添加防止缓存穿透
			v, err, _ := single.Do(ckey, func() (interface{}, error) {
				return in.Next(ctx)
			})
			if err != nil {
				return result, err
			}
			if !fprce && (v == nil || len(v.(gdb.Result)) < 1) {
				return result, err
			}
			_ = DB(name).GDB().GetCache().Set(ctx, ckey, v, DB(name).cachetime)

			return v.(gdb.Result), nil
		},
		Update: func(ctx context.Context, in *gdb.HookUpdateInput) (result sql.Result, err error) {

			next, err := in.Next(ctx)
			if cachetime < 0 {
				return next, err
			}
			if err == nil {
				err2 := DB(name).GDB().GetCore().ClearCache(ctx, guessPrimaryTableName(name, in.Table))
				if err2 != nil {
					g.Log().Debug(ctx, "update ClearCache err:", err2)
				}
			}
			return next, err
		},
		Insert: func(ctx context.Context, in *gdb.HookInsertInput) (res sql.Result, err error) {
			next, err := in.Next(ctx)
			if cachetime < 0 {
				return next, err
			}
			if err == nil {
				err2 := DB(name).GDB().GetCore().ClearCache(ctx, guessPrimaryTableName(name, in.Table))
				if err2 != nil {
					g.Log().Debug(ctx, "update ClearCache err:", err2)
				}
			}
			return next, err
		},
		Delete: func(ctx context.Context, in *gdb.HookDeleteInput) (res sql.Result, err error) {

			next, err := in.Next(ctx)
			if cachetime < 0 {
				return next, err
			}
			if err == nil {
				err2 := DB(name).GDB().GetCore().ClearCache(ctx, guessPrimaryTableName(name, in.Table))
				if err2 != nil {
					g.Log().Debug(ctx, "update ClearCache err:", err2)
				}
			}
			return next, err
		},
	}
}

func SetArgs(in *gdb.HookSelectInput) {
	// 时间格式处理
	var args = make([]any, len(in.Args))
	for i, arg := range in.Args {
		switch t := arg.(type) {
		case time.Time:
			args[i] = carbon.NewCarbon(t)
		case *time.Time:
			args[i] = carbon.NewCarbon(*t)
		case *gtime.Time:
			args[i] = carbon.NewCarbon(t.Time)
		case gtime.Time:
			args[i] = carbon.NewCarbon(t.Time)
		default:
			args[i] = t
		}
	}
	in.Args = args
}
