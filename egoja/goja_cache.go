package egoja

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/gogf/gf/v2/container/gmap"
	"goframe-ex/egoja/gojaapi"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	localCacheFunc = gmap.NewStrAnyMap(true)
	single         singleflight.Group
)

func FlushCache(id string, script string, f ...func(string, string) (string, string)) (*goja.Program, error) {
	if script == "" {
		localCacheFunc.Remove(id)
		return nil, nil
	}
	refresh, s, err := gojaapi.ScriptRefresh(script)
	if err != nil {
		return nil, err
	}

	if f != nil && len(f) > 0 {
		for _, f2 := range f {
			refresh, s = f2(refresh, s)
		}
	}
	script = refresh + fmt.Sprintf(mainFunc, s)
	prog, err := goja.Compile(id, script, false)
	if err != nil {
		return nil, err
	}
	localCacheFunc.Set(id, prog)
	return prog, nil
}

func GetCacheProgram(id, script string, f ...func(string, string) (string, string)) (*goja.Program, error) {
	pr, err2, _ := single.Do(id, func() (interface{}, error) {
		value, found := localCacheFunc.Search(id)
		if found && value != nil {
			return value, nil
		}
		return FlushCache(id, script, f...)
	})
	if err2 != nil {
		return nil, err2
	}
	return pr.(*goja.Program), nil
}

func ExecScriptById(ctx context.Context, id, script string, params map[string]any, d ...time.Duration) (goja.Value, error) {

	prog, err := GetCacheProgram(id, script)
	if err != nil {
		return nil, err
	}

	vm := gojaapi.GetNewVm()
	defer gojaapi.PutGoja(vm)
	for k, v := range params {
		vm.Set(k, v)
	}
	vm.Set("ctx", ctx)

	if len(d) > 0 && d[0] > 0 {
		time.AfterFunc(d[0], func() {
			vm.Interrupt("time out")
		})
	}

	return vm.RunProgram(prog)
}

func ExecScriptConversionTimeOutById(ctx context.Context, id, script string, params map[string]any, d time.Duration, f ...func(string, string) (string, string)) (goja.Value, error) {
	prog, err := GetCacheProgram(id, script, f...)
	if err != nil {
		return nil, err
	}

	vm := gojaapi.GetNewVm()
	defer gojaapi.PutGoja(vm)
	for k, v := range params {
		vm.Set(k, v)
	}
	vm.Set("ctx", ctx)

	if d > 0 {
		time.AfterFunc(d, func() {
			vm.Interrupt("time out")
		})
	}
	return vm.RunProgram(prog)
}

func ExecScriptConversionById(ctx context.Context, id, script string, params map[string]any, f ...func(string, string) (string, string)) (goja.Value, error) {
	return ExecScriptConversionTimeOutById(ctx, id, script, params, 0, f...)
}
