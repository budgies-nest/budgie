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

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(base.DockerModelRunnerContainerURL),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
			Temperature: openai.Opt(0.8),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You're a helpful assistant expert with Star Trek universe."),
				openai.UserMessage("Who is James T Kirk?"),
			},
		}),
		agents.WithLoggingEnabled(),
		agents.WithLogLevel(agents.LogLevelDebug),
	)
	if err != nil {
		panic(err)
	}
	response, err := bob.ChatCompletion(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Response from Bob:", response)

	// Disable logging
	fmt.Println("\n=== Disabling Logging ===")
	agents.DisableLogging()
}
