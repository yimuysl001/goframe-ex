package emcp

import "github.com/gogf/gf/v2/net/ghttp"

func RegisterConfig(s *ghttp.Server, config McpConfig) {
	sseServer := McpRun(config)

	s.Group(config.UrlPrefix, func(group *ghttp.RouterGroup) {
		group.GET(config.SSEPath, func(r *ghttp.Request) {
			sseServer.SSEHandler().ServeHTTP(r.Response.Writer, r.Request)
		})

		group.POST(config.MessagePath, func(r *ghttp.Request) {
			sseServer.MessageHandler().ServeHTTP(r.Response.Writer, r.Request)
		})

	})

}
