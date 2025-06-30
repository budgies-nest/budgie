# Completion Handlers - Quick Start Guide

> Completion handlers allow you to hook into the agent completion lifecycle to modify inputs and outputs.

## What Are Completion Handlers?

Completion handlers are functions that execute **before** or **after** AI completions, allowing you to:

- **Enhance prompts** before they reach the AI
- **Modify responses** after the AI generates them
- **Monitor performance** and collect metrics
- **Handle errors** gracefully
- **Implement custom logic** like filtering, caching, or analytics

## Quick Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/budgies-nest/budgie/agents"
    "github.com/budgies-nest/budgie/enums/base"
    "github.com/openai/openai-go"
)

func main() {
    agent, err := agents.NewAgent("MyAgent",
        // Standard configuration
        agents.WithDMR(base.DockerModelRunnerContainerURL),
        agents.WithParams(openai.ChatCompletionNewParams{
            Model: "k33g/qwen2.5:0.5b-instruct-q8_0",
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You are a helpful assistant."),
                openai.UserMessage("What is machine learning?"),
            },
        }),
        
        // 🔧 Before handler: Add context to improve responses
        agents.WithBeforeChatCompletion(func(ctx *agents.ChatCompletionContext) {
            contextMsg := openai.SystemMessage(
                "CONTEXT: Please provide detailed, educational answers with examples.")
            ctx.Agent.Params.Messages = append(
                []openai.ChatCompletionMessageParamUnion{contextMsg},
                ctx.Agent.Params.Messages...)
            
            fmt.Println("🎨 Enhanced prompt with educational context")
        }),
        
        // ✨ After handler: Add decorative elements
        agents.WithAfterChatCompletion(func(ctx *agents.ChatCompletionContext) {
            if ctx.Response != nil && ctx.Error == nil {
                original := *ctx.Response
                *ctx.Response = "🤖 " + original + "\n\n💡 Enhanced by completion handler!"
                
                fmt.Printf("⏱️  Completion took: %v\n", ctx.Duration)
                fmt.Printf("📝 Response length: %d characters\n", len(*ctx.Response))
            }
        }),
    )
    
    if err != nil {
        panic(err)
    }
    
    response, err := agent.ChatCompletion(context.Background())
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Final Response: %s\n", response)
}
```

## Handler Types

### Before Handlers
Execute **before** the AI model is called:
- `WithBeforeChatCompletion()` - Modify chat parameters
- `WithBeforeChatCompletionStream()` - Modify streaming parameters  
- `WithBeforeToolsCompletion()` - Modify tool parameters
- `WithBeforeAlternativeToolsCompletion()` - Modify alternative tool parameters

### After Handlers  
Execute **after** the AI model responds:
- `WithAfterChatCompletion()` - Modify chat responses
- `WithAfterChatCompletionStream()` - Modify streaming responses
- `WithAfterToolsCompletion()` - Modify tool results
- `WithAfterAlternativeToolsCompletion()` - Modify alternative tool results

## Context Objects

Each handler receives a rich **context object** with access to:

```go
type ChatCompletionContext struct {
    Agent     *Agent          // The agent instance
    Context   context.Context // Go context for cancellation
    StartTime time.Time       // When completion started
    Duration  time.Duration   // How long it took
    Error     error          // Any error that occurred
    Response  *string        // The AI response (modifiable)
}
```

## Common Use Cases

### 1. Response Enhancement
```go
agents.WithAfterChatCompletion(func(ctx *agents.ChatCompletionContext) {
    if ctx.Response != nil && ctx.Error == nil {
        *ctx.Response = "✨ " + *ctx.Response + " ✨"
    }
})
```

### 2. Content Filtering
```go
agents.WithAfterChatCompletion(func(ctx *agents.ChatCompletionContext) {
    if ctx.Response != nil && ctx.Error == nil {
        filtered := strings.ReplaceAll(*ctx.Response, "inappropriate", "***")
        *ctx.Response = filtered
    }
})
```

### 3. Performance Monitoring
```go
agents.WithAfterChatCompletion(func(ctx *agents.ChatCompletionContext) {
    fmt.Printf("Completion took %v\n", ctx.Duration)
    if ctx.Duration > 5*time.Second {
        fmt.Println("⚠️  Slow completion detected")
    }
})
```

### 4. Error Handling
```go
agents.WithAfterChatCompletion(func(ctx *agents.ChatCompletionContext) {
    if ctx.Error != nil {
        fallback := "I'm sorry, I'm having technical difficulties. Please try again."
        ctx.Response = &fallback
        ctx.Error = nil // Clear error
    }
})
```

### 5. Prompt Enhancement
```go
agents.WithBeforeChatCompletion(func(ctx *agents.ChatCompletionContext) {
    // Add context based on user's question
    contextMsg := openai.SystemMessage("Please provide step-by-step explanations.")
    ctx.Agent.Params.Messages = append(
        []openai.ChatCompletionMessageParamUnion{contextMsg},
        ctx.Agent.Params.Messages...)
})
```

## Multiple Handlers

You can register multiple handlers - they execute in registration order:

```go
agent, err := agents.NewAgent("MyAgent",
    // First handler: Add prefix
    agents.WithAfterChatCompletion(func(ctx *agents.ChatCompletionContext) {
        if ctx.Response != nil && ctx.Error == nil {
            *ctx.Response = "AI: " + *ctx.Response
        }
    }),
    // Second handler: Add metadata
    agents.WithAfterChatCompletion(func(ctx *agents.ChatCompletionContext) {
        if ctx.Response != nil && ctx.Error == nil {
            *ctx.Response += fmt.Sprintf("\n[Processed in %v]", ctx.Duration)
        }
    }),
)
```

## Why Context Pattern?

Instead of just passing the agent, handlers receive a **rich context object** that provides:

- **Lifecycle awareness** - Know if you're before/after completion
- **Rich metadata** - Access timing, errors, response data
- **Type safety** - Different contexts for different completion types
- **Extensibility** - Easy to add new information
- **Clean architecture** - Separation of concerns

The completion handlers system opens up powerful possibilities for customizing agent behavior without modifying core completion logic!