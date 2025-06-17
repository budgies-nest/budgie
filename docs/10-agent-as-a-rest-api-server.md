# Agent as REST API Server

## Initialize the agent
```golang
bob, err := agents.NewAgent("Bob",
    agents.WithDMR(context.Background(), base.DockerModelRunnerContainerURL),
    agents.WithParams(openai.ChatCompletionNewParams{
        Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
        Temperature: openai.Opt(0.8),
        Messages: []openai.ChatCompletionMessageParamUnion{
            openai.SystemMessage("You're a helpful assistant expert with Star Trek universe."),
        },
    }),
    agents.WithHTTPServer(agents.HTTPServerConfig{
        Port: "8080",
    }),
)
```

## Start the REST API server

```golang
// Start the HTTP server
fmt.Println("Starting HTTP server on port 8080...")
err = bob.StartHttpServer()
if err != nil {
    panic(err)
}
```

## Query the agent API

```bash
curl http://localhost:8080/api/chat \
-H "Content-Type: application/json" \
-d '
{
  "user": "who is James T Kirk?"
}' 
```


```bash
curl --no-buffer http://localhost:8080/api/chat-stream \
-H "Content-Type: application/json" \
-d '
{
  "user": "who is Jean-Luc Picard?"
}' 
```