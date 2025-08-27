package gojaapi

import (
	"github.com/dop251/goja"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/sync/singleflight"
	"strings"
	"sync"
)

type FuncStruct struct {
	//V      goja.Value
	//R      *goja.Runtime
	v      map[string]any
	script string
}

var (
	lock   sync.RWMutex
	single singleflight.Group
	// 注册go函数
	localImport = make(map[string]map[string]any)
	// 注册脚本函数
	localFunc = gmap.NewStrStrMap(true) //  gmap.NewStrAnyMap(true) //  make(map[string]FuncStruct)
	//localFuncStru = gmap.NewStrAnyMap(true)
	//
	scriptCacheFunc = make(map[string]*goja.Program)
	vmPool          = sync.Pool{
		New: func() interface{} {
			//vm := goja.New()
			//vm.Set("imports", Import)
			//vm.Set("importFunc", ImportFunc)
			//vm.Set("db", dbutil.DB)

			return NewVm()
		},
	}
	localCommonParameter = make(map[string]any)
)

func NewVm() *goja.Runtime {
	vm := goja.New()
	vm.Set("imports", Import)
	vm.Set("importFunc", ImportFunc)
	vm.Set("glogs", g.Log)
	vm.Set("glog", g.Log())
	//vm.Set("db", dbutil.DB)
	for k, v := range localCommonParameter {
		vm.Set(k, v)

	}

	return vm
}

// 注册通用变量
func RegisterCommonParameter(name string, value any) {
	localCommonParameter[name] = value
}

// 注册程序中的函数方法
func RegisterImport(name string, value map[string]any) {
	m, ok := localImport[name]
	if !ok {
		localImport[name] = value
		return
	}
	for k, v := range value {
		m[k] = v
	}
	localImport[name] = m

}
func RegisterFunc(name string, script string) error {

	refresh, s, err := ScriptRefresh(script)
	if err != nil {
		return err
	}

	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "const ") {
		n := strings.SplitN(s, "=", 2)
		if len(n) == 2 {
			var v = strings.Replace(n[0], "const ", "", -1)
			s = "const " + v + " = " + n[1] + "\n" + v
		}
	}

	localFunc.Set(name, refresh+"\n"+s)
	delete(scriptCacheFunc, name)
	return nil
}
func Import(name string) map[string]any {
	//lock.Lock()
	//defer lock.Unlock()
	return localImport[name]
}
func getCompiledScript(name, code string) (*goja.Program, error) {
	if prog, ok := scriptCacheFunc[name]; ok {
		return prog, nil
	}

	prog, err := goja.Compile(name, code, false)
	if err != nil {
		return nil, err
	}

	scriptCacheFunc[name] = prog
	return prog, nil
}
func ImportFunc(name string) any {
	lock.Lock()
	defer lock.Unlock()

	e, ok := localFunc.Search(name) //[replace]
	if !ok {
		panic("not function:" + name)
	}
	prog, err := getCompiledScript(name, e)
	if err != nil {
		panic(err)
	}

	vm := NewVm()

	v, err := vm.RunProgram(prog)
	if err != nil {
		panic(err)
	}

	return v.Export()
}
