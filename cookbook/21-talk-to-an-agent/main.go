package main

import (
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

func main() {
	modelRunnerBaseUrl := helpers.GetModelRunnerBaseUrl()

	sam, err := agents.NewAgent("Bob",
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
	)

	if err != nil {
		panic(err)
	}

	agentBaseURL := "http://0.0.0.0:8888"

	agentCard, err := sam.PingAgent(agentBaseURL)
	if err != nil {
		fmt.Println("Error pinging agent:", err)
		return
	}
	jsonAgentCard, err := agents.AgentCardToJSONString(agentCard)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("🤖 Agent Card:\n", jsonAgentCard)

	taskRequest := agents.TaskRequest{
		ID:     uuid.NewString(),
		Method: "message/send",
		Params: agents.AgentMessageParams{
			Message: agents.AgentMessage{
				Role: "user",
				Parts: []agents.TextPart{
					{
						Text: "What is the best pizza in the world?",
					},
				},
			},
			MetaData: map[string]any{
				"skill": "ask_for_something",
			},
		},
	}
	taskResponse, err := sam.SendToAgent(agentBaseURL, taskRequest)
	if err != nil {
		fmt.Println("Error sending task request:", err)
		return
	}

	jsonTaskResponse, err := agents.TaskResponseToJSONString(taskResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("🟢 Task Response JSON:\n", jsonTaskResponse)

	fmt.Println("🟣 Task Response Text:", taskResponse.Result.History[0].Parts[0].Text)

}
