package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/budgies-nest/budgie/agents"

	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func main() {
	modelRunnerBaseUrl := helpers.GetModelRunnerBaseUrl()

	//modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")
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
	fmt.Println("🛠️", toolsAgent.Name, "is ready to assist!", toolsAgent.Params.Tools)

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
			// NOTE: Start a tools detection

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

			// TOOLS DETECTION:
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

			// TOOL CALLS:
			results, err := toolsAgent.ExecuteToolCalls(detectedToolCalls,
				map[string]func(any) (any, error){
					"hello": func(args any) (any, error) {
						name := args.(map[string]any)["name"].(string)
						return fmt.Sprintf("👋 Hello %s! 🙂", name), nil
					},
				},
			)
			if err != nil {
				fmt.Println("Error executing tool calls:", err)
				return
			}
			fmt.Println("Tool Call Results:\n", results)

			//ctx.Agent.AddToolMessage(detectedToolCalls[0].ID, results[0]) // Assuming the first result is the one we want to use
			//ctx.Agent.AddUserMessage(results[0])
			ctx.Agent.AddSystemMessage("TELL THIS TO THE USER:" + results[0])
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

	chatAgent.HttpServer().HandleFunc("POST /api/info", func(response http.ResponseWriter, request *http.Request) {
		body := agents.GetBytesBody(request)
		// unmarshal the json data
		var data map[string]string
		err := json.Unmarshal(body, &data)
		if err != nil {
			response.Write([]byte("Error: " + err.Error()))
		}
		info := data["info"]

		response.Write([]byte("🤖: " + info))

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
