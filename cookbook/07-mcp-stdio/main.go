package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/enums/base"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func main() {

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(base.DockerModelRunnerContainerURL),
		agents.WithParams(
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
		agents.WithMCPStdioClient(
			context.Background(),
			"go",
			agents.STDIOCommandOptions{
				"run",
				"../../laboratory/mcp-stdio-server/main.go",
			},
			agents.EnvVars{},
		),
		agents.WithMCPStdioTools(context.Background(), []string{"say_hello"}),
	)

	if err != nil {
		panic(err)
	}
	fmt.Println("🤖 Bob is ready to assist!", bob.Params.Tools)

	// Generate the tools detection completion
	detectedToolCalls, err := bob.ToolsCompletion(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Number of Tool Calls:\n", len(detectedToolCalls))

	detectedToolCallsStr, err := helpers.ToolCallsToJSONString(detectedToolCalls)
	if err != nil {
		fmt.Println("Error converting tool calls to JSON string:", err)
		return
	}
	fmt.Println("Detected Tool Calls:\n", detectedToolCallsStr)

	results, err := bob.ExecuteMCPStdioToolCalls(context.Background(), detectedToolCalls)
	if err != nil {
		fmt.Println("Error executing tool calls:", err)
		return
	}
	fmt.Println("Results of Tool Calls:\n", results)

}
