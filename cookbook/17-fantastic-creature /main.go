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


func getModelRunnerBaseUrl() string {
	// Detect if running in a container or locally
	if helpers.DetectContainerEnvironment() == environments.Local {
		return base.DockerModelRunnerLocalURL
	}
	return base.DockerModelRunnerContainerURL
}

// TODO:
/*
add the necessary fields to handle the health, strength, ...
*/

func main() {

	modelRunnerBaseUrl := getModelRunnerBaseUrl()
	chatAgentName := "Werewolf"
	chatAgentSystemInstructions := `
	You are a werewolf
	Your name is Smity
	`

	// Enable global logging at Info level
	//agents.EnableLogging(agents.LogLevelInfo)

	chatAgent, err := agents.NewAgent(chatAgentName,
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model: "hf.co/unsloth/qwen3-0.6b-gguf:q4_k_m",
			Temperature: openai.Opt(0.8),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(chatAgentSystemInstructions),
				//openai.UserMessage("Who is James T Kirk?"),
			},
		}),
		agents.WithLoggingEnabled(),
		agents.WithLogLevel(agents.LogLevelError),
		agents.WithHTTPServer(agents.HTTPServerConfig{
			Port: "5050",
		}),
		agents.WithBeforeChatCompletionStream(func(ctx *agents.ChatCompletionStreamContext) {
			// TODO: make the tool detection here
			messages := ctx.Agent.Params.Messages
			// QUESTION: how to get the user message?
			for idx, message := range messages {
				content, _ := message.MarshalJSON()
				
				fmt.Println("🔴", idx, string(content))
			}
		}),
		agents.WithBeforeChatCompletionStream(func(ctx *agents.ChatCompletionStreamContext) {
			fmt.Println("👋 Hello World 🌎")
		}),
	)

	// QUESTION: how to handle the conversational memory?
	// NOTE: -> WithAfterCompletionStream (and Before) -> handle an array of messages

	if err != nil {
		panic(err)
	}
	// TODO: check how the context is handled with the REST API
	// _, err = bob.ChatCompletionStream(context.Background(), func(self *agents.Agent, content string, err error) error {
	// 	fmt.Print(content)
	// 	return nil
	// })

	chatAgent.HttpServer().HandleFunc("POST /api/fight", func(response http.ResponseWriter, request *http.Request) {
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
	err = chatAgent.StartHttpServer()

	if err != nil {
		panic(err)
	}
}
