package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
*/

func getModelRunnerBaseUrl() string {
	// Detect if running in a container or locally
	if helpers.DetectContainerEnvironment() == environments.Local {
		return base.DockerModelRunnerLocalURL
	}
	return base.DockerModelRunnerContainerURL
}

func main() {

	modelRunnerBaseUrl := getModelRunnerBaseUrl()

	// Enable global logging at Info level
	//agents.EnableLogging(agents.LogLevelInfo)

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(openai.ChatCompletionNewParams{
			//Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
			Model: "hf.co/unsloth/qwen3-0.6b-gguf:q4_k_m",
			//Model: "unsloth/qwen3-gguf:4B-UD-Q4_K_XL",
			Temperature: openai.Opt(0.8),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You're a helpful assistant expert with Star Trek universe."),
				openai.UserMessage("Who is James T Kirk?"),
			},
		}),
		agents.WithLoggingEnabled(),
		agents.WithLogLevel(agents.LogLevelError),
		agents.WithHTTPServer(agents.HTTPServerConfig{
			Port: "5050",
		}),
	)
	if err != nil {
		panic(err)
	}
	// TODO: check how the context is handled with the REST API
	// _, err = bob.ChatCompletionStream(context.Background(), func(self *agents.Agent, content string, err error) error {
	// 	fmt.Print(content)
	// 	return nil
	// })

	bob.HttpServer().HandleFunc("POST /api/fight", func(response http.ResponseWriter, request *http.Request) {
		body := agents.GetBytesBody(request)
		// unmarshal the json data
		var data map[string]string
		err := json.Unmarshal(body, &data)
		if err != nil {
			response.Write([]byte("Error: " + err.Error()))
		}
		info := data["info"]

		// QUESTION: how to add content type or somethinh help
		response.Write([]byte(info))

	})

	// Start the HTTP server
	fmt.Println("Starting HTTP server on port 5050...")
	err = bob.StartHttpServer()

	if err != nil {
		panic(err)
	}
}
