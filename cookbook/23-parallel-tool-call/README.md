# TUI Agent Example


## Start the application

**From a container**:
```bash
MODEL_RUNNER_BASE_URL=http://model-runner.docker.internal/engines/llama.cpp/v1 \
MODEL_RUNNER_CHAT_MODEL=ai/qwen2.5:latest \
MODEL_RUNNER_TOOLS_MODEL=hf.co/salesforce/xlam-2-3b-fc-r-gguf:q3_k_l \
go run main.go
```


**From a local machine**:
```bash
MODEL_RUNNER_BASE_URL=http://localhost:12434/engines/llama.cpp/v1 \
MODEL_RUNNER_CHAT_MODEL=ai/qwen2.5:latest \
MODEL_RUNNER_TOOLS_MODEL=hf.co/salesforce/xlam-2-3b-fc-r-gguf:q3_k_l \
go run main.go
```

## Talk to the Werewolf

- What is your name?
- What is your occupation?
- What is your favorite food?  
- What is your background story?
- What is your main quote?

## Test the tools

- what is your health value?
- set your health value to 200
- increase your health by 10
- decrease your health by 5
- what is your intelligence value?
- set your intelligence value to 100
- decrease your intelligence by 5