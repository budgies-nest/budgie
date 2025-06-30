package agents

import (
	"context"
	"errors"
	"time"

	"github.com/openai/openai-go"
)

// ChatCompletion handles the chat completion request using the DMR client.
// It sends the parameters set in the Agent and returns the response content or an error.
// It is a synchronous operation that waits for the completion to finish.
func (agent *Agent) ChatCompletion(ctx context.Context) (string, error) {
	start := time.Now()
	completion, err := agent.clientEngine.Chat.Completions.New(ctx, agent.Params)
	duration := time.Since(start)

	var response string
	var finalErr error

	if err != nil {
		finalErr = err
	} else if len(completion.Choices) > 0 {
		response = completion.Choices[0].Message.Content
	} else {
		finalErr = errors.New("no choices found")
	}

	agent.logger.LogChatCompletion(agent.Name, agent.Params, response, duration, finalErr)

	if finalErr != nil {
		return "", finalErr
	}
	return response, nil
}

// ChatCompletionStream handles the chat completion request using the DMR client in a streaming manner.
// It takes a callback function that is called for each chunk of content received.
// The callback function receives the Agent instance, the content of the chunk, and any error that occurred.
// It returns the accumulated response content and any error that occurred during the streaming process.
// The callback function should return an error if it wants to stop the streaming process.
func (agent *Agent) ChatCompletionStream(ctx context.Context, callBack func(self *Agent, content string, err error) error) (string, error) {
	start := time.Now()
	response := ""
	stream := agent.clientEngine.Chat.Completions.NewStreaming(ctx, agent.Params)
	var cbkRes error

	for stream.Next() {
		chunk := stream.Current()
		// Stream each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			cbkRes = callBack(agent, chunk.Choices[0].Delta.Content, nil)
			response += chunk.Choices[0].Delta.Content
		}

		if cbkRes != nil {
			break
		}
	}

	duration := time.Since(start)
	var finalErr error

	if cbkRes != nil {
		finalErr = cbkRes
	} else if err := stream.Err(); err != nil {
		finalErr = err
	} else if err := stream.Close(); err != nil {
		finalErr = err
	}

	agent.logger.LogChatCompletionStream(agent.Name, agent.Params, response, duration, finalErr)

	if finalErr != nil {
		return response, finalErr
	}
	return response, nil
}

// ToolsCompletion handles the tool calls completion request using the DMR client.
// It sends the parameters set in the Agent and returns the detected tool calls or an error.
// It is a synchronous operation that waits for the completion to finish.
func (agent *Agent) ToolsCompletion(ctx context.Context) ([]openai.ChatCompletionMessageToolCall, error) {
	start := time.Now()
	completion, err := agent.clientEngine.Chat.Completions.New(ctx, agent.Params)
	duration := time.Since(start)

	var detectedToolCalls []openai.ChatCompletionMessageToolCall
	var finalErr error

	if err != nil {
		finalErr = err
	} else {
		detectedToolCalls = completion.Choices[0].Message.ToolCalls
		if len(detectedToolCalls) == 0 {
			finalErr = errors.New("no tool calls detected")
		}
	}

	agent.logger.LogToolsCompletion(agent.Name, agent.Params, detectedToolCalls, duration, finalErr)

	if finalErr != nil {
		return nil, finalErr
	}
	return detectedToolCalls, nil
}
