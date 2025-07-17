package main

import (
	"context"
	"fmt"
	"os"

	"github.com/budgies-nest/budgie/agents"
	"github.com/openai/openai-go"
)

type Werewolf struct {
	Health       float64
	Strength     float64
	Agility      float64
	Intelligence float64
}

func main() {

	werewolf := Werewolf{
		Health:       100,
		Strength:     80,
		Agility:      70,
		Intelligence: 60,
	}

	modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")

	if modelRunnerBaseUrl == "" {
		panic("MODEL_RUNNER_BASE_URL environment variable is not set")
	}

	modelRunnerToolsModel := os.Getenv("MODEL_RUNNER_TOOLS_MODEL")
	if modelRunnerToolsModel == "" {
		panic("MODEL_RUNNER_TOOLS_MODEL environment variable is not set")
	}
	fmt.Println("Using Model Runner Tools Model:", modelRunnerToolsModel)

	toolsAgent, err := agents.NewAgent("tools_agent",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(
			openai.ChatCompletionNewParams{
				Model:             modelRunnerToolsModel,
				Temperature:       openai.Opt(0.0), // IMPORTANT: always set temperature to 0.0 for tools agents
				ParallelToolCalls: openai.Bool(true),
				Messages:          []openai.ChatCompletionMessageParamUnion{},
			},
		),
		agents.WithTools(toolsCatalog()),
	)
	if err != nil {
		panic(err)
	}

	msgContent := `
	- what is your health value?
	- set your health value to 200
	- increase your health by 10
	- decrease your health by 5
	- what is your intelligence value?
	- set your intelligence value to 100
	- decrease your intelligence by 5
	`
	toolsAgent.AddUserMessage(msgContent)

	detectedToolCalls, err := toolsAgent.ToolsCompletion(context.Background())

	if err != nil {
		fmt.Println("🔴 Error when detecting tool calls:", err)
		fmt.Println("🔍 No tool calls detected.")
		panic("Failed to detect tool calls")
	}

	numberOfToolCalls := len(detectedToolCalls)
	if numberOfToolCalls == 0 {
		fmt.Println("🔍 No tool calls detected.")
		panic("No tool calls detected")
	}

	fmt.Println("🔍 Detected tool calls:", len(detectedToolCalls))

	// TOOL CALLS:
	results, err := toolsAgent.ExecuteToolCalls(
		detectedToolCalls,
		toolsImplementation(&werewolf),
	)
	if err != nil {
		fmt.Println("🔴 Error executing tool calls:", err)
		panic("Failed to execute tool calls")
	}

	fmt.Println("🔧 Tool calls executed successfully. Results:", results)

}

func toolsCatalog() []openai.ChatCompletionToolParam {

	getHealth := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "get_health",
			Description: openai.String("Get the health of the Werewolf"),
			Parameters:  openai.FunctionParameters{},
		},
	}

	setHealth := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "set_health",
			Description: openai.String("Set the health of the Werewolf"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"value": map[string]string{
						"type":        "number",
						"description": "The new health value for the Werewolf.",
					},
				},
				"required": []string{"value"},
			},
		},
	}

	increaseHealth := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "increase_health",
			Description: openai.String("Increase the health of the Werewolf"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"amount": map[string]string{
						"type":        "number",
						"description": "The amount to increase the Werewolf's health by.",
					},
				},
				"required": []string{"amount"},
			},
		},
	}

	decreaseHealth := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "decrease_health",
			Description: openai.String("Decrease the health of the Werewolf"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"amount": map[string]string{
						"type":        "number",
						"description": "The amount to decrease the Werewolf's health by.",
					},
				},
				"required": []string{"amount"},
			},
		},
	}

	getIntelligence := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "get_intelligence",
			Description: openai.String("Get the intelligence of the Werewolf"),
			Parameters:  openai.FunctionParameters{},
		},
	}

	setIntelligence := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "set_intelligence",
			Description: openai.String("Set the intelligence of the Werewolf"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"value": map[string]string{
						"type":        "number",
						"description": "The new intelligence value for the Werewolf.",
					},
				},
				"required": []string{"value"},
			},
		},
	}

	increaseIntelligence := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "increase_intelligence",
			Description: openai.String("Increase the intelligence of the Werewolf"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"amount": map[string]string{
						"type":        "number",
						"description": "The amount to increase the Werewolf's intelligence by.",
					},
				},
				"required": []string{"amount"},
			},
		},
	}

	decreaseIntelligence := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "decrease_intelligence",
			Description: openai.String("Decrease the intelligence of the Werewolf"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"amount": map[string]string{
						"type":        "number",
						"description": "The amount to decrease the Werewolf's intelligence by.",
					},
				},
				"required": []string{"amount"},
			},
		},
	}

	return []openai.ChatCompletionToolParam{
		getHealth, setHealth, increaseHealth, decreaseHealth,
		getIntelligence, setIntelligence, increaseIntelligence, decreaseIntelligence,
		// Add more tools as needed
	}
}

