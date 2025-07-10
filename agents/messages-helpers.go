package agents

import (
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func (agent *Agent) AddUserMessage(content string) {
	agent.Params.Messages = append(agent.Params.Messages, openai.UserMessage(content))
}

func (agent *Agent) AddSystemMessage(content string) {
	agent.Params.Messages = append(agent.Params.Messages, openai.SystemMessage(content))
}

func (agent *Agent) AddAssistantMessage(content string) {
	agent.Params.Messages = append(agent.Params.Messages, openai.AssistantMessage(content))
}

func (agent *Agent) AddToolMessage(toolCallID, content string) {
	agent.Params.Messages = append(agent.Params.Messages, openai.ToolMessage(toolCallID, content))
}

func (agent *Agent) GetMessages() []openai.ChatCompletionMessageParamUnion {
	return agent.Params.Messages
}

func (agent *Agent) RemoveMessage(index int) {
	if index >= 0 && index < len(agent.Params.Messages) {
		agent.Params.Messages = append(agent.Params.Messages[:index], agent.Params.Messages[index+1:]...)
	}
}

func (agent *Agent) ClearMessages() {
	agent.Params.Messages = []openai.ChatCompletionMessageParamUnion{}
}

func (agent *Agent) GetLastUserMessage() *openai.ChatCompletionMessageParamUnion {
	for i := len(agent.Params.Messages) - 1; i >= 0; i-- {
		if agent.Params.Messages[i].OfUser != nil {
			return &agent.Params.Messages[i]
		}
	}
	return nil
}

func (agent *Agent) GetLastUserMessageContent() (string, error) {
	lastUserMessage := agent.GetLastUserMessage()
	if lastUserMessage == nil {
		return "", nil
	}

	msgMap, err := helpers.MessageToMap(*lastUserMessage)
	if err != nil {
		return "", err
	}

	return msgMap["content"], nil
}

func (agent *Agent) RemoveLastUserMessage() {
	for i := len(agent.Params.Messages) - 1; i >= 0; i-- {
		if agent.Params.Messages[i].OfUser != nil {
			agent.Params.Messages = append(agent.Params.Messages[:i], agent.Params.Messages[i+1:]...)
			return
		}
	}
}

func (agent *Agent) GetLastAssistantMessage() *openai.ChatCompletionMessageParamUnion {
	for i := len(agent.Params.Messages) - 1; i >= 0; i-- {
		if agent.Params.Messages[i].OfAssistant != nil {
			return &agent.Params.Messages[i]
		}
	}
	return nil
}

func (agent *Agent) GetLastAssistantMessageContent() (string, error) {
	lastAssistantMessage := agent.GetLastAssistantMessage()
	if lastAssistantMessage == nil {
		return "", nil
	}

	msgMap, err := helpers.MessageToMap(*lastAssistantMessage)
	if err != nil {
		return "", err
	}

	return msgMap["content"], nil
}

func (agent *Agent) RemoveLastAssistantMessage() {
	for i := len(agent.Params.Messages) - 1; i >= 0; i-- {
		if agent.Params.Messages[i].OfAssistant != nil {
			agent.Params.Messages = append(agent.Params.Messages[:i], agent.Params.Messages[i+1:]...)
			return
		}
	}
}
