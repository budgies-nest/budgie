# MCP Stdio tool Call

## Initialize the agent

```golang
bob, err := agents.NewAgent("Bob",
    agents.WithDMR(base.DockerModelRunnerContainerURL),
    agents.WithParams(
        openai.ChatCompletionNewParams{
            Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
            Temperature: openai.Opt(0.0),
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.UserMessage(`
                    Say Hello to Bob	
                `),
            },
            ParallelToolCalls: openai.Bool(false),
        },
    ),
    // Define the MCP Stdio Server
    agents.WithMCPStdioClient(
        context.Background(),
        "go",
        agents.STDIOCommandOptions{
            "run",
            "./mcp-stdio-server/main.go",
        },
        agents.EnvVars{},
    ),
    // Define the MCP Stdio Tools and filter the tool(s) to be used
    agents.WithMCPStdioTools([]string{"say_hello"}),
)
```

## Run the tool completion

```golang
// Generate the tools detection completion
detectedToolCalls, err := bob.ToolsCompletion()
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Number of Tool Calls:\n", len(detectedToolCalls))
```

## Execute the tool call

```golang
results, err := bob.ExecuteMCPStdioToolCalls(detectedToolCalls)
if err != nil {
    fmt.Println("Error executing tool calls:", err)
    return
}
fmt.Println("Results of Tool Calls:\n", results)
```