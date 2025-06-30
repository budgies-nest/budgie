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

## Log Types and Output

### Chat Completion Logs

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "type": "chat_completion",
  "agent_name": "my-agent",
  "data": {
    "model": "gpt-4o-mini",
    "messages_count": 2,
    "max_tokens": null,
    "temperature": 0.7,
    "response_length": 42,
    "duration_ms": 1500
  },
  "message": "Chat completion successful"
}
```

### Streaming Chat Completion Logs

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "type": "chat_completion_stream",
  "agent_name": "my-agent",
  "data": {
    "model": "gpt-4o-mini",
    "messages_count": 2,
    "response_length": 156,
    "duration_ms": 2300
  },
  "message": "Chat completion stream successful"
}
```

### Tools Completion Logs

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "type": "tools_completion",
  "agent_name": "my-agent",
  "data": {
    "model": "gpt-4o-mini",
    "messages_count": 2,
    "tools_count": 3,
    "tool_calls_count": 1,
    "duration_ms": 1800,
    "tool_names": ["calculator"]
  },
  "message": "Tools completion successful with 1 tool calls"
}
```

### Tool Execution Logs

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "type": "tool_execution",
  "agent_name": "my-agent",
  "data": {
    "tool_name": "calculator",
    "response_length": 4,
    "duration_ms": 5,
    "args": {"operation": "add", "a": 2, "b": 3},
    "response": "5"
  },
  "message": "Tool 'calculator' executed successfully"
}
```

### MCP Tool Execution Logs

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "type": "mcp_tool_execution",
  "agent_name": "my-agent",
  "data": {
    "tool_name": "filesystem_read",
    "client_type": "stdio",
    "response_length": 1024,
    "duration_ms": 150
  },
  "message": "MCP tool 'filesystem_read' executed successfully via stdio"
}
```

### Error Logs

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "error",
  "type": "chat_completion",
  "agent_name": "my-agent",
  "data": {
    "model": "gpt-4o-mini",
    "messages_count": 2,
    "duration_ms": 500
  },
  "error": "API rate limit exceeded"
}
```

## Complete Example with Tools

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
    // Enable global logging
    agents.EnableLogging(agents.LogLevelDebug)
    
    // Add calculator tool
    calculatorTool := openai.ChatCompletionToolParam{
        Function: openai.FunctionDefinitionParam{
            Name:        "calculator",
            Description: openai.String("Performs basic arithmetic operations"),
            Parameters: openai.FunctionParameters{
                "type": "object",
                "properties": map[string]interface{}{
                    "operation": map[string]interface{}{
                        "type":        "string",
                        "description": "The operation (add, subtract, multiply, divide)",
                    },
                    "a": map[string]interface{}{
                        "type":        "number",
                        "description": "First number",
                    },
                    "b": map[string]interface{}{
                        "type":        "number",
                        "description": "Second number",
                    },
                },
                "required": []string{"operation", "a", "b"},
            },
        },
    }
    
    // Create agent with logging and local model
    agent, err := agents.NewAgent("calculator-agent",
        agents.WithDMR(base.DockerModelRunnerContainerURL),
        agents.WithParams(openai.ChatCompletionNewParams{
            Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
            Temperature: openai.Opt(0.0), // Important for tool detection
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You are a calculator assistant. Use the calculator tool to perform arithmetic."),
                openai.UserMessage("Calculate 15 * 3 using the calculator tool"),
            },
            ParallelToolCalls: openai.Bool(false), // Better for small models
        }),
        agents.WithTools([]openai.ChatCompletionToolParam{calculatorTool}),
    )
    if err != nil {
        panic(err)
    }
    
    // Tool implementation
    tools := map[string]func(any) (any, error){
        "calculator": func(args any) (any, error) {
            argsMap := args.(map[string]interface{})
            operation := argsMap["operation"].(string)
            a := argsMap["a"].(float64)
            b := argsMap["b"].(float64)
            
            switch operation {
            case "add":
                return a + b, nil
            case "subtract":
                return a - b, nil
            case "multiply":
                return a * b, nil
            case "divide":
                if b == 0 {
                    return nil, fmt.Errorf("division by zero")
                }
                return a / b, nil
            default:
                return nil, fmt.Errorf("unknown operation: %s", operation)
            }
        },
    }
    
    // This will log the tools completion
    toolCalls, err := agent.ToolsCompletion(context.Background())
    if err != nil {
        fmt.Printf("Error in tools completion: %v\n", err)
        return
    }
    
    // This will log each tool execution
    responses, err := agent.ExecuteToolCalls(toolCalls, tools)
    if err != nil {
        fmt.Printf("Error executing tools: %v\n", err)
        return
    }
    
    for i, response := range responses {
        fmt.Printf("Tool response %d: %s\n", i+1, response)
    }
    
    // Log custom error
    agents.GetGlobalLogger().LogError("calculator-agent", "validation_error", 
        "Invalid input detected", fmt.Errorf("user provided negative number"), 
        map[string]interface{}{
            "user_input": "-5",
            "operation": "sqrt",
        })
    
    // Disable logging for subsequent operations
    agents.DisableLogging()
    
    // This won't be logged
    agent.AddUserMessage("This won't be logged")
    agent.ChatCompletion(context.Background())
    
    fmt.Println("Demo completed!")
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