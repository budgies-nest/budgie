package agents

// Chat completion handler options

// WithBeforeChatCompletion adds a before chat completion handler
func WithBeforeChatCompletion(handler BeforeChatCompletionHandler) AgentOption {
	return func(agent *Agent) {
		agent.completionHandlers.BeforeChatCompletion = append(
			agent.completionHandlers.BeforeChatCompletion, handler)
	}
}

// WithAfterChatCompletion adds an after chat completion handler
func WithAfterChatCompletion(handler AfterChatCompletionHandler) AgentOption {
	return func(agent *Agent) {
		agent.completionHandlers.AfterChatCompletion = append(
			agent.completionHandlers.AfterChatCompletion, handler)
	}
}

// Chat completion stream handler options

// WithBeforeChatCompletionStream adds a before chat completion stream handler
func WithBeforeChatCompletionStream(handler BeforeChatCompletionStreamHandler) AgentOption {
	return func(agent *Agent) {
		agent.completionHandlers.BeforeChatCompletionStream = append(
			agent.completionHandlers.BeforeChatCompletionStream, handler)
	}
}

// WithAfterChatCompletionStream adds an after chat completion stream handler
func WithAfterChatCompletionStream(handler AfterChatCompletionStreamHandler) AgentOption {
	return func(agent *Agent) {
		agent.completionHandlers.AfterChatCompletionStream = append(
			agent.completionHandlers.AfterChatCompletionStream, handler)
	}
}

// Tools completion handler options

// WithBeforeToolsCompletion adds a before tools completion handler
func WithBeforeToolsCompletion(handler BeforeToolsCompletionHandler) AgentOption {
	return func(agent *Agent) {
		agent.completionHandlers.BeforeToolsCompletion = append(
			agent.completionHandlers.BeforeToolsCompletion, handler)
	}
}

// WithAfterToolsCompletion adds an after tools completion handler
func WithAfterToolsCompletion(handler AfterToolsCompletionHandler) AgentOption {
	return func(agent *Agent) {
		agent.completionHandlers.AfterToolsCompletion = append(
			agent.completionHandlers.AfterToolsCompletion, handler)
	}
}

// Alternative tools completion handler options

// WithBeforeAlternativeToolsCompletion adds a before alternative tools completion handler
func WithBeforeAlternativeToolsCompletion(handler BeforeAlternativeToolsCompletionHandler) AgentOption {
	return func(agent *Agent) {
		agent.completionHandlers.BeforeAlternativeToolsCompletion = append(
			agent.completionHandlers.BeforeAlternativeToolsCompletion, handler)
	}
}

// WithAfterAlternativeToolsCompletion adds an after alternative tools completion handler
func WithAfterAlternativeToolsCompletion(handler AfterAlternativeToolsCompletionHandler) AgentOption {
	return func(agent *Agent) {
		agent.completionHandlers.AfterAlternativeToolsCompletion = append(
			agent.completionHandlers.AfterAlternativeToolsCompletion, handler)
	}
}