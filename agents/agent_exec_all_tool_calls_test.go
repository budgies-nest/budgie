package agents

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/budgies-nest/budgie/enums/base"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

// go test -v -run TestExecAllToolCalls
func TestExecAllToolCalls(t *testing.T) {

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
				Model:       "k33g/llama-xlam-2:8b-fc-r-q2_k", // NOTE: this model is able to detect several tool calls in a single request
				Temperature: openai.Opt(0.0),
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(`
						Add 10 and 32
						Add 12 and 30
						Add 40 and 2
						Add 5 and 37						
					`),
				},
				// Enable parallel tool calls to detect all tool calls in a single request
				// IMPORTANT: not all the models are able to detect several tool calls in a single request
				ParallelToolCalls: openai.Bool(true),
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

	if len(detectedToolCalls) != 4 {
		t.Errorf("😡 Expected 1 tool call, but got %d", len(detectedToolCalls))
	}

	detectedToolCallsStr, err := helpers.ToolCallsToJSONString(detectedToolCalls)
	if err != nil {
		t.Fatalf("😡 Failed to convert tool calls to JSON string: %v", err)
	}
	fmt.Println("Detected Tool Calls:\n", detectedToolCallsStr)

	results, err := bob.ExecuteToolCalls(detectedToolCalls,
		map[string]func(any) (any, error){

			"add": func(args any) (any, error) {
				a := args.(map[string]any)["a"].(float64)
				b := args.(map[string]any)["b"].(float64)
				return a + b, nil
			},
		},
	)
	if err != nil {
		t.Fatalf("😡 Failed to execute tool calls: %v", err)
	}
	if len(results) != 4 {
		t.Errorf("😡 Expected 4 result, but got %d", len(results))
	}

	if strings.Join(results, "") != "42424242" {
		t.Errorf("😡 Expected result to be 42, but got %v", results[0])
	}

}
