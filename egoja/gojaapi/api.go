package gojaapi

import (
	"errors"
	"github.com/dop251/goja"
	"path"
	"strings"
)

func GetNewVm() *goja.Runtime {
	//vm := goja.New()
	//vm.Set("imports", Import)
	//vm.Set("importFunc", ImportFunc)
	//vm.Set("db", dbutil.DB)
	//vm := vmPool.Get().(*goja.Runtime)
	// 将 VM 放回池中以供将来重用
	//defer vmPool.Put(vm)
	return vmPool.Get().(*goja.Runtime) //  vm

	//return NewVm()

}
func PutGoja(vm *goja.Runtime) {
	vmPool.Put(vm)
}

func ScriptRefresh(script string) (string, string, error) {
	var (
		imporsb strings.Builder
		mainsb  strings.Builder
	)
	for _, s := range strings.Split(script, "\n") {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "import ") {
			s = strings.Replace(s, "import ", "", 1)
			n := strings.SplitN(s, "\"", 2)
			if len(n) != 2 {
				return "", "", errors.New("import 设置有错误:" + s)
			}
			if strings.TrimSpace(n[0]) == "" {
				base := path.Base(n[1])
				base = strings.Replace(base, "\"", "", -1)
				//if strings.Contains(base, ".") {
				//	imporsb.WriteString("var " + strings.SplitN(base, ".", 2)[1])
				//} else {
				imporsb.WriteString("var " + base)
				//}

			} else {
				imporsb.WriteString("var " + n[0])
			}

			if strings.HasPrefix(n[1], "@/") {
				replace := strings.Replace(n[1], "@/", "", 1)
				imporsb.WriteString("= importFunc(\"" + replace + ")\n")
			} else {
				imporsb.WriteString("= imports(\"" + n[1] + ")\n")
			}

		} else {
			mainsb.WriteString(s)
			mainsb.WriteString("\n")
		}

	}
	//imporsb.WriteString(mainsb.String())

	return imporsb.String(), mainsb.String(), nil

}

func ScriptOutRefresh(script string) (string, string, error) {
	var (
		imporsb strings.Builder
		mainsb  strings.Builder
	)
	for _, s := range strings.Split(script, "\n") {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "import ") {
			s = strings.Replace(s, "import ", "", 1)
			n := strings.SplitN(s, "\"", 2)
			if len(n) != 2 {
				return "", "", errors.New("import 设置有错误:" + s)
			}
			if strings.TrimSpace(n[0]) == "" {
				base := path.Base(n[1])
				base = strings.Replace(base, "\"", "", -1)
				//if strings.Contains(base, ".") {
				//	imporsb.WriteString("var " + strings.SplitN(base, ".", 2)[1])
				//} else {
				imporsb.WriteString("var " + base)
				//}

			} else {
				imporsb.WriteString("var " + n[0])
			}

			if strings.HasPrefix(n[1], "@/") {
				replace := strings.Replace(n[1], "@/", "", 1)
				imporsb.WriteString("= importFunc(\"" + replace + ")\n")
			} else {
				imporsb.WriteString("= imports(\"" + n[1] + ")\n")
			}

		} else {
			if strings.HasPrefix(s, "return ") {
				replace := strings.TrimSpace(strings.Replace(s, "return ", "", 1))
				if replace != "" && replace != ";" {
					s = strings.Replace(s, "return ", "return  _out.out=", 1)
				}
			}

			mainsb.WriteString(s)
			mainsb.WriteString("\n")
		}

	}

	return imporsb.String(), mainsb.String(), nil

}
