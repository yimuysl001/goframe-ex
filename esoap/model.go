package esoap

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	mysoap "goframe-ex/esoap/soap"
)

type Model struct {
	Names   string                                                             //方法名
	DoThing func(ctx context.Context, model Model, str string) (string, error) //数据处理具体方法
}

type ModelRequest struct {
	Xml string
}

type ModelResponse struct {
	ResultInfo string
}

// todo 方法名
func (u Model) Name() string {
	return u.Names
}

// todo  接收数据
func (u Model) ReqStruct() interface{} {
	return ModelRequest{}
}

// todo 返回数据
func (u Model) RespStruct() interface{} {
	return ModelResponse{}
}

// todo  model 数据处理
func (u Model) Do(ctx context.Context, req interface{}, resp interface{}) error {
	var err error
	re := req.(*ModelRequest)
	res := resp.(*ModelResponse)
	g.Log().Info(ctx, "=================", u.Name(), "开始=========================")
	res.ResultInfo, err = u.DoThing(ctx, u, re.Xml)
	if err != nil {
		g.Log().Error(ctx, err.Error())
		res.ResultInfo = err.Error()
	}
	//mylog.Trace(res.ResultInfo)
	g.Log().Info(ctx, "=================", u.Name(), "结束=========================")
	return nil
}

type WebInter interface {
	WebService(prefix string, Name string) func(group *ghttp.RouterGroup)
}

type Webmodel struct {
	Models    []Model //数据处理实例
	NameSpace string  // http://ip:端口
}

// s.Group("/webservice", webmodel.WebService("/webservice", consts.YxConfig.MethodName))
func (w *Webmodel) WebService(prefix string, Name string) func(group *ghttp.RouterGroup) {
	//if w.registerfuc==nil{
	//	panic("未设置注册方式")
	//}
	return func(group *ghttp.RouterGroup) {
		//group.Middleware(mid.AuthMiddleware)
		my := mysoap.NewServer(Name, w.NameSpace+prefix+"/"+Name)
		w.registerIServerMethod(my)
		group.GET("/"+Name, my.Handler)
		group.POST("/"+Name, my.Handler)
	}
}
func (w *Webmodel) registerIServerMethod(s mysoap.IServer) {
	if len(w.Models) == 0 {
		panic("未设置方法实例")
	}
	for _, model := range w.Models {
		if model.DoThing == nil {
			panic(model.Names + "未设置处理方法")
		}
		s.RegisterMethod(model)
	}
}
