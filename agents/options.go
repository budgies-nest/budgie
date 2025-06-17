package agents

import (
	"context"

	"github.com/budgies-nest/budgie/enums/base"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func WithDMR(ctx context.Context, baseURL string) AgentOption {
	return func(agent *Agent) {
		agent.ctx = ctx
		agent.clientEngine = openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(""),
		)
	}
}

func WithOpenAI(ctx context.Context, apiKey string) AgentOption {
	return func(agent *Agent) {
		agent.ctx = ctx
		agent.clientEngine = openai.NewClient(
			option.WithBaseURL(base.OpenAIURL),
			option.WithAPIKey(apiKey),
		)
	}
}

func WithOpenAIURL(ctx context.Context, baseURL string, apiKey string) AgentOption {
	return func(agent *Agent) {
		agent.ctx = ctx
		agent.clientEngine = openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(apiKey),
		)
	}
}

// TODO: add more client options

// WithParams sets the parameters for the Agent's chat completion requests.
func WithParams(params openai.ChatCompletionNewParams) AgentOption {
	return func(agent *Agent) {
		agent.Params = params
	}
}

// WithEmbeddingParams sets the parameters for the Agent's embedding requests.
func WithEmbeddingParams(embeddingParams openai.EmbeddingNewParams) AgentOption {
	return func(agent *Agent) {
		agent.EmbeddingParams = embeddingParams
	}
}

// WithTools sets the tools for the Agent's chat completion requests.
// It allows the Agent to use specific tools during the chat completion process.
// IMPORTANT: The tools are appended to the existing tools in the Agent's parameters.
func WithTools(tools []openai.ChatCompletionToolParam) AgentOption {
	return func(agent *Agent) {
		agent.Params.Tools = append(agent.Params.Tools, tools...)
	}
}

// QUESTION: how to handle the MCP Tools?
