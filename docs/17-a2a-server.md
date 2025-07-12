# Creating an A2A Server Agent with Budgie

> **✋ IMPORTANT**: 
> - This A2A protocol implementation is a subset of the A2A specification.
> - This is a work in progress and may not cover all aspects of the A2A protocol.


This guide explains how to create an Agent-to-Agent (A2A) server using the Budgie framework. An A2A server allows your agent to communicate with other agents over HTTP using JSON-RPC protocol.

## Overview

An A2A server agent can:
- Expose itself as an HTTP server with agent discovery capabilities
- Handle task requests from other agents
- Process different types of skills/tasks based on metadata
- Return structured responses following the JSON-RPC 2.0 specification

## Basic Setup

### 1. Import Required Packages

```go
import (
    "context"
    "fmt"

    "github.com/budgies-nest/budgie/agents"
    "github.com/budgies-nest/budgie/enums/base"
    "github.com/budgies-nest/budgie/enums/environments"
    "github.com/budgies-nest/budgie/helpers"
    "github.com/openai/openai-go"
)
```

### 2. Create the Agent with A2A Server Configuration

```go
bob, err := agents.NewAgent("Bob",
    agents.WithDMR(modelRunnerBaseUrl),
    agents.WithParams(
        openai.ChatCompletionNewParams{
            Model: "hf.co/bartowski/smollm2-135m-instruct-gguf:q4_k_m",
            Temperature: openai.Opt(0.0),
            Messages:    []openai.ChatCompletionMessageParamUnion{},
        },
    ),
    agents.WithA2AServer(agents.A2AServerConfig{Port: "8888"}),
    // ... other configurations
)
```

### 3. Configure Agent Card

The agent card provides discovery information for other agents:

```go
agents.WithAgentCard(agents.AgentCard{
    Name:        "Bob",
    Description: "A helpful assistant with expertise in the Star Trek universe.",
    URL:         "http://localhost:8888",
    Version:     "1.0.0",
    Skills: []map[string]any{
        {
            "id":          "ask_for_something",
            "name":        "Ask for something",
            "description": "Bob is using a small language model to answer questions",
        },
        {
            "id":          "say_hello_world",
            "name":        "Say Hello World",
            "description": "Bob can say hello world",
        },
    },
})
```

### 4. Implement the Agent Callback

The callback function handles incoming task requests:

```go
agents.WithAgentCallback(func(ctx *agents.AgentCallbackContext) (agents.TaskResponse, error) {
    fmt.Printf("=> Processing task request: %s\n", ctx.TaskRequest.ID)
    
    // Extract user message
    userMessage := ctx.TaskRequest.Params.Message.Parts[0].Text
    fmt.Printf("=> UserMessage: %s\n", userMessage)
    fmt.Printf("=> TaskRequest Metadata: %v\n", ctx.TaskRequest.Params.MetaData)

    var systemMessage, userPrompt string

    // Handle different skills based on metadata
    switch ctx.TaskRequest.Params.MetaData["skill"] {
    case "ask_for_something":
        systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
        userPrompt = userMessage

    case "greetings":
        systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
        userPrompt = "Greetings to " + userMessage + " with emojis and use his name."

    default:
        systemMessage = "You are Bob, a simple A2A agent. You can answer questions."
        userPrompt = "Be nice, and explain that " + fmt.Sprintf("%v", ctx.TaskRequest.Params.MetaData["skill"]) + " is not a valid task ID."
    }

    // Process the request
    ctx.Agent.AddSystemMessage(systemMessage)
    ctx.Agent.AddUserMessage(userPrompt)

    answer, err := ctx.Agent.ChatCompletion(context.Background())
    if err != nil {
        fmt.Printf("L Error during chat completion: %v\n", err)
        return agents.TaskResponse{}, err
    }

    // Create response task
    responseTask := agents.TaskResponse{
        ID:             ctx.TaskRequest.ID,
        JSONRpcVersion: "2.0",
        Result: agents.Result{
            Status: agents.TaskStatus{
                State: "completed",
            },
            History: []agents.AgentMessage{
                {
                    Role: "assistant",
                    Parts: []agents.TextPart{
                        {
                            Text: answer,
                            Type: "text",
                        },
                    },
                },
            },
            Kind:     "task",
            Metadata: map[string]any{},
        },
    }

    return responseTask, nil
})
```

### 5. Start the Server

```go
fmt.Println("> Starting A2A server on port", bob.A2AServerConfig().Port)

errSrv := bob.StartA2AServer()
if errSrv != nil {
    panic(errSrv)
}
```

## Agent Discovery

Once your A2A server is running, other agents can discover it by making a GET request to:

```
http://localhost:8888/.well-known/agent.json
```

This returns the agent card with available skills and metadata.

## Sending Tasks to the Agent

Send tasks using JSON-RPC 2.0 format via POST requests:

```json
{
    "jsonrpc": "2.0",
    "id": "1111",
    "method": "message/send",
    "params": {
        "message": {
            "role": "user",
            "parts": [
                {
                    "text": "What is the best pizza in the world?"
                }
            ]
        },
        "metadata": {
            "skill": "ask_for_something"
        }
    }
}
```


A complete example can be found in `/cookbook/20-agent-as-a2a-server/`.