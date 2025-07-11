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
	ctx := context.Background()

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(base.DockerModelRunnerLocalURL),
		agents.WithParams(
			openai.ChatCompletionNewParams{
				Model:       "ai/qwen2.5:latest",
				Temperature: openai.Opt(0.0),
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(`
						search for the latest news about Docker and Kubernetes.	
					`),
				},
				//ParallelToolCalls: openai.Bool(false),
			},
		),
		agents.WithMCPStdioClient(
			ctx,
			"docker",
			agents.STDIOCommandOptions{
				"mcp",
				"gateway",
				"run",
			},
			agents.EnvVars{},
		),
		agents.WithMCPStdioTools(ctx, []string{"fetch_content", "search"}),
	)

	if err != nil {
		panic(err)
	}
	fmt.Println("🤖 Bob is ready to assist!")
	for i, tool := range bob.Params.Tools {
		fmt.Printf("🛠️ Tool %d: %s\n", i+1, tool.Function.Name)
		fmt.Printf("  Description: %s\n", tool.Function.Description)
		fmt.Printf("  Parameters: %s\n", tool.Function.Parameters)
	}

	// Generate the tools detection completion
	detectedToolCalls, err := bob.ToolsCompletion(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("⚡️ Number of Tool Calls:\n", len(detectedToolCalls))

	detectedToolCallsStr, err := helpers.ToolCallsToJSONString(detectedToolCalls)
	if err != nil {
		fmt.Println("Error converting tool calls to JSON string:", err)
		return
	}
	fmt.Println("🛠️ Detected Tool Calls:\n", detectedToolCallsStr)

	results, err := bob.ExecuteMCPStdioToolCalls(ctx, detectedToolCalls)
	if err != nil {
		fmt.Println("Error executing tool calls:", err)
		return
	}
	fmt.Println("📝 Results of Tool Calls:\n", results)

}
