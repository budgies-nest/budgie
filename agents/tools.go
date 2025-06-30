package agents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
		start := time.Now()
		toolResponse, err := toolFunc(args)
		duration := time.Since(start)

		responseStr := fmt.Sprintf("%v", toolResponse)
		if err != nil {
			responseStr = fmt.Sprintf("%v", err)
			responses = append(responses, responseStr)
			agent.logger.LogToolExecution(agent.Name, toolCall.Function.Name, args, responseStr, duration, err)
		} else {
			responses = append(responses, responseStr)
			agent.Params.Messages = append(
				agent.Params.Messages,
				openai.ToolMessage(
					responseStr,
					toolCall.ID,
				),
			)
			agent.logger.LogToolExecution(agent.Name, toolCall.Function.Name, args, responseStr, duration, nil)
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
		start := time.Now()
		toolResponse, err := agent.mpcStdioClient.CallTool(ctx, request)
		duration := time.Since(start)

		if err != nil {
			responseStr := fmt.Sprintf("%v", err)
			responses = append(responses, responseStr)
			agent.logger.LogMCPToolExecution(agent.Name, toolCall.Function.Name, args, responseStr, "stdio", duration, err)
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
				agent.logger.LogMCPToolExecution(agent.Name, toolCall.Function.Name, args, result, "stdio", duration, nil)
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
		start := time.Now()
		toolResponse, err := agent.mcpStreamableHTTPClient.CallTool(ctx, request)
		duration := time.Since(start)

		if err != nil {
			responseStr := fmt.Sprintf("%v", err)
			responses = append(responses, responseStr)
			agent.logger.LogMCPToolExecution(agent.Name, toolCall.Function.Name, args, responseStr, "http", duration, err)
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
				agent.logger.LogMCPToolExecution(agent.Name, toolCall.Function.Name, args, result, "http", duration, nil)
			}
		}

	}
	if len(responses) == 0 {
		return nil, errors.New("no tool responses found")
	}
	return responses, nil
}
