package egoja

import (
	"github.com/dop251/goja"
	"goframe-ex/egoja/gojaapi"
	_ "goframe-ex/egoja/pkgs"
	"strings"
)

const mainFunc = `
function Main() {
	%s
}
Main()
`

// ExecSimple 简易数据处理
func ExecSimple(script string, params map[string]any) (goja.Value, error) {
	vm := gojaapi.GetNewVm()
	defer gojaapi.PutGoja(vm)

	for k, v := range params {
		vm.Set(k, v)
	}

	return vm.RunString(script)
}
func JsonConversion(refresh, s string) (string, string) {
	s = strings.TrimSpace(s)
	if (strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) ||
		(strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) {
		s = "return  JSON.stringify(" + s + ")"
	}
	return refresh, s
}
