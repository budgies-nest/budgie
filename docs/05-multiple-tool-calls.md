# Multiple Tool Calls
> Experimental feature: Multiple tool calls in a single prompt.

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
    agents.WithDMR(base.DockerModelRunnerContainerURL),
    agents.WithParams(
        openai.ChatCompletionNewParams{
			Model: "k33g/llama-xlam-2:8b-fc-r-q2_k",
            Temperature: openai.Opt(0.0),
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.UserMessage(`
                    Add 10 and 32
                    Add 12 and 30
                    Add 40 and 2
                    Add 5 and 37			
                `),
            },
            ParallelToolCalls: openai.Bool(true),
        },
    ),
    agents.WithTools([]openai.ChatCompletionToolParam{addTool}),
)
```

> ✋ Few models are able to call/detect multiple tools in a single prompt. Therefore, it is important to test your model to verify its effectiveness with multiple tool calls.
> - Always set `Temperature` to `0.0`.
> - Set `ParallelToolCalls` to **`true`**.
>
> `k33g/llama-xlam-2:8b-fc-r-q2_k` works well with multiple tool calls, but is not able to dectect the tool calls with `ParallelToolCalls` set to **`false`**.

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
 4
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
        "id": "I1-1"
    },
    {
        "function": {
            "arguments": {
                "a": 12,
                "b": 30
            },
            "name": "add"
        },
        "id": "I1-2"
    },
    {
        "function": {
            "arguments": {
                "a": 40,
                "b": 2
            },
            "name": "add"
        },
        "id": "I1-3"
    },
    {
        "function": {
            "arguments": {
                "a": 5,
                "b": 37
            },
            "name": "add"
        },
        "id": "I1-4"
    }
]
```

## Execute the tool calls
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
 [42 42 42 42]
```
