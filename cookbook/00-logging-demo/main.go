package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/enums/base"
	"github.com/openai/openai-go"
)

func main() {
	// Enable global logging at Info level
	agents.EnableLogging(agents.LogLevelInfo)

	// Test chat completion with logging
	fmt.Println("=== Testing Chat Completion with Logging ===")
	
	chatAgent, err := agents.NewAgent("chat-demo-agent",
		agents.WithDMR(base.DockerModelRunnerContainerURL),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
			Temperature: openai.Opt(0.0),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a helpful assistant that responds briefly."),
				openai.UserMessage("What is 2+2?"),
			},
		}),
		agents.WithLoggingEnabled(),
		agents.WithLogLevel(agents.LogLevelDebug),
	)
	if err != nil {
		fmt.Printf("Error creating chat agent: %v\n", err)
		return
	}

	response, err := chatAgent.ChatCompletion(context.Background())
	if err != nil {
		fmt.Printf("Error in chat completion: %v\n", err)
	} else {
		fmt.Printf("Response: %s\n", response)
	}

	// Test with tools and logging
	fmt.Println("\n=== Testing Tools with Logging ===")
	
	// Add a simple calculator tool
	calculatorTool := openai.ChatCompletionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "calculator",
			Description: openai.String("Performs basic arithmetic operations"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"operation": map[string]interface{}{
						"type":        "string",
						"description": "The operation to perform (add, subtract, multiply, divide)",
					},
					"a": map[string]interface{}{
						"type":        "number",
						"description": "First number",
					},
					"b": map[string]interface{}{
						"type":        "number",
						"description": "Second number",
					},
				},
				"required": []string{"operation", "a", "b"},
			},
		},
	}

	toolAgent, err := agents.NewAgent("tool-demo-agent",
		agents.WithDMR(base.DockerModelRunnerContainerURL),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
			Temperature: openai.Opt(0.0), // Important for tool detection
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a calculator assistant. Use the calculator tool to perform arithmetic."),
				openai.UserMessage("Calculate 15 * 3 using the calculator tool"),
			},
			ParallelToolCalls: openai.Bool(false), // Better for small models
		}),
		agents.WithTools([]openai.ChatCompletionToolParam{calculatorTool}),
	)
	if err != nil {
		fmt.Printf("Error creating tool agent: %v\n", err)
		return
	}

	// Tool implementation
	tools := map[string]func(any) (any, error){
		"calculator": func(args any) (any, error) {
			argsMap := args.(map[string]interface{})
			operation := argsMap["operation"].(string)
			a := argsMap["a"].(float64)
			b := argsMap["b"].(float64)

			switch operation {
			case "add":
				return a + b, nil
			case "subtract":
				return a - b, nil
			case "multiply":
				return a * b, nil
			case "divide":
				if b == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return a / b, nil
			default:
				return nil, fmt.Errorf("unknown operation: %s", operation)
			}
		},
	}

	toolCalls, err := toolAgent.ToolsCompletion(context.Background())
	if err != nil {
		fmt.Printf("Error in tools completion: %v\n", err)
		return
	}

	responses, err := toolAgent.ExecuteToolCalls(toolCalls, tools)
	if err != nil {
		fmt.Printf("Error executing tools: %v\n", err)
		return
	}

	for i, response := range responses {
		fmt.Printf("Tool response %d: %s\n", i+1, response)
	}

	// Test error logging
	fmt.Println("\n=== Testing Error Logging ===")
	agents.GetGlobalLogger().LogError("demo-agent", "test_error", "This is a test error", fmt.Errorf("simulated error"), map[string]interface{}{
		"context": "demonstration",
		"user_id": 12345,
	})

	// Test streaming
	fmt.Println("\n=== Testing Stream Completion with Logging ===")
	
	streamAgent, err := agents.NewAgent("stream-demo-agent",
		agents.WithDMR(base.DockerModelRunnerContainerURL),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
			Temperature: openai.Opt(0.8),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a helpful assistant."),
				openai.UserMessage("Tell me a short joke about programming"),
			},
		}),
	)
	if err != nil {
		fmt.Printf("Error creating stream agent: %v\n", err)
		return
	}

	fmt.Print("Streaming response: ")
	streamResponse, err := streamAgent.ChatCompletionStream(context.Background(), func(self *agents.Agent, content string, err error) error {
		if err != nil {
			return err
		}
		fmt.Print(content)
		return nil
	})
	if err != nil {
		fmt.Printf("Error in stream completion: %v\n", err)
	} else {
		fmt.Printf("\nFull response: %s\n", streamResponse)
	}

	// Disable logging
	fmt.Println("\n=== Disabling Logging ===")
	agents.DisableLogging()
	
	// This should not be logged
	chatAgent.AddUserMessage("This message should not be logged")
	_, err = chatAgent.ChatCompletion(context.Background())
	if err != nil {
		fmt.Printf("Error in chat completion: %v\n", err)
	}

	fmt.Println("Demo completed!")
}