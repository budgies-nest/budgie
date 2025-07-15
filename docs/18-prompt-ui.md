# Prompt UI - Interactive Terminal Chat Interface

The `Prompt` method provides an interactive Terminal User Interface (TUI) that allows users to chat with an agent in real-time. This creates a conversational experience similar to ChatGPT but running directly in your terminal.

## Overview

The Prompt method creates a chat loop where users can:
- Type messages and receive responses from the agent
- Use streaming or non-streaming completions
- Interrupt responses with Ctrl+C
- Exit gracefully with `/bye` command
- Maintain conversation history throughout the session

## Basic Usage

```go
package main

import (
    "github.com/budgies-nest/budgie/agents"
    "github.com/budgies-nest/budgie/openai"
)

func main() {
    // Create an agent
    agent, err := agents.NewAgent("Assistant",
        agents.WithDMR("http://localhost:3000"),
        agents.WithParams(openai.ChatCompletionNewParams{
            Model:       "ai/qwen2.5:latest",
            Temperature: openai.Opt(0.8),
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You're a helpful assistant."),
            },
        }),
    )
    if err != nil {
        panic(err)
    }

    // Start interactive prompt
    err = agent.Prompt(agents.PromptConfig{
        UseStreamCompletion: true,
    })
    if err != nil {
        panic(err)
    }
}
```

## PromptConfig Structure

The `PromptConfig` struct allows you to customize the TUI experience:

```go
type PromptConfig struct {
    UseStreamCompletion        bool   // Enable streaming responses
    StartingMessage            string // Welcome message
    ExplanationMessage         string // User instructions
    PromptTitle                string // Input prompt title
    ThinkingPrompt             string // Processing indicator
    InterruptInstructions      string // Interrupt help text
    CompletionInterruptMessage string // Interruption message
    GoodbyeMessage             string // Exit message
}
```

### Configuration Fields

| Field | Default | Description |
|-------|---------|-------------|
| `UseStreamCompletion` | `false` | Enable streaming responses for real-time output |
| `StartingMessage` | `"🤖 Starting TUI for agent: {name}"` | Message displayed when TUI starts |
| `ExplanationMessage` | `"Type your questions below. Use '/bye' to quit or Ctrl+C to interrupt completions."` | Instructions shown to user |
| `PromptTitle` | `"💬 Chat with {name}"` | Title displayed in input prompt |
| `ThinkingPrompt` | `"🤔 "` | Indicator shown while agent is processing |
| `InterruptInstructions` | `"(Press Ctrl+C to interrupt)"` | Instructions for interrupting responses |
| `CompletionInterruptMessage` | `"🚫 Completion was interrupted\n"` | Message shown when response is interrupted |
| `GoodbyeMessage` | `"👋 Goodbye!"` | Message shown when user exits with `/bye` |

## Advanced Usage

### Custom Configuration

```go
err = agent.Prompt(agents.PromptConfig{
    UseStreamCompletion:        true,
    StartingMessage:            "🖖 Welcome to the Star Trek Assistant!",
    ExplanationMessage:         "Ask me anything about the Star Trek universe. Type '/bye' to quit or Ctrl+C to interrupt responses.",
    PromptTitle:                "🚀 Star Trek Query",
    ThinkingPrompt:             "🤖 ",
    InterruptInstructions:      "(Press Ctrl+C to interrupt)",
    CompletionInterruptMessage: "⚠️ Response was interrupted\n",
    GoodbyeMessage:             "🖖 Live long and prosper!",
})
```

### Streaming vs Non-Streaming

**Streaming Mode** (`UseStreamCompletion: true`):
- Responses appear character by character as they're generated
- More interactive and responsive feeling
- Can be interrupted mid-response with Ctrl+C
- Uses `agent.ChatCompletionStream()`

**Non-Streaming Mode** (`UseStreamCompletion: false`):
- Complete response appears at once
- Less interactive but more stable
- Uses `agent.ChatCompletion()`

## Features

### Interactive Commands

- **Regular Input**: Type any message and press Enter to send
- **Exit Command**: Type `/bye` to quit the session gracefully
- **Interrupt**: Press Ctrl+C to interrupt ongoing completions

### Conversation Management

The Prompt method automatically:
- Adds user messages to the agent's conversation history
- Adds assistant responses to the conversation history
- Maintains context throughout the session
- Preserves conversation state between messages

### Error Handling

- Graceful handling of interruptions
- Context cancellation support
- Proper cleanup of signal handling
- Error messages for completion failures

## Implementation Details

The Prompt method is implemented in `agents/tui-prompt.go` and:

1. **Initializes UI**: Sets up default messages if not provided
2. **Creates Forms**: Uses the `huh` library for interactive input
3. **Handles Input**: Processes user messages and special commands
4. **Manages Completions**: Supports both streaming and non-streaming modes
5. **Provides Interruption**: Allows Ctrl+C to stop ongoing completions
6. **Maintains History**: Automatically manages conversation state

## Related Methods

The Prompt method works with several other agent methods:

- `AddUserMessage(content string)` - Adds user input to conversation
- `AddAssistantMessage(content string)` - Adds assistant response
- `ChatCompletion(ctx context.Context)` - Synchronous completion
- `ChatCompletionStream(ctx context.Context, callback)` - Streaming completion
- `GetMessages()` - Returns conversation history
- `ClearMessages()` - Clears conversation history

## Example: Complete Application

```go
package main

import (
    "log"
    "github.com/budgies-nest/budgie/agents"
    "github.com/budgies-nest/budgie/openai"
)

func main() {
    // Create a specialized agent
    codeReviewer, err := agents.NewAgent("CodeReviewer",
        agents.WithDMR("http://localhost:3000"),
        agents.WithParams(openai.ChatCompletionNewParams{
            Model:       "ai/qwen2.5:latest",
            Temperature: openai.Opt(0.3),
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You are a senior software engineer who reviews code and provides constructive feedback."),
            },
        }),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Start interactive session with custom configuration
    err = codeReviewer.Prompt(agents.PromptConfig{
        UseStreamCompletion:        true,
        StartingMessage:            "👨‍💻 Code Review Assistant Ready!",
        ExplanationMessage:         "Paste your code or describe what you need reviewed. Type '/bye' to quit.",
        PromptTitle:                "📝 Code Review",
        ThinkingPrompt:             "🔍 ",
        InterruptInstructions:      "(Press Ctrl+C to interrupt review)",
        CompletionInterruptMessage: "⚠️ Review was interrupted\n",
        GoodbyeMessage:             "👋 Happy coding!",
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

This creates a specialized code review assistant with a customized interface that provides an interactive experience for code reviews.