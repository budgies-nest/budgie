# Alternative tool call detection
> With small models, tool calls and parallel tool calls may not work as expected

> 🧪 Experimental feature

The SLMs (small Language Models) are not very good on multiple tool calls detection. They often fail to detect all the tool calls, or they detect it but do not provide the correct parameters.

The budgie agent provides an alternative tool call detection mechanism that can be used to detect the tool calls in the user query with the `AlternativeToolsCompletion` method.

## Define the tools

```go
addTool := openai.ChatCompletionToolParam{
    Function: openai.FunctionDefinitionParam{
        Name:        "add",
        Description: openai.String("add two numbers"),
        Parameters: openai.FunctionParameters{
            "type": "object",
            "properties": map[string]interface{}{
                "a": map[string]string{
                    "type":        "number",
                    "description": "The first number to add.",
                },
                "b": map[string]string{
                    "type":        "number",
                    "description": "The second number to add.",
                },
            },
            "required": []string{"a", "b"},
        },
    },
}

multiplyTool := openai.ChatCompletionToolParam{
    Function: openai.FunctionDefinitionParam{
        Name:        "multiply",
        Description: openai.String("multiply two numbers"),
        Parameters: openai.FunctionParameters{
            "type": "object",
            "properties": map[string]interface{}{
                "a": map[string]string{
                    "type":        "number",
                    "description": "The first number to multiply.",
                },
                "b": map[string]string{
                    "type":        "number",
                    "description": "The second number to multiply.",
                },
            },
            "required": []string{"a", "b"},
        },
    },
}


sayHelloTool := openai.ChatCompletionToolParam{
    Function: openai.FunctionDefinitionParam{
        Name:        "say_hello",
        Description: openai.String("Say hello to the given person name"),
        Parameters: openai.FunctionParameters{
            "type": "object",
            "properties": map[string]interface{}{
                "name": map[string]string{
                    "type": "string",
                },
            },
            "required": []string{"name"},
        },
    },
}
```

## Initialize the agent

```go
bob, err := agents.NewAgent("Bob",
    agents.WithDMR(base.DockerModelRunnerContainerURL),
    agents.WithParams(
        openai.ChatCompletionNewParams{
            Model: "ai/qwen2.5:latest",
            Temperature: openai.Opt(0.0), 
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.UserMessage(`
                    Add 10 and 32
                    Add 12 and 30
                    Say Hello to Bob
                    Add 40 and 2
                    Add 5 and 37
                    Say Hey to John Doe
                    Multiply 2 and 3					
                `),
            },
        },
    ),
    agents.WithTools([]openai.ChatCompletionToolParam{addTool, sayHelloTool, multiplyTool}),
)
```

## Run the alternative completion to detect the tool calls

```go
detectedToolCalls, err := bob.AlternativeToolsCompletion() 
```

## Print the detected tool calls

```go
detectedToolCallsStr, err := helpers.ToolCallsToJSONString(detectedToolCalls)
if err != nil {
    fmt.Println("Error converting tool calls to JSON string:", err)
    return
}
fmt.Println("Detected Tool Calls:\n", detectedToolCallsStr)
```

## Execute the detected tool calls

```go
results, err := bob.ExecuteToolCalls(detectedToolCalls,
    map[string]func(any) (any, error){

        "add": func(args any) (any, error) {
            a := args.(map[string]any)["a"].(float64)
            b := args.(map[string]any)["b"].(float64)
            return a + b, nil
        },

        "say_hello": func(args any) (any, error) {
            name := args.(map[string]any)["name"].(string)
            return fmt.Sprintf("Hello, %s!", name), nil
        },

        "multiply": func(args any) (any, error) {
            a := args.(map[string]any)["a"].(float64)
            b := args.(map[string]any)["b"].(float64)
            return a * b, nil
        },
    },
)
```