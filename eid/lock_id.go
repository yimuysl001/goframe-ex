package eid

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"sync"
	"time"
)

var (
	lock       sync.Mutex
	lockMap    = make(map[string]int64)
	getIdUrl   string
	startInt   int64                  = 1000000
	thresholds                        = []int64{10, 100, 1000, 10000, 100000, 1000000}
	startDate  func(key string) int64 = func(key string) int64 {
		return time.Now().Unix() * startInt
	}
)

func SetStartInt(i int64) {
	for _, t := range thresholds {
		if i <= t {
			startInt = t
			return
		}
	}
	startInt = thresholds[len(thresholds)-1] // 超过最大值时取最后一个阈值
}

func SetStartLockIdFunc(startFunc func(key string) int64) {
	startDate = startFunc
}

func GetLockId(key string) int64 {
	lock.Lock()
	defer lock.Unlock()
	u, ok := lockMap[key]
	if !ok {
		u = startDate(key)
	}
	u = u + 1
	lockMap[key] = u

	return u

}

func GetLockIdServer(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"code": 200,
		"data": GetLockId(r.Get("key").String()),
	})
	r.Exit()

}

func SetLockIdUrl(url string) {
	getIdUrl = url
}

func GetLockIdClient(ctx context.Context, key string) int64 {
	if getIdUrl == "" {
		return GetLockId(key)
	}

	get, err := g.Client().Get(ctx, getIdUrl+"?key="+key)
	if err != nil {
		g.Log().Error(ctx, "GetIdClient err:", err)
		return GetLockId(key)
	}
	defer get.Close()

	all := get.ReadAll()

	g2 := gjson.New(all).Get("data")
	if g2.IsNil() {
		return GetLockId(key)
	}
	return g2.Int64()

}
