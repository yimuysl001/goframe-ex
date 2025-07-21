/**
cron 添加有效期处理
*/

package ecron

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gtimer"
	"time"
)

type (
	JobFunc   = gtimer.JobFunc
	CheckFunc = func(ctx context.Context) error
)

var (
	localCheckFunc    CheckFunc = nil
	localCheckFuncMap           = make(map[string]CheckFunc)
)

func SetCheckFuncMap(name string, f CheckFunc) {
	localCheckFuncMap[name] = f
}

func SetCheckFunc(f CheckFunc) {
	localCheckFunc = f
}

func checkCron(ctx context.Context, name ...string) bool {
	if localCheckFunc != nil {
		err := localCheckFunc(ctx)
		if err != nil {
			g.Log().Error(ctx, "cron check job failed:", err)
			return false
		}
	}

	if len(name) < 2 || name[1] == "" {
		return true
	}

	checkFunc, ok := localCheckFuncMap[name[1]]
	if !ok {
		return true
	}
	err := checkFunc(ctx)
	if err != nil {
		g.Log().Error(ctx, "cron SetCheckFuncMap job failed:", err)
		return false
	}

	return true
}
func Add(ctx context.Context, pattern string, job JobFunc, name ...string) (*gcron.Entry, error) {
	return gcron.Add(ctx, pattern, func(ctx context.Context) {
		if !checkCron(ctx, name...) {
			return
		}
		job(ctx)

	}, name...)
}
func AddSingleton(ctx context.Context, pattern string, job JobFunc, name ...string) (*gcron.Entry, error) {
	return gcron.AddSingleton(ctx, pattern, func(ctx context.Context) {
		if !checkCron(ctx, name...) {
			return
		}
		job(ctx)

	}, name...)
}
func AddOnce(ctx context.Context, pattern string, job JobFunc, name ...string) (*gcron.Entry, error) {
	return gcron.AddOnce(ctx, pattern, func(ctx context.Context) {
		if !checkCron(ctx, name...) {
			return
		}
		job(ctx)

	}, name...)
}
func AddTimes(ctx context.Context, pattern string, times int, job JobFunc, name ...string) (*gcron.Entry, error) {
	return gcron.AddTimes(ctx, pattern, times, func(ctx context.Context) {
		if !checkCron(ctx, name...) {
			return
		}
		job(ctx)

	}, name...)
}

// DelayAdd adds a timed task to default cron object after `delay` time.
func DelayAdd(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	gcron.DelayAdd(ctx, delay, pattern, func(ctx context.Context) {
		if !checkCron(ctx, name...) {
			return
		}
		job(ctx)

	}, name...)
}

// DelayAddSingleton adds a singleton timed task after `delay` time to default cron object.
func DelayAddSingleton(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	gcron.DelayAddSingleton(ctx, delay, pattern, func(ctx context.Context) {
		if !checkCron(ctx, name...) {
			return
		}
		job(ctx)

	}, name...)
}

// DelayAddOnce adds a timed task after `delay` time to default cron object.
// This timed task can be run only once.
func DelayAddOnce(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	gcron.DelayAddOnce(ctx, delay, pattern, func(ctx context.Context) {
		if !checkCron(ctx, name...) {
			return
		}
		job(ctx)

	}, name...)
}

// DelayAddTimes adds a timed task after `delay` time to default cron object.
// This timed task can be run specified times.
func DelayAddTimes(ctx context.Context, delay time.Duration, pattern string, times int, job JobFunc, name ...string) {
	gcron.DelayAddTimes(ctx, delay, pattern, times, func(ctx context.Context) {
		if !checkCron(ctx, name...) {
			return
		}
		job(ctx)

	}, name...)
}
