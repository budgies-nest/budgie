# Messages Helpers

The messages helpers provide convenient methods for managing chat messages in an agent's conversation context.

## Adding Messages

### AddUserMessage
```golang
agent.AddUserMessage("Hello, how are you?")
```

### AddSystemMessage
```golang
agent.AddSystemMessage("You are a helpful assistant.")
```

### AddAssistantMessage
```golang
agent.AddAssistantMessage("I'm doing well, thank you!")
```

### AddToolMessage
```golang
agent.AddToolMessage("tool_call_123", "Tool execution result")
```

## Retrieving Messages

### GetMessages
```golang
messages := agent.GetMessages()
// Returns all messages as []openai.ChatCompletionMessageParamUnion
```

### GetLastUserMessage
```golang
lastMsg := agent.GetLastUserMessage()
// Returns pointer to last user message or nil if none found
```

### GetLastUserMessageContent
```golang
content, err := agent.GetLastUserMessageContent()
if err != nil {
    // Handle error
}
// Returns content string of last user message
```

### GetLastAssistantMessage
```golang
lastMsg := agent.GetLastAssistantMessage()
// Returns pointer to last assistant message or nil if none found
```

### GetLastAssistantMessageContent
```golang
content, err := agent.GetLastAssistantMessageContent()
if err != nil {
    // Handle error
}
// Returns content string of last assistant message
```

## Removing Messages

### RemoveMessage
```golang
agent.RemoveMessage(2) // Remove message at index 2
```

### RemoveLastUserMessage
```golang
agent.RemoveLastUserMessage()
// Removes the most recent user message
```

### RemoveLastAssistantMessage
```golang
agent.RemoveLastAssistantMessage()
// Removes the most recent assistant message
```

### ClearMessages
```golang
agent.ClearMessages()
// Removes all messages from the agent
```

## Example Usage

```golang
// Create agent with initial system message
agent, err := agents.NewAgent("Assistant",
    agents.WithParams(
        openai.ChatCompletionNewParams{
            Model: "gpt-3.5-turbo",
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You are a helpful coding assistant."),
            },
        }
    ),
)

// Add user question
agent.AddUserMessage("How do I create a Go struct?")

// Get response
response, err := agent.ChatCompletion(context.Background())
if err != nil {
    panic(err)
}

// Add assistant's response to conversation
agent.AddAssistantMessage(response)

// Continue conversation
agent.AddUserMessage("Can you show me an example?")

// Get last user message content
lastQuestion, err := agent.GetLastUserMessageContent()
if err != nil {
    panic(err)
}
println("User asked:", lastQuestion)

// Remove last message if needed
agent.RemoveLastUserMessage()
```