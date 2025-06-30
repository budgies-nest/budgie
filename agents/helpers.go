package agents

import "github.com/openai/openai-go"

func (agent *Agent) AddUserMessage(content string) {
	agent.Params.Messages = append(agent.Params.Messages, openai.UserMessage(content))
}

func (agent *Agent) AddSystemMessage(content string) {
	agent.Params.Messages = append(agent.Params.Messages, openai.SystemMessage(content))
}

func (agent *Agent) AddAssistantMessage(content string) {
	agent.Params.Messages = append(agent.Params.Messages, openai.AssistantMessage(content))
}

func (agent *Agent) AddTool(tool openai.ChatCompletionToolParam) {
	agent.Params.Tools = append(agent.Params.Tools, tool)
}

func (agent *Agent) AddTools(tools []openai.ChatCompletionToolParam) {
	agent.Params.Tools = append(agent.Params.Tools, tools...)
}

func (agent *Agent) SetModel(model string) {
	agent.Params.Model = model
}

func (agent *Agent) SetTemperature(temperature float64) {
	agent.Params.Temperature = openai.Opt(temperature)
}

func (agent *Agent) SetMaxTokens(maxTokens int64) {
	agent.Params.MaxTokens = openai.Opt(maxTokens)
}

func WithOpenAIClient(apiKey, baseURL string) AgentOption {
	return WithOpenAIURL(baseURL, apiKey)
}

func WithModel(model string) AgentOption {
	return func(agent *Agent) {
		agent.Params.Model = model
	}
}

func WithSystemInstructions(instructions string) AgentOption {
	return func(agent *Agent) {
		agent.Params.Messages = append(agent.Params.Messages, openai.SystemMessage(instructions))
	}
}

func WithTemperature(temperature float64) AgentOption {
	return func(agent *Agent) {
		agent.Params.Temperature = openai.Opt(temperature)
	}
}

func WithMaxTokens(maxTokens int64) AgentOption {
	return func(agent *Agent) {
		agent.Params.MaxTokens = openai.Opt(maxTokens)
	}
}