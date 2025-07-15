# TUI Agent Example

This example demonstrates how to use the `bob.Prompt()` method to create an interactive Terminal User Interface (TUI) for chatting with an agent.

## Features

- Interactive chat interface using the `huh` library
- Support for both streaming and non-streaming completions
- Simple command handling (`/bye` to quit)
- Ctrl+C interrupt handling for stopping completions
- Conversation history maintained throughout the session

## Usage

```bash
go run main.go
```

## Configuration

The `PromptConfig` struct allows you to configure the behavior and messages:

```go
bob.Prompt(agents.PromptConfig{
    UseStreamCompletion: true, // Enable streaming for real-time responses
    StartingMessage:     "🖖 Welcome to the Star Trek Assistant!",
    ExplanationMessage:  "Ask me anything about the Star Trek universe. Type '/bye' to quit or Ctrl+C to interrupt responses.",
    PromptTitle:         "🚀 Star Trek Query",
})
```

### Configuration Options

- `UseStreamCompletion: true` - Enables streaming completions where responses appear character by character
- `UseStreamCompletion: false` - Uses regular completions where the full response appears at once
- `StartingMessage` - Custom welcome message (default: "🤖 Starting TUI for agent: {name}")
- `ExplanationMessage` - Custom instructions (default: "Type your questions below. Use '/bye' to quit or Ctrl+C to interrupt completions.")
- `PromptTitle` - Custom title for the input prompt (default: "💬 Chat with {name}")

## Commands

- Type any question or message to chat with the agent
- Type `/bye` to quit the application
- Press Ctrl+C during completion to interrupt the response and return to the prompt

## Example Session

```
🖖 Welcome to the Star Trek Assistant!
Ask me anything about the Star Trek universe. Type '/bye' to quit or Ctrl+C to interrupt responses.

🚀 Star Trek Query
┃ Tell me about Captain Picard

🤔 Captain Jean-Luc Picard is one of the most iconic characters in the Star Trek universe...

🚀 Star Trek Query
┃ What about Data?

🤔 Data is an android officer serving aboard the USS Enterprise...

🚀 Star Trek Query
┃ /bye

👋 Goodbye!
```