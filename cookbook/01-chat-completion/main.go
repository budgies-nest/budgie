package main

import (
	"context"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/enums/base"
	"github.com/openai/openai-go"
)

func main() {

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
	)
	if err != nil {
		panic(err)
	}
	response, err := bob.ChatCompletion(context.Background())
	if err != nil {
		panic(err)
	}
	println("Response from Bob:", response)
}
