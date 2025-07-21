package edb

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
	"goframe-ex/ecache"
	"path"
	"sync"
	"time"
)

var (
	localDb  = gmap.NewStrAnyMap(true)
	one      = sync.Once{}
	initflag = false
)

type DataBase struct {
	gdb       gdb.DB
	name      string
	cacheKey  string
	nullCache bool
	cachetime time.Duration
}

func DB(name ...string) *DataBase {
	n := gdb.DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	setFunc := localDb.GetOrSetFuncLock(n, func() interface{} {
		db, _ := gdb.NewByGroup(n)
		if db == nil {
			db = g.DB(n)
		}

		// 添加缓存
		db.GetCache().SetAdapter(ecache.Cache("db_" + n))

		// 添加日志
		get, _ := g.Cfg().Get(context.TODO(), "database.logger")
		getmap := get.Map()
		if p, ok := getmap["path"]; !ok {
			getmap["path"] = "logs/sql/" + n
		} else {
			getmap["path"] = path.Join(gconv.String(p), n)
		}
		var lglog = glog.New().Clone()

		_ = lglog.SetConfigWithMap(getmap)
		db.SetLogger(lglog)

		return &DataBase{gdb: db, name: n}

	})

	if setFunc == nil {
		panic(n + ":db fail")
	}

	return setFunc.(*DataBase)

}

func (d *DataBase) GDB() gdb.DB {
	return d.gdb
}

func CheckDB(name ...string) (check bool) {

	err := g.Try(gctx.GetInitCtx(), func(ctx context.Context) {
		err := DB(name...).GDB().PingMaster()
		check = err == nil
	})

	if err != nil {
		check = false
	}

	return check

}

func RemoveDb(key string) {
	if key == "" {
		key = gdb.DefaultGroupName
	}

	localDb.Remove(key)
}

func SetDb(group string, nodes gdb.ConfigGroup) error {
	return g.Try(context.TODO(), func(ctx context.Context) {
		err := gdb.SetConfigGroup(group, nodes)
		if err != nil {
			panic(err)
		}
		RemoveDb(group)
		_ = DB(group)

	})

}

func CheckInit() bool {
	one.Do(func() {
		initflag = CheckDB()
	})
	return initflag
}
