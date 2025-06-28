# MCP Streamable Http tool Call

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
	agents.WithMCPStreamableHttpClient(context.Background(), "http://localhost:9090", agents.StreamableHttpOptions{}),
	agents.WithMCPStreamableHttpTools([]string{"say_hello"}),
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
results, err := bob.ExecuteMCPStreamableHTTPToolCalls(detectedToolCalls)
if err != nil {
    fmt.Println("Error executing tool calls:", err)
    return
}
fmt.Println("Results of Tool Calls:\n", results)
```