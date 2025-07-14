package agents

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (agent *Agent) AddToolToMCPServer(tool mcp.Tool, handler server.ToolHandlerFunc) {
	if agent.mcpServer == nil {
		// If the MCP server is not initialized, we cannot add tools
		return
	}
	// Add the tool to the MCP server
	agent.mcpServer.AddTool(tool, handler)
}

func (agent *Agent) StartMCPHttpServer() error {
	return server.NewStreamableHTTPServer(agent.mcpServer,
		server.WithEndpointPath(agent.mcpServerConfig.Endpoint),
	).Start(":" + agent.mcpServerConfig.Port)
}

func (agent *Agent) MCPServerConfig() MCPServerConfig {
	return agent.mcpServerConfig
}