package agents

import (
	"github.com/openai/openai-go"
)

func (agent *Agent) AddTool(tool openai.ChatCompletionToolParam) {
	agent.Params.Tools = append(agent.Params.Tools, tool)
}

func (agent *Agent) AddTools(tools []openai.ChatCompletionToolParam) {
	agent.Params.Tools = append(agent.Params.Tools, tools...)
}

