package agents

import (
	"context"
	"errors"
	"slices"

	"github.com/budgies-nest/budgie/enums/constants"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
)

type StreamableHttpOptions []string

// NOTE: this is subject to change
func WithMCPStreamableHttpClient(ctx context.Context, mcpHttpServerUrl string, options StreamableHttpOptions) AgentOption {

	return func(agent *Agent) {

		httpTransport, err := transport.NewStreamableHTTP(mcpHttpServerUrl) // TODO: add the options
		if err != nil {
			agent.optionError = err // TODO: check if the error is used in the Agent constructor
			return
		}

		mcpClient := client.NewClient(httpTransport)
		if err := mcpClient.Start(ctx); err != nil {
			agent.optionError = err
			return
		}

		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    constants.MCPStreamableHTTPClientName,
			Version: constants.MCPStreamableHTTPClientVersion,
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		//initResult, err := mcpClient.Initialize(ctx, initRequest)
		_, err = mcpClient.Initialize(ctx, initRequest)
		if err != nil {
			// Failed to initialize
			agent.optionError = err
			return
		}

		agent.mcpStreamableHTTPClient = mcpClient
		// TODO: make a logger for the agent
		/*
			fmt.Printf(
				"Initialized with server: %s %s\n\n",
				initResult.ServerInfo.Name,
				initResult.ServerInfo.Version,
			)
		*/
	}

}

func WithMCPStreamableHttpTools(ctx context.Context, toolsFilter []string) AgentOption {
	return func(agent *Agent) {
		// Get the tools from the MCP client
		toolsRequest := mcp.ListToolsRequest{}
		mcpTools, err := agent.mcpStreamableHTTPClient.ListTools(ctx, toolsRequest)
		if err != nil {
			agent.optionError = err
			return
		}

		if len(toolsFilter) == 0 {
			// If no tools are specified, use all available tools
			// Convert the tools to OpenAI format
			convertedTools := helpers.ConvertMCPToolsToOpenAITools(mcpTools)
			agent.Params.Tools = append(agent.Params.Tools, convertedTools...)

		} else {
			filteredTools := []openai.ChatCompletionToolParam{}
			convertedTools := helpers.ConvertMCPToolsToOpenAITools(mcpTools)

			// filter convertedTools with toolsFilter and add them to filteredTools
			for _, tool := range convertedTools {
				if slices.Contains(toolsFilter, tool.Function.Name) {
					filteredTools = append(filteredTools, tool) // No need to check other filters for this tool
				}
			}
			if len(filteredTools) == 0 {
				agent.optionError = errors.New("no tools found matching the filter")
				return
			}
			// Append the filtered tools to the agent's parameters
			agent.Params.Tools = append(agent.Params.Tools, filteredTools...)

		}

	}
}
