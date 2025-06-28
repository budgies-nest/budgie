# Agent as a MCP Server

## Initialize the agent

```golang
bob, err := agents.NewAgent("Bob",
    agents.WithDMR(base.DockerModelRunnerContainerURL),
    agents.WithEmbeddingParams(
        openai.EmbeddingNewParams{
            Model: "ai/mxbai-embed-large",
        },
    ),
    agents.WithRAGMemory(chunks),
    agents.WithMCPStreamableHttpServer(agents.MCPServerConfig{
        Port:     "9090",
        Version:  "v1",
        Name:     "mcp-bob",
        Endpoint: "/mcp",
    }),
)
```

## Add a tool to the agent

```golang
searchInDoc := mcp.NewTool("question_about_something",
    mcp.WithDescription(`Find an answer in the internal database.`),
    mcp.WithString("question",
        mcp.Required(),
        mcp.Description("Search question"),
    ),
)

bob.AddToolToMCPServer(searchInDoc, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := request.GetArguments()
    question := args["question"].(string)
    similarities, err := bob.RAGMemorySearchSimilaritiesWithText(question, 0.7)
    if err != nil {
        return nil, err
    }
    content := ""
    for _, similarity := range similarities {
        content += similarity + "\n"
    }
    return mcp.NewToolResultText(content), nil
})
```

## Start the MCP server

```golang
fmt.Println("MCP Streamable HTTP server Agent is running on port 9090")
bob.StartMCPHttpServer()
```
