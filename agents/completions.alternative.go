package agents

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
)
// NOTE: this is subject to change in the future, as we are still experimenting with the best way to handle tool calls detection.
func (agent *Agent) AltenativeToolsCompletion(ctx context.Context) ([]openai.ChatCompletionMessageToolCall, error) {

	systemContentIntroduction := `You have access to the following tools:`
	catalog := agent.Params.Tools

	toolsJson, err := json.Marshal(catalog)
	if err != nil {
		return nil, errors.New("error marshalling tools to JSON: " + err.Error())
	}
	
	toolsContent := "[AVAILABLE_TOOLS]" + string(toolsJson) + "[/AVAILABLE_TOOLS]"

	systemContentInstructions := `If the question of the user matched the description of a tool, the tool will be called.
	To call a tool, respond with a JSON object with the following structure: 
	[
		{
			"name": <name of the called tool>,
			"arguments": {
				<name of the argument>: <value of the argument>
			}
		},
	]
	
	search the name of the tool in the list of tools with the Name field
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
		return nil, err

	}
	if len(completion.Choices) == 0 {
	}
	result := completion.Choices[0].Message.Content
	if result == "" {
		agent.Params.Tools = catalog // Restore the tools in case of error
		return nil, errors.New("no tool calls detected")
	}

	agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("Return all function calls wrapped in a container object with a 'function_calls' key."),
		openai.UserMessage(result),
	}

	agent.Params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONObject: &openai.ResponseFormatJSONObjectParam{
			Type: "json_object",
		},
	}

	completionNext, err := agent.clientEngine.Chat.Completions.New(ctx, agent.Params)

	if err != nil {
		agent.Params.Tools = catalog // Restore the tools in case of error
		return nil, errors.New("error in the next step of tool calls completion: " + err.Error())
	}

	if len(completionNext.Choices) == 0 {
		agent.Params.Tools = catalog // Restore the tools in case of error
		return nil, errors.New("no choices found in the next step of tool calls completion")
	}
	resultNext := completionNext.Choices[0].Message.Content
	if resultNext == "" {
		agent.Params.Tools = catalog // Restore the tools in case of error
		return nil, errors.New("no tool calls detected in the next step")
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
		return nil, errors.New("error unmarshalling tool calls JSON: " + errJson.Error())
	}
	if len(commands.FunctionCalls) == 0 {
		agent.Params.Tools = catalog // Restore the tools in case of error
		return nil, errors.New("no tool calls detected after unmarshalling")
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
			return nil, errors.New("error marshalling command arguments to JSON: " + err.Error())
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

	agent.Params.Tools = catalog // Restore the tools in case of error

	return toolCalls, nil
}
