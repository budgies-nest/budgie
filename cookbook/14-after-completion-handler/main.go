package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/enums/base"
	"github.com/openai/openai-go"
)

func main() {

	// Enable global logging at Info level
	agents.EnableLogging(agents.LogLevelInfo)

	alice, err := agents.NewAgent("Alice",
		agents.WithDMR(base.DockerModelRunnerContainerURL),
		agents.WithParams(openai.ChatCompletionNewParams{
			//Model:       "ai/gemma3n:2B-Q4_K_M",
			//Model: "ai/qwen2.5:0.5B-F16",
			Model: "k33g/qwen2.5:0.5b-instruct-q2_k",

			Temperature: openai.Opt(0.1),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a helpful assistant."),
				openai.UserMessage("What is the best Pizza int the world?"),
			},
		}),
		// Add a handler to modify the response
		agents.WithAfterChatCompletion(func(ctx *agents.ChatCompletionContext) {
			if ctx.Response != nil && ctx.Error == nil {
				// Add a friendly prefix to the response
				original := *ctx.Response
				*ctx.Response = "🍕🍕🍕 " + original
			}
		}),
	)
	if err != nil {
		panic(err)
	}

	response, err := alice.ChatCompletion(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Modified Response: %s\n", response)

}
