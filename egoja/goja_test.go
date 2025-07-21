package egoja

import (
	"context"
	"fmt"
	"goframe-ex/egoja/gojaapi"
	"sync"
	"testing"
)

// 写函数方法
const funcs1 = `
 MyLib = {
	version: "1.0.0",
	fib:  function(n) {
			if(n==0||n==1) {
				return n
			}
		return this.fib(n-2)+this.fib(n-1)
	},
	fmtTest:  function(n) {
		//var fmt = imports("fmt")
		//fmt.Println("=============n======================")
		return n
	},

}

`

// BenchmarkGojaCache-12    	    2169	    652156 ns/op
func BenchmarkGojaCache(b *testing.B) {
	gojaapi.RegisterFunc("bdutil", funcs1)
	gojaapi.RegisterCommonParameter("ccc", "bbbb")
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			simple, err := ExecScriptById(context.TODO(), "AAA", `import "@/bdutil"
				return bdutil.fib(aaa??10) `, map[string]any{"aaa": 15, "cccc": "aaaa"})

			fmt.Println(simple, err)
		}()
	}
	wg.Wait()
}

func BenchmarkArgs(b *testing.B) {
	gojaapi.RegisterCommonParameter("TArgs", tArgs)
	for i := 0; i < b.N; i++ {
		_, err := ExecScriptById(context.TODO(), "AAA", `var aaa= ["a","b","ccc",ddd]
TArgs("a",  ...aaa)
`, map[string]any{
			"ddd": "123546",
		})
		fmt.Println(err)
	}
}

func BenchmarkArgsThread(b *testing.B) {
	gojaapi.RegisterCommonParameter("TArgs", tArgs)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := ExecScriptById(context.TODO(), "AAA", `var aaa= ["a","b","ccc",ddd]
TArgs("a",  ...aaa)
`, map[string]any{
				"ddd": "123546",
			})
			fmt.Println(err)
		}()

	}
	wg.Wait()
}
func tArgs(b string, args ...any) {
	fmt.Println(b)
	fmt.Println(args...)
}
