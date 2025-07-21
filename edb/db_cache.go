package edb

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"strings"
	"time"
)

// 设置缓存
func (d *DataBase) SetCache(key string, cachetime time.Duration, nullCache ...bool) *DataBase {
	return &DataBase{
		name:      d.name,
		cacheKey:  fmt.Sprintf(`SelectCache:%s:%s`, d.name, key),
		cachetime: cachetime,
		gdb:       d.gdb,
		nullCache: len(nullCache) > 0 && nullCache[0],
	}

}
func (d *DataBase) Model(tableNameOrStruct ...interface{}) *gdb.Model {
	if d.cachetime < 1 {
		return d.gdb.Model(tableNameOrStruct...)
	}
	return d.gdb.Model(tableNameOrStruct...).Hook(DbHook(d.name, d.nullCache))

}
func (d *DataBase) CtxModel(ctx context.Context, tableNameOrStruct ...interface{}) *gdb.Model {
	if d.cachetime < 1 {
		return d.gdb.Model(tableNameOrStruct...).Ctx(ctx).Safe(true)
	}
	return d.gdb.Model(tableNameOrStruct...).Ctx(ctx).Hook(DbHook(d.name, false)).Safe(true)

}

func (d *DataBase) CacheModel(ctx context.Context, tableNameOrStruct ...interface{}) *gdb.Model {
	return d.SetCache("", time.Hour).gdb.Model(tableNameOrStruct...).Ctx(ctx).
		Hook(DbHook(d.name, d.nullCache)).Safe(true)
}

func (d *DataBase) RemoveCache(ctx context.Context) error {
	keys, err := d.gdb.GetCache().KeyStrings(ctx)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	for _, v := range keys {
		if strings.HasPrefix(v, "SelectCache:") {
			_, err = d.gdb.GetCache().Remove(ctx, v)
			if err != nil {
				return err
			}
		}
	}

	return nil

}
