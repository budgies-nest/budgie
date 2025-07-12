package agents

// TODO: add args to save to file
func WithMemoryVectorStore(storeFilePath string) AgentOption {
	return func(agent *Agent) {

		agent.storeFilePath = storeFilePath

	}
}
