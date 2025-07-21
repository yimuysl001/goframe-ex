package econf

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gurkankaymak/hocon"
	"goframe-ex/epongo"
	"strings"
)

var (
	localHoconPath = []string{"config.conf", "config/config.conf", "manifest/config/config.conf"}
	exPath         = "|"
)

type AdapterContent struct {
	jsonVar *gvar.Var // The pared JSON object for configuration content, type: *gjson.Json.
}

func (a *AdapterContent) SetContent(content string) error {
	j, err := gjson.LoadContent([]byte(content), true)
	if err != nil {
		return gerror.Wrap(err, `load configuration content failed`)
	}
	a.jsonVar.Set(j)
	return nil
}

func (a *AdapterContent) Available(ctx context.Context, resource ...string) (ok bool) {
	if a.jsonVar.IsNil() {
		return false
	}
	return true
}

func (a *AdapterContent) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if a.jsonVar.IsNil() {
		return nil, nil
	}
	return a.jsonVar.Val().(*gjson.Json).Get(pattern).Val(), nil
}

func (a *AdapterContent) Data(ctx context.Context) (data map[string]interface{}, err error) {
	if a.jsonVar.IsNil() {
		return nil, nil
	}

	return a.jsonVar.Val().(*gjson.Json).Var().Map(), nil
}

func NewAdapterFile(filename ...string) *AdapterContent {
	var ctx = gctx.GetInitCtx()
	d, err := g.Cfg(filename...).Data(ctx)
	if err != nil {
		panic(err)
	}
	viewer := d["viewer"] // 有特殊处理
	datajson := gjson.New(d)
	err = datajson.Remove("viewer")
	if err != nil {
		panic(err)
	}
	yamlString, err := datajson.ToJsonIndentString()
	if err != nil {
		panic(err)
	}
	content, err := epongo.ParseContent(gctx.GetInitCtx(), yamlString, nil)
	if err != nil {
		panic(err)
	}
	json := gjson.New(content)
	err = json.Set("viewer", viewer)
	if err != nil {
		panic(err)
	}
	if len(filename) < 1 || filename[0] == "" || filename[0] == "config" {
		confMap := getConf()
		for k, v := range confMap {
			err = json.Set(k, v)
			if err != nil {
				g.Log().Error(ctx, "hocon导入数据出错:", err)
				continue
			}
		}
		g.Log().Info(ctx, "修改本地文件配置")
	}

	get := json.Get("gf.gcfg.file")

	// 避免无限调用
	if get.IsNil() || get.IsEmpty() || get.String() == "" || strings.Contains(exPath, "|"+get.String()+"|") {
		a := &AdapterContent{
			jsonVar: gvar.New(json, true),
		}
		return a
	}
	exPath = exPath + get.String() + "|"
	file := NewAdapterFile(get.String())

	data, err := file.Data(ctx)
	if err != nil {
		panic(err)
	}
	for k, v := range data {
		err = json.Set(k, v)
		if err != nil {
			g.Log().Error(ctx, get.String()+"导入失败:", err)
			continue
		}
	}

	a := &AdapterContent{
		jsonVar: gvar.New(json, true),
	}

	return a

}

func NewAdapterContent(contents ...string) *AdapterContent {

	var ctx = gctx.GetInitCtx()

	d, err := g.Cfg().Data(ctx)
	if err != nil {
		panic(err)
	}

	viewer := d["viewer"] // 有特殊处理

	datajson := gjson.New(d)

	err = datajson.Remove("viewer")
	if err != nil {
		panic(err)
	}
	yamlString, err := datajson.ToJsonIndentString()
	if err != nil {
		panic(err)
	}
	content, err := epongo.ParseContent(gctx.GetInitCtx(), yamlString, nil)
	if err != nil {
		panic(err)
	}
	json := gjson.New(content)
	err = json.Set("viewer", viewer)
	if err != nil {
		panic(err)
	}

	// 刷新数据参数
	for _, c := range contents {
		if c == "" {
			continue
		}
		nc, err2 := epongo.ParseContent(gctx.GetInitCtx(), c, nil)
		if err2 != nil {
			g.Log().Error(ctx, "导入数据出错：", err2)
			continue
		}
		cjson := gjson.New(nc)
		if cjson.IsNil() {
			g.Log().Error(ctx, "数据转换json失败:", nc)
			continue
		}
		for k, v := range cjson.Map() {
			err = json.Set(k, v)
			if err != nil {
				g.Log().Error(ctx, "导入数据出错:", err)
				continue
			}
		}

	}

	confMap := getConf()
	for k, v := range confMap {
		err = json.Set(k, v)
		if err != nil {
			g.Log().Error(ctx, "hocon导入数据出错:", err)
			continue
		}
	}

	//ctx := gctx.New()
	g.Log().Info(ctx, "修改本地文件配置")

	a := &AdapterContent{
		jsonVar: gvar.New(json, true),
	}

	return a
}

func getConf() map[string]any {
	var confMap = make(map[string]any)
	for _, path := range localHoconPath {
		if !gfile.Exists(path) {
			continue
		}
		contents := gfile.GetContents(path)
		content, err := epongo.ParseContent(gctx.GetInitCtx(), contents, nil)
		if err != nil {
			panic(err)
		}
		parseString, err := hocon.ParseString(content)
		if err != nil {
			panic(err)
		}

		confMap = gjson.New(parseString.GetRoot()).Map()
		break

	}

	return confMap

}
