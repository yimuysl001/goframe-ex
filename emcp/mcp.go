package emcp

import "github.com/mark3labs/mcp-go/server"

/**
mcp:
    name: Name_MCP
    version: v1.0.0
    sse_path: /sse
    message_path: /message
    url_prefix: ''
*/

type McpConfig struct {
	Name        string
	Version     string
	SSEPath     string
	MessagePath string
	UrlPrefix   string
}

func DefaultConfig() McpConfig {
	return McpConfig{
		Name:        "emcp",
		Version:     "1.0",
		SSEPath:     "/sse",
		MessagePath: "/message",
		UrlPrefix:   "/",
	}
}

func McpRun(config McpConfig) *server.SSEServer {
	s := server.NewMCPServer(
		config.Name,
		config.Version,
	)
	RegisterAllTools(s)

	return server.NewSSEServer(s,
		server.WithSSEEndpoint(config.SSEPath),
		server.WithMessageEndpoint(config.MessagePath),
		server.WithBaseURL(config.UrlPrefix))
}
