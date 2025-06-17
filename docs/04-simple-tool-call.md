# Simple Tool Call

## Define a tool
```golang
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
```

## Initialize the agent

```golang
bob, err := agents.NewAgent("Bob",
    agents.WithDMR(context.Background(), base.DockerModelRunnerContainerURL),
    agents.WithParams(
        openai.ChatCompletionNewParams{
            Model: "k33g/qwen2.5:0.5b-instruct-q8_0",
            Temperature: openai.Opt(0.0),
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.UserMessage(`
                    Add 10 and 32			
                `),
            },
            ParallelToolCalls: openai.Bool(false),
        },
    ),
    agents.WithTools([]openai.ChatCompletionToolParam{addTool}),
)
```

> ✋ Small models are not very good at function calling, then plan to make one call at a time and everything will be fine.
> - Always set `Temperature` to `0.0`.
> - Always set `ParallelToolCalls` to **`false`** for small models.

## Run the tool completion
```golang
// Generate the tools detection completion
detectedToolCalls, err := bob.ToolsCompletion()
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Number of Tool Calls:\n", len(detectedToolCalls))

detectedToolCallsStr, err := helpers.ToolCallsToJSONString(detectedToolCalls)
if err != nil {
    fmt.Println("Error converting tool calls to JSON string:", err)
    return
}
fmt.Println("Detected Tool Calls:\n", detectedToolCallsStr)
```

> Output:
```json
Number of Tool Calls:
 1
Detected Tool Calls:
 [
    {
        "function": {
            "arguments": {
                "a": 10,
                "b": 32
            },
            "name": "add"
        },
        "id": "mHvFxbwlvvIqwKLxm51oesCWR9fSzqwu"
    }
]
```

## Execute the tool call
```golang
results, err := bob.ExecuteToolCalls(detectedToolCalls,
    map[string]func(any) (any, error){

        "add": func(args any) (any, error) {
            a := args.(map[string]any)["a"].(float64)
            b := args.(map[string]any)["b"].(float64)
            return a + b, nil
        },
    },
)
if err != nil {
    fmt.Println("Error executing tool calls:", err)
    return
}
fmt.Println("Results of Tool Calls:\n", results)

```

> Output:
```text
Results of Tool Calls:
 [42]
```