// TODO: check the arguments provided to the tool calls
func toolsImplementation(werewolf *Werewolf) map[string]func(any) (any, error) {
	return map[string]func(any) (any, error){
		"get_health": func(args any) (any, error) {
			fmt.Println("🔧 Executing tool call: get_health with args:", args)
			return fmt.Sprintf("TELL THIS TO THE USER: 🐺 The Werewolf's health is %f.", werewolf.Health), nil
		},
		"set_health": func(args any) (any, error) {
			fmt.Println("🔧 Executing tool call: set_health with args:", args)
			newHealth := args.(map[string]any)["value"].(float64)
			werewolf.Health = newHealth
			return fmt.Sprintf("TELL THIS TO THE USER: 🐺 The Werewolf's health has been set to %f.", werewolf.Health), nil
		},
		"increase_health": func(args any) (any, error) {
			fmt.Println("🔧 Executing tool call: increase_health with args:", args)
			amount := args.(map[string]any)["amount"].(float64)
			werewolf.Health += amount
			return fmt.Sprintf("TELL THIS TO THE USER: 🐺 The Werewolf's health has been increased by %f. New health is %f.", amount, werewolf.Health), nil
		},
		"decrease_health": func(args any) (any, error) {
			fmt.Println("🔧 Executing tool call: decrease_health with args:", args)
			amount := args.(map[string]any)["amount"].(float64)
			werewolf.Health -= amount
			return fmt.Sprintf("TELL THIS TO THE USER: 🐺 The Werewolf's health has been decreased by %f. New health is %f.", amount, werewolf.Health), nil
		},
		"get_intelligence": func(args any) (any, error) {
			fmt.Println("🔧 Executing tool call: get_intelligence with args:", args)
			return fmt.Sprintf("TELL THIS TO THE USER: 🐺 The Werewolf's intelligence is %f.", werewolf.Intelligence), nil
		},
		"set_intelligence": func(args any) (any, error) {
			fmt.Println("🔧 Executing tool call: set_intelligence with args:", args)
			newIntelligence := args.(map[string]any)["value"].(float64)
			werewolf.Intelligence = newIntelligence
			return fmt.Sprintf("TELL THIS TO THE USER: 🐺 The Werewolf's intelligence has been set to %f.", werewolf.Intelligence), nil
		},
		"increase_intelligence": func(args any) (any, error) {
			fmt.Println("🔧 Executing tool call: increase_intelligence with args:", args)
			amount := args.(map[string]any)["amount"].(float64)
			werewolf.Intelligence += amount
			return fmt.Sprintf("TELL THIS TO THE USER: 🐺 The Werewolf's intelligence has been increased by %f. New intelligence is %f.", amount, werewolf.Intelligence), nil
		},
		"decrease_intelligence": func(args any) (any, error) {
			fmt.Println("🔧 Executing tool call: decrease_intelligence with args:", args)
			amount := args.(map[string]any)["amount"].(float64)
			werewolf.Intelligence -= amount
			return fmt.Sprintf("TELL THIS TO THE USER: 🐺 The Werewolf's intelligence has been decreased by %f. New intelligence is %f.", amount, werewolf.Intelligence), nil
		},
	}
}
