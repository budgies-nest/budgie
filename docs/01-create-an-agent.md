# Create an agent

```golang
bob, err := agents.NewAgent("Bob",
    agents.WithDMR(context.Background(), base.DockerModelRunnerContainerURL),
    agents.WithParams(
        openai.ChatCompletionNewParams{
            Model: "k33g/qwen2.5:0.5b-instruct-q8_0",
            Temperature: openai.Opt(0.8),
            Messages: []openai.ChatCompletionMessageParamUnion{
                openai.SystemMessage("You're a helpful assistant expert with Star Trek universe."),
                openai.UserMessage("Who is James T Kirk?"),
            },
        }
    ),
)
```
> Requirements: `docker model pull k33g/qwen2.5:0.5b-instruct-q8_0`

## Base URL enumeration

```golang
DockerModelRunnerLocalURL = "http://localhost:12434/engines/llama.cpp/v1"
DockerModelRunnerContainerURL = "http://model-runner.docker.internal/engines/llama.cpp/v1"

// subject to change
DockerModelRunnerDockerCEURL = "http://172.17.0.1:12434/engines/llama.cpp/v1"
DockerModelRunnerDockerCloudURL = "http://172.17.0.1:12435/engines/llama.cpp/v1"

OpenAIURL = "https://api.openai.com/v1"
```