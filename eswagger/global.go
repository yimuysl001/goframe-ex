/**

knife4j 风格业务处理实现

*/

package eswagger

import (
	"context"
	_ "embed"
	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"path"
	"strings"
)

var (
	//go:embed doc/data.bin
	databin []byte

	CryptoKey = []byte("x76cgqt36i9c863bzmotuf8626dxiwu0")
)

const apijson = "/api.json"

type SwaggerInfo struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Version        string `json:"version"`
	TermsOfService string `json:"termsOfService"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	Email          string `json:"email"`
}

func InitWeb(ctx context.Context) bool {
	serverRoot := g.Cfg().MustGet(ctx, "server.serverRoot").String()
	if serverRoot == "" { // 未配置路径，不需要处理swagger
		return false
	}

	binContent, err := gaes.Decrypt(databin, CryptoKey)
	if err != nil {
		g.Log().Error(ctx, err)
		return false
	}

	if !strings.Contains(serverRoot, ":") { // 全路径不处理
		var pwd = gfile.Pwd()
		pwd = strings.ReplaceAll(pwd, "\\", "/")
		serverRoot = path.Join(pwd, serverRoot)
	}
	if !strings.HasSuffix(serverRoot, "/") {
		serverRoot = serverRoot + "/"
	}

	if err := gres.Add(string(binContent), serverRoot); err != nil {
		g.Log().Error(ctx, err)
		return false
	}

	g.Log().Debug(ctx, "swagger 添加完成")
	return true
}

func EnhanceOpenAPIDoc(ctx context.Context, s *ghttp.Server) {
	openapi := s.GetOpenApi()
	//openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	//openapi.Config.CommonResponseDataField = `Data`
	var swaggerInfo = new(SwaggerInfo)

	err := g.Cfg().MustGet(ctx, "plugin.swagger."+s.GetName()).Scan(swaggerInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	// API description.
	openapi.Info = goai.Info{
		Title:          swaggerInfo.Title,
		Description:    swaggerInfo.Description,
		Version:        swaggerInfo.Version,
		TermsOfService: swaggerInfo.TermsOfService,
		Contact: &goai.Contact{
			Name:  swaggerInfo.Name,
			URL:   swaggerInfo.Url,
			Email: swaggerInfo.Email,
		},
	}
}

func InitSwagger(ctx context.Context, f func(ctx context.Context) *goai.OpenApiV3, name ...any) {
	//swaggerPath := g.Cfg().MustGet(ctx, "server.swaggerPath").String()
	//if swaggerPath == "" {
	//	return
	//}
	openapiPath := g.Cfg().MustGet(ctx, "server.openapiPath").String()
	if openapiPath == "" {
		return
	}
	if !strings.HasPrefix(openapiPath, "/") {
		openapiPath = "/" + openapiPath
	}

	if !InitWeb(ctx) {
		return
	}
	s := g.Server(name...)
	EnhanceOpenAPIDoc(ctx, s)

	if f != nil {
		s.BindHandler(apijson, func(r *ghttp.Request) {
			r.Response.Header().Set("connection", "keep-alive")
			r.Response.Header().Set("content-type", "application/json;charset=UTF-8")

			//r.Response.RedirectTo(g.Cfg().MustGet(r.Context(), "server.openapiPath").String())
			r.Response.Write(f(r.Context()))
		})
	}

	if openapiPath != apijson {
		s.BindHandler(apijson, func(r *ghttp.Request) {
			r.Response.RedirectTo(openapiPath)
		})
	}

}
