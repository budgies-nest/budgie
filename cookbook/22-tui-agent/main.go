package main

import (
	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func main() {

	modelRunnerBaseUrl := helpers.GetModelRunnerBaseUrl()

	// Create a new agent named Bob
	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       "ai/qwen2.5:latest",
			Temperature: openai.Opt(0.8),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You're a helpful assistant expert with Star Trek universe."),
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	// Start the TUI prompt with custom messages
	err = bob.Prompt(agents.PromptConfig{
		UseStreamCompletion:        true, // Set to false for non-streaming completion
		StartingMessage:            "🖖 Welcome to the Star Trek Assistant!",
		ExplanationMessage:         "Ask me anything about the Star Trek universe. Type '/bye' to quit or Ctrl+C to interrupt responses.",
		PromptTitle:                "🚀 Star Trek Query",
		ThinkingPrompt:             "🤖 ",
		InterruptInstructions:      "(Press Ctrl+C to interrupt)",
		CompletionInterruptMessage: "⚠️ Response was interrupted\n",
		GoodbyeMessage:             "🖖 Live long and prosper!",
	})
	if err != nil {
		panic(err)
	}
}
