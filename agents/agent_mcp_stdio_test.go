package agents

import (
	"context"
	"fmt"
	"testing"

	"github.com/budgies-nest/budgie/enums/base"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

// go test -v -run TestMCPStdioOneToolCall
func TestMCPStdioOneToolCall(t *testing.T) {

	bob, err := NewAgent("Bob",
		WithDMR(base.DockerModelRunnerContainerURL),
		WithParams(
			openai.ChatCompletionNewParams{
				Model: "k33g/qwen2.5:0.5b-instruct-q8_0",
				//Model:       "k33g/llama-xlam-2:8b-fc-r-q2_k",
				Temperature: openai.Opt(0.0), // IMPORTANT: set temperature to 0.0 to ensure the agent uses the tool
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(`
						Say Hello to Bob	
					`),
				},
				ParallelToolCalls: openai.Bool(false),
			},
		),
		WithMCPStdioClient(
			context.Background(),
			"go",
			STDIOCommandOptions{
				"run",
				"../laboratory/mcp-stdio-server/main.go",
			},
			EnvVars{},
		),
		WithMCPStdioTools(context.Background(), []string{"say_hello"}),
	)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}
	// Generate the tools detection completion
	detectedToolCalls, err := bob.ToolsCompletion(context.Background())
	if err != nil {
		t.Fatalf("Failed to get tools completion: %v", err)
	}
	fmt.Println("Number of Tool Calls:\n", len(detectedToolCalls))

	detectedToolCallsStr, err := helpers.ToolCallsToJSONString(detectedToolCalls)
	if err != nil {
		t.Fatalf("Error converting tool calls to JSON string: %v", err)
	}
	fmt.Println("Detected Tool Calls:\n", detectedToolCallsStr)

	results, err := bob.ExecuteMCPStdioToolCalls(context.Background(), detectedToolCalls)
	if err != nil {
		t.Fatalf("Failed to execute tool calls: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, but got %d", len(results))
	}
	fmt.Println("Results:\n", results)
	if results[0] != "Hello Bob" {
		t.Errorf("Expected result 'Hello Bob!', but got '%s'", results[0])
	}

}
