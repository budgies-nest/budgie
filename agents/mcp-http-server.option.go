package agents

import (
	"github.com/mark3labs/mcp-go/server"
)

type MCPServerConfig struct {
	Name     string
	Version  string
	Port     string
	Endpoint string
}

func WithMCPStreamableHttpServer(mcpServerConfig MCPServerConfig) AgentOption {
	return func(agent *Agent) {
		agent.mcpServerConfig = mcpServerConfig

		if mcpServerConfig.Endpoint == "" {
			// Default endpoint path for MCP server
			agent.mcpServerConfig.Endpoint = "/mcp"
		}
		if mcpServerConfig.Port == "" {
			// Default port for MCP server
			agent.mcpServerConfig.Port = "9090"
		}

		// Create MCP server
		agent.mcpServer = server.NewMCPServer(
			mcpServerConfig.Name,
			mcpServerConfig.Version,
		)

	}
}
