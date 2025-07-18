package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func main() {
	modelRunnerBaseUrl := helpers.GetModelRunnerBaseUrl()

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(
			openai.ChatCompletionNewParams{
				//Model:       "ai/qwen2.5:latest",
				// REF: https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF
				Model:       "hf.co/bartowski/smollm2-135m-instruct-gguf:q4_k_m",
				Temperature: openai.Opt(0.0),
				Messages:    []openai.ChatCompletionMessageParamUnion{},
			},
		),
		agents.WithA2AServer(agents.A2AServerConfig{Port: "8888"}),
		agents.WithAgentCard(agents.AgentCard{
			Name:        "Bob",
			Description: "A helpful assistant with expertise in the Star Trek universe.",
			URL:         "http://localhost:8888",
			Version:     "1.0.0",
			//Capabilities: map[string]any{},
			Skills: []map[string]any{
				{
					"id":          "ask_for_something",
					"name":        "Ask for something",
					"description": "Bob is using a small language model to answer questions",
				},
				{
					"id":          "say_hello_world",
					"name":        "Say Hello World",
					"description": "Bob can say hello world",
				},
			},
		}),
		agents.WithAgentCallback(func(ctx *agents.AgentCallbackContext) (agents.TaskResponse, error) {

			fmt.Printf("🟢 Processing task request: %s\n", ctx.TaskRequest.ID)
			// Extract user message
			userMessage := ctx.TaskRequest.Params.Message.Parts[0].Text
			fmt.Printf("🔵 UserMessage: %s\n", userMessage)
			fmt.Printf("🟡 TaskRequest Metadata: %v\n", ctx.TaskRequest.Params.MetaData)

			var systemMessage, userPrompt string

			switch ctx.TaskRequest.Params.MetaData["skill"] {
			case "ask_for_something":
				systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
				userPrompt = userMessage

			case "greetings":
				systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
				userPrompt = "Greetings to " + userMessage + " with emojis and use his name."

			default:
				systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
				userPrompt = "Be nice, and explain that " + fmt.Sprintf("%v", ctx.TaskRequest.Params.MetaData["skill"]) + " is not a valid task ID."
			}

			ctx.Agent.AddSystemMessage(systemMessage)
			ctx.Agent.AddUserMessage(userPrompt)

			answer, err := ctx.Agent.ChatCompletion(context.Background())
			if err != nil {
				fmt.Printf("❌ Error during chat completion: %v\n", err)
				return agents.TaskResponse{}, err
			}
			fmt.Printf("🤖 Generated response: %s\n", answer)

			// Create response task
			responseTask := agents.TaskResponse{
				ID:             ctx.TaskRequest.ID,
				JSONRpcVersion: "2.0",
				Result: agents.Result{
					Status: agents.TaskStatus{
						State: "completed",
					},
					History: []agents.AgentMessage{
						{
							Role: "assistant",
							Parts: []agents.TextPart{
								{
									Text: answer,
									Type: "text",
								},
							},
						},
					},
					Kind:     "task",
					Metadata: map[string]any{},
				},
			}

			return responseTask, nil

		}),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("🤖 Starting A2A server on port", bob.A2AServerConfig().Port)

	errSrv := bob.StartA2AServer()
	if errSrv != nil {
		panic(errSrv)
	}

}
