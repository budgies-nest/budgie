package agents

import (
	"net/http"

	"github.com/budgies-nest/budgie/rag"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/server"
	"github.com/openai/openai-go"
)

type Agent struct {
	clientEngine    openai.Client
	Name            string
	Params          openai.ChatCompletionNewParams
	EmbeddingParams openai.EmbeddingNewParams
	
	//Store           rag.MemoryVectorStore
	Store rag.VectorStore
	storeFilePath string

	mcpServerConfig MCPServerConfig
	mcpServer       *server.MCPServer

	httpServerConfig HTTPServerConfig
	httpServer       *http.ServeMux

	//ToolCalls []openai.ChatCompletionMessageToolCall
	//Instructions openai.ChatCompletionMessageParamUnion

	Metadata map[string]any

	optionError error

	// MCP Clients
	mpcStdioClient          *client.Client
	mcpStreamableHTTPClient *client.Client
	// TODO: QUESTION: how to filter the tools?
	// NOTE: by filteryng the tools catalog

	// Logger
	logger *Logger
}

type AgentOption func(*Agent)

// NewAgent creates a new Agent instance with the provided options.
// It applies all the options to the Agent and returns it.
// If any option sets an error, it returns the error instead of the Agent.
// The Agent can be configured with various options such as DMR client, parameters, tools, and memory.
func NewAgent(name string, options ...AgentOption) (*Agent, error) {

	agent := &Agent{}
	agent.Name = name
	agent.logger = GetGlobalLogger()
	// Apply all options
	for _, option := range options {
		option(agent)
	}
	if agent.optionError != nil {
		return nil, agent.optionError
	}
	return agent, nil
}

