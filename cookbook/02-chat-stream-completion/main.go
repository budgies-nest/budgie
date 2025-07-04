package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/enums/base"
	"github.com/budgies-nest/budgie/enums/environments"

	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

/*

curl http://localhost:12434/engines/llama.cpp/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{
	"model": "ai/qwen2.5:latest",
	"messages": [
		{
			"role": "system",
			"content": "You are a helpful assistant."
		},
		{
			"role": "user",
			"content": "Please write 500 words about the fall of Rome."
		}
	],
	"stream": true
}'

curl http://localhost:12434/engines/llama.cpp/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{
	"model": "unsloth/qwen3-gguf:4B-UD-Q4_K_XL",
	"messages": [
		{
			"role": "system",
			"content": "You are a helpful assistant."
		},
		{
			"role": "user",
			"content": "Please write 500 words about the fall of Rome."
		}
	],
	"stream": true
}'


*/

func main() {

	fmt.Println("🐳", helpers.DetectContainerEnvironment())
	modelRunnerBaseUrl := ""
	if helpers.DetectContainerEnvironment() == environments.Local {
		modelRunnerBaseUrl = base.DockerModelRunnerLocalURL
	} else {
		modelRunnerBaseUrl = base.DockerModelRunnerContainerURL
	}

	// Enable global logging at Info level
	agents.EnableLogging(agents.LogLevelInfo)

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(openai.ChatCompletionNewParams{
			//Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
			Model:       "ai/qwen2.5:latest",
			//Model: "unsloth/qwen3-gguf:4B-UD-Q4_K_XL",
			Temperature: openai.Opt(0.8),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You're a helpful assistant expert with Star Trek universe."),
				openai.UserMessage("Who is James T Kirk?"),
			},
		}),
		agents.WithLoggingEnabled(),
		agents.WithLogLevel(agents.LogLevelError),
	)
	if err != nil {
		panic(err)
	}
	_, err = bob.ChatCompletionStream(context.Background(), func(self *agents.Agent, content string, err error) error {
		fmt.Print(content)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
