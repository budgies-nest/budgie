package main

import (
	"fmt"

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
			},
		}),
		agents.WithHTTPServer(agents.HTTPServerConfig{
			Port: "8080",
		}),
	)
	if err != nil {
		panic(err)
	}
	// Start the HTTP server
	fmt.Println("Starting HTTP server on port 8080...")
	err = bob.StartHttpServer()
	if err != nil {
		panic(err)
	}
}
