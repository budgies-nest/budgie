package agents

import (
	"fmt"

	"github.com/budgies-nest/budgie/rag"
	"github.com/openai/openai-go"
)

// WithRAGMemory initializes the Agent with a RAG memory using the provided chunks.
// It creates a MemoryVectorStore and saves the embeddings of the chunks into it.
// The chunks should be pre-processed text data that will be used for retrieval-augmented generation (RAG).
// It returns an AgentOption that can be used to configure the agent.
func WithRAGMemory(chunks []string) AgentOption {
	return func(agent *Agent) {
		// -------------------------------------------------
		// Create a vector store
		// -------------------------------------------------
		agent.Store = rag.MemoryVectorStore{
			Records: make(map[string]rag.VectorRecord),
		}

		// -------------------------------------------------
		// Create and save the embeddings from the chunks
		// -------------------------------------------------
		for _, chunk := range chunks {

			agent.EmbeddingParams.Input = openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(chunk),
			}
			embeddingsResponse, err := agent.clientEngine.Embeddings.New(agent.ctx, agent.EmbeddingParams)

			if err != nil {
				agent.optionError = fmt.Errorf("failed to create embedding for chunk: %w", err)
				return
			} else {
				_, errSave := agent.Store.Save(rag.VectorRecord{
					Prompt:    chunk,
					Embedding: embeddingsResponse.Data[0].Embedding,
				})
				if errSave != nil {
					agent.optionError = errSave
					return
					// QUESTION: How to handle the error?
					// TODO: do some samples to define what to do
					// IMPORTANT! ...
				}
			}
		}
	}
}
