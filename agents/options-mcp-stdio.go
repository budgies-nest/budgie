package agents

import (
	"context"
	"errors"

	"slices"

	"github.com/budgies-nest/budgie/enums/constants"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
)

type STDIOCommandOptions []string
type EnvVars []string

// WithMCPSTDIOClient initializes the Agent with an MCP client using the provided command.
// It runs the command to connect to the MCP server and sets up the client transport.
// The command should be a valid command that can be executed in the environment where the agent runs.
// It returns an AgentOption that can be used to configure the agent.
func WithMCPStdioClient(ctx context.Context, cmd string, options STDIOCommandOptions, envvars EnvVars) AgentOption {
	return func(agent *Agent) {
		//agent.ctx = ctx

		mcpClient, err := client.NewStdioMCPClient(
			cmd,
			envvars, // Environment variables for the MCP client
			options...,
		)
		if err != nil {
			agent.optionError = err // TODO: check if the error is used in the Agent constructor
			return
		}
		// QUESTION: Should I defer the client close here?
		//defer mcpClient.Close()

		// Initialize the client
		// 	fmt.Println("Initializing client...") // TODO: make a logger for the agent
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    constants.MCPStdioClientName,
			Version: constants.MCPStdioClientVersion,
		}

		//initResult, err := mcpClient.Initialize(ctx, initRequest)
		_, err = mcpClient.Initialize(ctx, initRequest)

		if err != nil {
			// Failed to initialize
			agent.optionError = err
			return
		}
		agent.mpcStdioClient = mcpClient
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

// WithMCPStdioTools fetches the tools from the MCP server and sets them in the agent.
// It filters the tools based on the provided names and converts them to OpenAI format.
// It requires the MCP server to be running and accessible at the specified address.
// The tools are expected to be in the format defined by the MCP server.
// It returns an AgentOption that can be used to configure the agent.
// The tools are fetched using the MCP client and converted to OpenAI format.
// If no toolsFilter are specified, all available tools are used.
// If toolsFilter is specified, only the tools matching the filter are used.
// IMPORTANT: The tools are appended to the existing tools in the Agent's parameters.
func WithMCPStdioTools(ctx context.Context, toolsFilter []string) AgentOption {
	return func(agent *Agent) {
		// Get the tools from the MCP client
		toolsRequest := mcp.ListToolsRequest{}
		mcpTools, err := agent.mpcStdioClient.ListTools(ctx, toolsRequest)
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
