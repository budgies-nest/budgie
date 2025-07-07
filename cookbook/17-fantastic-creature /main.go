package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	chatAgentName := "Bob"
	chatAgentSystemInstructions := `
	You are a useful agent
	Your name is Bob
	`

	toolsAgent, err := agents.NewAgent("tools_agent",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(
			openai.ChatCompletionNewParams{
				Model: "ai/qwen2.5:latest",
				//Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
				Temperature:       openai.Opt(0.0), // IMPORTANT: set temperature to 0.0 to ensure the agent uses the tool
				ParallelToolCalls: openai.Bool(false),
			},
		),
		agents.WithTools(toolsIndex()),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(toolsAgent.Name, "is ready to assist!", toolsAgent.Params.Tools)

	// Enable global logging at Info level
	//agents.EnableLogging(agents.LogLevelInfo)

	chatAgent, err := agents.NewAgent(chatAgentName,
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(openai.ChatCompletionNewParams{
			//Model: "ai/qwen2.5:1.5B-F16",
			Model:       "ai/qwen2.5:latest",
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

			fmt.Println(strings.Repeat("=", 50))
			fmt.Println("👋 Before Chat Completion Stream")
			fmt.Println("Agent Name:", ctx.Agent.Name)
			fmt.Println("Agent Model:", ctx.Agent.Params.Model)
			fmt.Println(strings.Repeat("=", 50))
			fmt.Println("Last user message:")
			msgContent, _ := ctx.Agent.GetLastUserMessageContent()
			fmt.Println("📝 Content:", msgContent)

			fmt.Println(strings.Repeat("=", 50))

			// Display the messages list
			displayMessagesList(ctx.Agent.Params.Messages)

			toolsAgent.AddUserMessage(msgContent)
			// NOTE: it's simpler to call it with the fight API? or not ...
			detectedToolCalls, err := toolsAgent.ToolsCompletion(context.Background())
			// Generate the tools detection completion
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

		}),

		agents.WithAfterChatCompletionStream(func(ctx *agents.ChatCompletionStreamContext) {
			answer := *ctx.Response
			fmt.Println(strings.Repeat("=", 50))
			fmt.Println("👋 After Chat Completion Stream", answer)
			fmt.Println(strings.Repeat("=", 50))

			// Display the messages list
			displayMessagesList(ctx.Agent.Params.Messages)

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

func displayMessagesList(messages []openai.ChatCompletionMessageParamUnion) {
	for i, message := range messages {
		msg := map[string]string{
			"role":    "unknown",
			"content": "",
		}
		msg, _ = helpers.MessageToMap(message)

		switch {
		case message.OfUser != nil:
			fmt.Printf("🟢 Message %d (User): %s\n", i, msg["content"])
		case message.OfSystem != nil:
			fmt.Printf("🟣 Message %d (System): %s\n", i, msg["content"])
		case message.OfAssistant != nil:
			fmt.Printf("🟡 Message %d (Assistant): %s\n", i, msg["content"])
		case message.OfTool != nil:
			fmt.Printf("🟢 Message %d (Tool): %s\n", i, msg["content"])
		default:
			fmt.Printf("🔴 Message %d (Unknown): %s\n", i, msg["content"])
		}
	}
}

func toolsIndex() []openai.ChatCompletionToolParam {
	hello := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "hello",
			Description: openai.String("Say hello to the user"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{
						"type":        "string",
						"description": "The name of the user to greet.",
					},
				},
				"required": []string{"name"},
			},
		},
	}
	return []openai.ChatCompletionToolParam{
		hello,
	}
}
