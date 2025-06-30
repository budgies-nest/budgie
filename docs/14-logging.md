# Logging

Budgie provides a comprehensive logging system that allows you to track completions, tool executions, and errors. The logging system is configurable and can be enabled or disabled as needed.

## Features

- **Configurable log levels**: Off, Error, Info, Debug
- **Enable/disable toggle**: Can be activated/deactivated at runtime
- **Structured JSON output**: Consistent format with timestamps and metadata
- **Multiple log types**: Chat completions, streaming, tools, and errors
- **Performance tracking**: Duration and response metrics
- **Global and per-agent configuration**: Flexible logging setup

## Log Levels

```go
// Available log levels
agents.LogLevelOff    // No logging
agents.LogLevelError  // Only errors
agents.LogLevelInfo   // Info and errors
agents.LogLevelDebug  // All logs including debug info
```

## Global Logging Configuration

### Enable/Disable Logging

```go
package main

import (
    "github.com/budgies-nest/budgie/agents"
)

func main() {
    // Enable logging globally at Info level
    agents.EnableLogging(agents.LogLevelInfo)
    
    // Disable logging globally
    agents.DisableLogging()
    
    // Get current global logger
    logger := agents.GetGlobalLogger()
    
    // Set custom global logger
    customLogger := agents.NewLogger(agents.LogLevelDebug, true)
    agents.SetGlobalLogger(customLogger)
}
```

## Per-Agent Logging Configuration

### Using Agent Options

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
    // Agent with logging enabled
    agent, err := agents.NewAgent("my-agent",
        agents.WithDMR(base.DockerModelRunnerContainerURL),
        agents.WithParams(openai.ChatCompletionNewParams{
            Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
            Temperature: openai.Opt(0.0),
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You are a helpful assistant."),
                openai.UserMessage("What is 2+2?"),
            },
        }),
        agents.WithLoggingEnabled(),
        agents.WithLogLevel(agents.LogLevelDebug),
    )
    if err != nil {
        panic(err)
    }
    
    // This will be logged
    response, err := agent.ChatCompletion(context.Background())
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Response: %s\n", response)
}
```

### Using Custom Logger

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
    // Create custom logger
    customLogger := agents.NewLogger(agents.LogLevelInfo, true)
    
    agent, err := agents.NewAgent("custom-logger-agent",
        agents.WithDMR(base.DockerModelRunnerContainerURL),
        agents.WithParams(openai.ChatCompletionNewParams{
            Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
            Temperature: openai.Opt(0.0),
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You are a helpful assistant."),
                openai.UserMessage("Hello!"),
            },
        }),
        agents.WithLogger(customLogger),
    )
    if err != nil {
        panic(err)
    }
    
    response, err := agent.ChatCompletion(context.Background())
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Response: %s\n", response)
}
```


## Runtime Log Control

You can control logging at runtime:

```go
// Check if logging is enabled
if agents.GetGlobalLogger().IsEnabled() {
    fmt.Println("Logging is currently enabled")
}

// Change log level
agents.GetGlobalLogger().SetLevel(agents.LogLevelError)

// Enable/disable
agents.GetGlobalLogger().SetEnabled(false)
```

## Best Practices

1. **Use appropriate log levels**: 
   - `LogLevelError` for production environments
   - `LogLevelInfo` for development and debugging
   - `LogLevelDebug` for detailed troubleshooting

2. **Monitor performance**: The logging system tracks duration for all operations

3. **Structured data**: All logs are JSON formatted for easy parsing and analysis

4. **Global vs per-agent**: Use global logging for consistent behavior, per-agent for specific requirements

5. **Runtime control**: Enable/disable logging based on environment or user preferences

## Environment Variables

You can also control logging via environment variables in your applications:

```go
func setupLogging() {
    if os.Getenv("BUDGIE_LOGGING") == "true" {
        level := agents.LogLevelInfo
        if os.Getenv("BUDGIE_LOG_LEVEL") == "debug" {
            level = agents.LogLevelDebug
        } else if os.Getenv("BUDGIE_LOG_LEVEL") == "error" {
            level = agents.LogLevelError
        }
        agents.EnableLogging(level)
    }
}
```

The logging system provides comprehensive visibility into your agent's operations while maintaining performance and configurability.