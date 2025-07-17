package base

const (
	DockerModelRunnerLocalURL = "http://localhost:12434/engines/llama.cpp/v1"
	DockerModelRunnerContainerURL = "http://model-runner.docker.internal/engines/llama.cpp/v1"
	// NOTE: subject to change
	DockerModelRunnerDockerCEURL = "http://172.17.0.1:12434/engines/llama.cpp/v1"
	DockerModelRunnerDockerCloudURL = "http://172.17.0.1:12435/engines/llama.cpp/v1"

	OpenAIURL = "https://api.openai.com/v1"

	OllamaLocalURL = "http://localhost:11434/v1"
	OllamaContainerURL = "http://host.docker.internal:11434/v1"
)