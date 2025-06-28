package agents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
)

// ExecuteToolCalls executes the tool calls detected by the Agent.
// QUESTION: Should I return []any instead of []string?
func (agent *Agent) ExecuteToolCalls(detectedtToolCalls []openai.ChatCompletionMessageToolCall, toolsImpl map[string]func(any) (any, error)) ([]string, error) {
	responses := []string{}
	for _, toolCall := range detectedtToolCalls {
		// Check if the tool is implemented
		toolFunc, ok := toolsImpl[toolCall.Function.Name]

		if !ok { // NOTE: the tool is not implemented
			//return nil, fmt.Errorf("tool %s not implemented", toolCall.Function.Name)
			//fmt.Printf("✋ tool %s not implemented", toolCall.Function.Name)
			continue
		}

		var args map[string]any
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			return nil, err
		}

		// Call the tool with the arguments
		toolResponse, err := toolFunc(args)
		if err != nil {
			responses = append(responses, fmt.Sprintf("%v", err))
		} else {
			responses = append(responses, fmt.Sprintf("%v", toolResponse))
			agent.Params.Messages = append(
				agent.Params.Messages,
				openai.ToolMessage(
					fmt.Sprintf("%v", toolResponse),
					toolCall.ID,
				),
			)
		}
	}
	if len(responses) == 0 {
		return nil, errors.New("no tool responses found")
	}
	return responses, nil
}


// TODO: check what will happend if the tool does not xist
// ExecuteMCPStdioToolCalls executes the tool calls detected by the Agent using the MCP STDIO client.
func (agent *Agent) ExecuteMCPStdioToolCalls(ctx context.Context, detectedtToolCalls []openai.ChatCompletionMessageToolCall) ([]string, error) {
	responses := []string{}
	for _, toolCall := range detectedtToolCalls {

		var args map[string]any
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			return nil, err
		}

		// NOTE: Call the MCP tool with the arguments
		request := mcp.CallToolRequest{}
		request.Params.Name = toolCall.Function.Name
		request.Params.Arguments = args

		// Call the tool with the arguments thanks to the MCP client
		toolResponse, err := agent.mpcStdioClient.CallTool(ctx, request)
		if err != nil {
			responses = append(responses, fmt.Sprintf("%v", err))
		} else {
			if toolResponse != nil && len(toolResponse.Content) > 0 {
				// TODO: test if the content is a TextContent 
				result := toolResponse.Content[0].(mcp.TextContent).Text

				agent.Params.Messages = append(
					agent.Params.Messages,
					openai.ToolMessage(
						result,
						toolCall.ID,
					),
				)
				responses = append(responses, result)
			}
		}

	}
	if len(responses) == 0 {
		return nil, errors.New("no tool responses found")
	}
	return responses, nil
}

// ExecuteMCPStreamableHTTPToolCalls executes the tool calls detected by the Agent using the MCP Streamable HTTP client.
func (agent *Agent) ExecuteMCPStreamableHTTPToolCalls(ctx context.Context, detectedtToolCalls []openai.ChatCompletionMessageToolCall) ([]string, error) {
	responses := []string{}
	for _, toolCall := range detectedtToolCalls {

		var args map[string]any
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			return nil, err
		}

		// NOTE: Call the MCP tool with the arguments
		request := mcp.CallToolRequest{}
		request.Params.Name = toolCall.Function.Name
		request.Params.Arguments = args

		// Call the tool with the arguments thanks to the MCP client
		toolResponse, err := agent.mcpStreamableHTTPClient.CallTool(ctx, request)
		if err != nil {
			responses = append(responses, fmt.Sprintf("%v", err))
		} else {
			if toolResponse != nil && len(toolResponse.Content) > 0 {
				// TODO: test if the content is a TextContent 
				result := toolResponse.Content[0].(mcp.TextContent).Text

				agent.Params.Messages = append(
					agent.Params.Messages,
					openai.ToolMessage(
						result,
						toolCall.ID,
					),
				)
				responses = append(responses, result)
			}
		}

	}
	if len(responses) == 0 {
		return nil, errors.New("no tool responses found")
	}
	return responses, nil
}
