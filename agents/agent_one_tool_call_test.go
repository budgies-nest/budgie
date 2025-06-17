package agents

import (
	"context"
	"fmt"
	"testing"

	"github.com/budgies-nest/budgie/enums/base"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

// go test -v -run TestOneToolCall
func TestOneToolCall(t *testing.T) {

	addTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "add",
			Description: openai.String("add two numbers"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]string{
						"type":        "number",
						"description": "The first number to add.",
					},
					"b": map[string]string{
						"type":        "number",
						"description": "The second number to add.",
					},
				},
				"required": []string{"a", "b"},
			},
		},
	}

	bob, err := NewAgent("Bob",
		WithDMR(context.Background(), base.DockerModelRunnerContainerURL),
		WithParams(
			openai.ChatCompletionNewParams{
				Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
				Temperature: openai.Opt(0.0),
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(`
						Add 10 and 32			
					`),
				},
				ParallelToolCalls: openai.Bool(false),
			},
		),
		WithTools([]openai.ChatCompletionToolParam{addTool}),
	)
	if err != nil {
		t.Fatalf("😡 Failed to create agent: %v", err)
	}

	// Generate the tools detection completion
	detectedToolCalls, err := bob.ToolsCompletion()
	if err != nil {
		t.Fatalf("😡 Failed to get tools completion: %v", err)
	}
	fmt.Println("Number of Tool Calls:\n", len(detectedToolCalls))

	if len(detectedToolCalls) != 1 {
		t.Errorf("😡 Expected 1 tool call, but got %d", len(detectedToolCalls))
	}

	detectedToolCallsStr, err := helpers.ToolCallsToJSONString(detectedToolCalls)
	if err != nil {
		t.Fatalf("😡 Failed to convert tool calls to JSON string: %v", err)
	}
	fmt.Println("Detected Tool Calls:\n", detectedToolCallsStr)

}
