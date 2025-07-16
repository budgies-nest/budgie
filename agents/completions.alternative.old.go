package agents

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)
// IMPORTANT:
// NOTE: I keep this file for reference, but I will not use it anymore.
// TODO: Remove this file in the future.
func (agent *Agent) OldAlternativeToolsCompletion(ctx context.Context) ([]openai.ChatCompletionMessageToolCall, error) {
	start := time.Now()

	// Create context for handlers
	handlerCtx := &AlternativeToolsCompletionContext{
		CompletionContext: CompletionContext{
			Agent:     agent,
			Context:   ctx,
			StartTime: start,
		},
	}

	// Call before handlers
	for _, handler := range agent.completionHandlers.BeforeAlternativeToolsCompletion {
		handler(handlerCtx)
	}

	//systemContentIntroduction := `You have access to the following tools:`

	systemContentIntroduction := `You are an AI assistant with access to various tools. Your task is to analyze user input and identify ALL possible tool calls that can be made.
	IMPORTANT: You must process the ENTIRE user input and identify ALL tool calls, not just the first few. Each line or request in the user input should be analyzed separately.
	You have access to the following tools:
	`

	catalog := agent.Params.Tools

	toolsJson, err := json.Marshal(catalog)
	if err != nil {
		finalErr := errors.New("error marshalling tools to JSON: " + err.Error())
		duration := time.Since(start)
		handlerCtx.Duration = duration
		handlerCtx.Error = finalErr
		for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
			handler(handlerCtx)
		}
		agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, finalErr)
		return nil, finalErr
	}

	toolsContent := "[AVAILABLE_TOOLS]" + string(toolsJson) + "[/AVAILABLE_TOOLS]"

	// systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
	// To call a tool, respond with a JSON object with the following structure:
	// [
	// 	{
	// 		"name": <name of the called tool>,
	// 		"arguments": {
	// 			<name of the argument>: <value of the argument>
	// 		}
	// 	},
	// ]

	// search the name of the tool in the list of tools with the Name field
	// `

	systemContentInstructions := `INSTRUCTIONS:
	1. Read the ENTIRE user input carefully
	2. Process each line/request separately
	3. For each request, check if it matches any tool description
	4. If multiple tool calls are needed, include ALL of them in your response
	5. NEVER stop processing until you've analyzed the complete input

	TOOL MATCHING RULES:
	- Match tool calls based on the "description" field of each tool
	- Use the exact "name" field from the tool definition
	- Provide all required arguments as specified in the tool's parameters

	RESPONSE FORMAT:
	When you find tool calls, respond with a JSON array containing ALL identified tool calls:
	[
		{
			"name": "<exact_tool_name_from_catalog>",
			"arguments": {
				"<parameter_name>": "<parameter_value>"
			}
		},
		{
			"name": "<next_tool_name>",
			"arguments": {
				"<parameter_name>": "<parameter_value>"
			}
		}
	]

	EXAMPLES:
	Input: "Say hello to John. Add 5 and 10. Make vulcan salute to Spock."
	Output: [
		{"name": "send_message", "arguments": {"name": "John"}},
		{"name": "operation", "arguments": {"number1": 5, "number2": 10, "number3": 8}},
		{"name": "greetings", "arguments": {"name": "Jane"}}
	]

	If no tool calls are found, respond with an empty array: []

	CRITICAL: You must analyze the COMPLETE user input and identify ALL possible tool calls. Do not stop after finding the first few matches.
	`

	instructionMessages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemContentIntroduction + "\n" + toolsContent + "\n" + systemContentInstructions),
	}

	currentAgentMessages := agent.Params.Messages

	// Combine messages using append
	newSetOfMessages := append([]openai.ChatCompletionMessageParamUnion{}, instructionMessages...)
	newSetOfMessages = append(newSetOfMessages, currentAgentMessages...)

	// Add the user message to the new set of messages
	agent.Params.Messages = newSetOfMessages

	// IMPORTANT: Deactivate the tools for the next step of the completion
	agent.Params.Tools = nil
	// IMPORTANT: at the end of the function, we will restore the tools to the original state with the catalog variable

	completion, err := agent.clientEngine.Chat.Completions.New(ctx, agent.Params)
	if err != nil {
		agent.Params.Tools = catalog // Restore the tools in case of error
		duration := time.Since(start)
		handlerCtx.Duration = duration
		handlerCtx.Error = err
		for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
			handler(handlerCtx)
		}
		agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, err)
		return nil, err
	}
	if len(completion.Choices) == 0 {
	}
	result := completion.Choices[0].Message.Content
	if result == "" {
		agent.Params.Tools = catalog // Restore the tools in case of error
		finalErr := errors.New("no tool calls detected")
		duration := time.Since(start)
		handlerCtx.Duration = duration
		handlerCtx.Error = finalErr
		for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
			handler(handlerCtx)
		}
		agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, finalErr)
		return nil, finalErr
	}

	agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
		//openai.SystemMessage("Return all function calls wrapped in a container object with a 'function_calls' key."),
		openai.UserMessage(result),
	}
	/*
		🚀 Next Result: {"function_calls": [
		{"name": "vulcan_salute", "arguments": {"name": "Spock"}},
		{"name": "say_hello", "arguments": {"name": "John Doe"}},
		{"name": "addition", "arguments": {"number1": 10, "number2": 32}},
		{"name": "vulcan_salute", "arguments": {"name": "Bob Morane"}},
		{"name": "say_hello", "arguments": {"name": "Jane Doe"}},
		{"name": "who_is", "arguments": {"name": "Jean-Luc Picard"}},
		{"name": "say_hello", "arguments": {"name": "Philippe"}},
		{"name": "addition", "arguments": {"number1": 5, "number2": 37}},
		{"name": "vulcan_salute", "arguments": {"name": "Sam"}},
		{"name": "say_hello", "arguments": {"name": "Alice"}},
		{"name": "vulcan_salute", "arguments": {"name": "Bob"}}
		]}
	*/

	responseFormat := openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
			Type: "json_schema",
			JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:        "function_calls",
				Description: openai.String("Function calls data structure"),
				Schema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"function_calls": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"name": map[string]any{
										"type":        "string",
										"description": "The name of the function to call",
									},
									"arguments": map[string]any{
										"type":        "object",
										"description": "The arguments to pass to the function",
									},
								},
								"required":             []string{"name", "arguments"},
								"additionalProperties": false,
							},
							"description": "Array of function calls to execute",
						},
					},
					"required":             []string{"function_calls"},
					"additionalProperties": false,
				},
			},
		},
	}

	// agent.Params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
	// 	OfJSONObject: &openai.ResponseFormatJSONObjectParam{
	// 		Type: "json_object",
	// 	},
	// }

	agent.Params.ResponseFormat = responseFormat

	completionNext, err := agent.clientEngine.Chat.Completions.New(ctx, agent.Params)

	if err != nil {
		agent.Params.Tools = catalog // Restore the tools in case of error
		finalErr := errors.New("error in the next step of tool calls completion: " + err.Error())
		duration := time.Since(start)
		handlerCtx.Duration = duration
		handlerCtx.Error = finalErr
		for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
			handler(handlerCtx)
		}
		agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, finalErr)
		return nil, finalErr
	}

	if len(completionNext.Choices) == 0 {
		agent.Params.Tools = catalog // Restore the tools in case of error
		finalErr := errors.New("no choices found in the next step of tool calls completion")
		duration := time.Since(start)
		handlerCtx.Duration = duration
		handlerCtx.Error = finalErr
		for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
			handler(handlerCtx)
		}
		agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, finalErr)
		return nil, finalErr
	}
	resultNext := completionNext.Choices[0].Message.Content
	if resultNext == "" {
		agent.Params.Tools = catalog // Restore the tools in case of error
		finalErr := errors.New("no tool calls detected in the next step")
		duration := time.Since(start)
		handlerCtx.Duration = duration
		handlerCtx.Error = finalErr
		for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
			handler(handlerCtx)
		}
		agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, finalErr)
		return nil, finalErr
	}

	type Command struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}

	type FunctionCalls struct {
		FunctionCalls []Command `json:"function_calls"`
	}

	//var commands []Command
	var commands FunctionCalls

	errJson := json.Unmarshal([]byte(resultNext), &commands)
	if errJson != nil {
		agent.Params.Tools = catalog // Restore the tools in case of error
		finalErr := errors.New("error unmarshalling tool calls JSON: " + errJson.Error())
		duration := time.Since(start)
		handlerCtx.Duration = duration
		handlerCtx.Error = finalErr
		for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
			handler(handlerCtx)
		}
		agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, finalErr)
		return nil, finalErr
	}
	if len(commands.FunctionCalls) == 0 {
		agent.Params.Tools = catalog // Restore the tools in case of error
		finalErr := errors.New("no tool calls detected after unmarshalling")
		duration := time.Since(start)
		handlerCtx.Duration = duration
		handlerCtx.Error = finalErr
		for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
			handler(handlerCtx)
		}
		agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, finalErr)
		return nil, finalErr
	}

	// Create a []openai.ChatCompletionMessageToolCall from the commands
	toolCalls := make([]openai.ChatCompletionMessageToolCall, len(commands.FunctionCalls))

	/*
	   // The ID of the tool call.
	   ID string `json:"id,required"`
	   // The function that the model called.
	   Function ChatCompletionMessageToolCallFunction `json:"function,required"`
	   // The type of the tool. Currently, only `function` is supported.
	   Type constant.Function `json:"type,required"`

	*/

	for i, command := range commands.FunctionCalls {
		// transform command.Arguments to  JSON string
		argumentsJson, err := json.Marshal(command.Arguments)
		if err != nil {
			agent.Params.Tools = catalog // Restore the tools in case of error
			finalErr := errors.New("error marshalling command arguments to JSON: " + err.Error())
			duration := time.Since(start)
			handlerCtx.Duration = duration
			handlerCtx.Error = finalErr
			for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
				handler(handlerCtx)
			}
			agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, nil, duration, finalErr)
			return nil, finalErr
		}

		toolCalls[i] = openai.ChatCompletionMessageToolCall{
			ID: uuid.New().String(), // Generate a unique ID for the tool call
			Function: openai.ChatCompletionMessageToolCallFunction{
				Name:      command.Name,
				Arguments: string(argumentsJson),
			},
			Type: "function",
		}
	}

	agent.Params.Tools = catalog // Restore the tools

	duration := time.Since(start)

	// Update handler context with results
	handlerCtx.Duration = duration
	handlerCtx.Error = nil
	handlerCtx.ToolCalls = &toolCalls

	// Call after handlers
	for _, handler := range agent.completionHandlers.AfterAlternativeToolsCompletion {
		handler(handlerCtx)
	}

	// Add logging for AlternativeToolsCompletion (was missing)
	agent.logger.LogAlternativeToolsCompletion(agent.Name, agent.Params, toolCalls, duration, nil)

	return toolCalls, nil
}
