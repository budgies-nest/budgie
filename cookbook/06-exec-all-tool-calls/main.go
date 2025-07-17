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

func main() {
	modelRunnerBaseUrl := getModelRunnerBaseUrl()

	addTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "add",
			Description: openai.String("add two numbers"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]string{
						"type":        "number",
						"description": "The first number to add.",
					},
					"b": map[string]string{
						"type":        "number",
						"description": "The second number to add.",
					},
				},
				"required": []string{"a", "b"},
			},
		},
	}

	sayHelloTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "say_hello",
			Description: openai.String("Say hello to the given person name"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"name"},
			},
		},
	}

	/*
		Some small model are able to detect several tool calls in a single request.
		So let ParallelToolCalls set to openai.Bool(true), then the model will detect all the tool calls in a single request.

		Models that are able to detect several tool calls in a single request:
		- ignaciolopezluna020/llama-xlam:8B-Q4_K_M
		- k33g/llama-xlam-2:8b-fc-r-q2_k
		  - https://huggingface.co/Salesforce/Llama-xLAM-2-8b-fc-r-gguf
		  - Llama-xLAM-2-8B-fc-r-Q2_K.gguf

	*/
	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(
			openai.ChatCompletionNewParams{
				//Model: "k33g/llama-xlam-2:8b-fc-r-q2_k",
				Model: "hf.co/salesforce/xlam-2-3b-fc-r-gguf:q3_k_l",
				// NOTE: this model is able to detect several tool calls in a single request
				// NOTE: but it is bad with ParallelToolCalls: openai.Bool(false)
				Temperature: openai.Opt(0.0), // IMPORTANT: set temperature to 0.0 to ensure the agent uses the tool
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.UserMessage(`
						Add 10 and 32
						Add 12 and 30
						Say Hello to Bob
						Add 40 and 2
						Add 5 and 37					
					`),
				},
				ParallelToolCalls: openai.Bool(true),
				// TODO: make more test cases with this option
				// NOTE: with small models, parallel tool calls may not work as expected
			},
		),
		agents.WithTools([]openai.ChatCompletionToolParam{addTool, sayHelloTool}),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("🤖 Bob is ready to assist!", bob.Params.Tools)

	// Generate the tools detection completion
	detectedToolCalls, err := bob.ToolsCompletion(context.Background())
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

	results, err := bob.ExecuteToolCalls(detectedToolCalls,
		map[string]func(any) (any, error){

			"add": func(args any) (any, error) {
				a := args.(map[string]any)["a"].(float64)
				b := args.(map[string]any)["b"].(float64)
				return a + b, nil
			},

			"say_hello": func(args any) (any, error) {
				name := args.(map[string]any)["name"].(string)
				return fmt.Sprintf("Hello, %s!", name), nil
			},
		},
	)
	if err != nil {
		fmt.Println("Error executing tool calls:", err)
		return
	}
	fmt.Println("Results of Tool Calls:\n", results)

}

func getModelRunnerBaseUrl() string {
	// Detect if running in a container or locally
	if helpers.DetectContainerEnvironment() == environments.Local {
		return base.DockerModelRunnerLocalURL
	}
	return base.DockerModelRunnerContainerURL
}
