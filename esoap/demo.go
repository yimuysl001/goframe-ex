package esoap

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

// GetWebModelTest nameSpace 命名空间 ip:端口
// prefix 服务前缀 /开头
// name 服务路径
func GetWebModelTest(s *ghttp.Server, nameSpace, prefix, name string, models ...Model) {

	var webmodel = Webmodel{
		NameSpace: "http://" + nameSpace,
	}

	if models != nil && models[0].Name() != "" {
		webmodel.Models = append(webmodel.Models, models...)
	} else {
		webmodel.Models = append(webmodel.Models, Model{
			Names: "Test",
			DoThing: func(ctx context.Context, model Model, str string) (string, error) {
				return "测试成功：" + str, nil
			},
		})
	}
	s.Group(prefix, webmodel.WebService(prefix, name))

	g.Log().Info(gctx.GetInitCtx(), "webservice 请求地址：", "http://"+nameSpace+"/"+prefix+"/"+name+"?wsdl")

}
