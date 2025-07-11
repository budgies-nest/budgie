package main

import (
	"context"
	"fmt"
	"os"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func main() {

	modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")
	fmt.Println("Using Model Runner Base URL:", modelRunnerBaseUrl)

	if modelRunnerBaseUrl == "" {
		panic("MODEL_RUNNER_BASE_URL environment variable is not set")
	}

	modelRunnerChatModel := os.Getenv("MODEL_RUNNER_CHAT_MODEL")
	fmt.Println("Using Model Runner Chat Model:", modelRunnerChatModel)

	if modelRunnerChatModel == "" {
		panic("MODEL_RUNNER_CHAT_MODEL environment variable is not set")
	}

	ctx := context.Background()

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(
			openai.ChatCompletionNewParams{
				Model:       modelRunnerChatModel,
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
			"socat",
			agents.STDIOCommandOptions{
				"STDIO",
				"TCP:mcp-gateway:8811",
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
