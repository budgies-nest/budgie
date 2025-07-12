package agents

import (
	"context"
	"time"

	"github.com/openai/openai-go"
)

// CompletionContext provides context information for completion handlers
type CompletionContext struct {
	Agent     *Agent
	Context   context.Context
	StartTime time.Time
	Duration  time.Duration
	Error     error
}

// ChatCompletionContext provides specific context for chat completion handlers
type ChatCompletionContext struct {
	CompletionContext
	Response *string // Pointer to allow modification
}

// ChatCompletionStreamContext provides specific context for streaming chat completion handlers
type ChatCompletionStreamContext struct {
	CompletionContext
	Response *string // Pointer to allow modification
	Callback func(self *Agent, content string, err error) error
}

// ToolsCompletionContext provides specific context for tools completion handlers
type ToolsCompletionContext struct {
	CompletionContext
	ToolCalls *[]openai.ChatCompletionMessageToolCall // Pointer to allow modification
}

// AlternativeToolsCompletionContext provides specific context for alternative tools completion handlers
type AlternativeToolsCompletionContext struct {
	CompletionContext
	ToolCalls *[]openai.ChatCompletionMessageToolCall // Pointer to allow modification
}

// AgentCallbackContext provides specific context for agent callback handlers
type AgentCallbackContext struct {
	CompletionContext
	TaskRequest  *TaskRequest  // Pointer to allow modification
	TaskResponse *TaskResponse // Pointer to allow modification
}

// Handler types for before and after completion events
type (
	// Chat completion handlers
	BeforeChatCompletionHandler func(*ChatCompletionContext)
	AfterChatCompletionHandler  func(*ChatCompletionContext)

	// Chat completion stream handlers
	BeforeChatCompletionStreamHandler func(*ChatCompletionStreamContext)
	AfterChatCompletionStreamHandler  func(*ChatCompletionStreamContext)

	// Tools completion handlers
	BeforeToolsCompletionHandler func(*ToolsCompletionContext)
	AfterToolsCompletionHandler  func(*ToolsCompletionContext)

	// Alternative tools completion handlers
	BeforeAlternativeToolsCompletionHandler func(*AlternativeToolsCompletionContext)
	AfterAlternativeToolsCompletionHandler  func(*AlternativeToolsCompletionContext)
)

// CompletionHandlers holds all completion handlers for an agent
type CompletionHandlers struct {
	// Chat completion handlers
	BeforeChatCompletion []BeforeChatCompletionHandler
	AfterChatCompletion  []AfterChatCompletionHandler

	// Chat completion stream handlers
	BeforeChatCompletionStream []BeforeChatCompletionStreamHandler
	AfterChatCompletionStream  []AfterChatCompletionStreamHandler

	// Tools completion handlers
	BeforeToolsCompletion []BeforeToolsCompletionHandler
	AfterToolsCompletion  []AfterToolsCompletionHandler

	// Alternative tools completion handlers
	BeforeAlternativeToolsCompletion []BeforeAlternativeToolsCompletionHandler
	AfterAlternativeToolsCompletion  []AfterAlternativeToolsCompletionHandler
}

// NewCompletionHandlers creates a new CompletionHandlers instance
func NewCompletionHandlers() *CompletionHandlers {
	return &CompletionHandlers{
		BeforeChatCompletion:             make([]BeforeChatCompletionHandler, 0),
		AfterChatCompletion:              make([]AfterChatCompletionHandler, 0),
		BeforeChatCompletionStream:       make([]BeforeChatCompletionStreamHandler, 0),
		AfterChatCompletionStream:        make([]AfterChatCompletionStreamHandler, 0),
		BeforeToolsCompletion:            make([]BeforeToolsCompletionHandler, 0),
		AfterToolsCompletion:             make([]AfterToolsCompletionHandler, 0),
		BeforeAlternativeToolsCompletion: make([]BeforeAlternativeToolsCompletionHandler, 0),
		AfterAlternativeToolsCompletion:  make([]AfterAlternativeToolsCompletionHandler, 0),
	}
}