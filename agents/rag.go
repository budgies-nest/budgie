package agents

import (
	"github.com/budgies-nest/budgie/rag"
	"github.com/openai/openai-go"
)

// RAGMemorySearchSimilaritiesWithText searches for similar records in the RAG memory using the provided text.
// It creates an embedding from the text and searches for records with cosine similarity above the specified limit.
// It returns a slice of strings containing the prompts of the similar records and an error if any occurred.
// If no similar records are found, it returns an empty slice.
// It requires the DMR client to be initialized and the embedding parameters to be set in the Agent.
// The limit parameter specifies the minimum cosine similarity score for a record to be considered similar.
// It returns an error if the embedding creation fails or if the search operation fails.
func (agent *Agent) RAGMemorySearchSimilaritiesWithText(text string, limit float64) ([]string, error) {
	// Create the embedding from the question
	agent.EmbeddingParams.Input = openai.EmbeddingNewParamsInputUnion{
		OfString: openai.String(text),
	}
	embeddingResponse, err := agent.clientEngine.Embeddings.New(agent.ctx, agent.EmbeddingParams)
	if err != nil {
		return nil, err
	}
	// -------------------------------------------------
	// Create a vector record from the user embedding
	// -------------------------------------------------
	embeddingFromText := rag.VectorRecord{
		Embedding: embeddingResponse.Data[0].Embedding,
	}

	similarities, _ := agent.Store.SearchSimilarities(embeddingFromText, limit)
	var results []string
	for _, similarity := range similarities {
		results = append(results, similarity.Prompt)
	}
	return results, nil

}

// RAGMemorySearchSimilaritiesWith searches for similar records in the RAG memory using the provided embedding.
// It creates an embedding from the input and searches for records with cosine similarity above the specified limit.
// It returns a slice of strings containing the prompts of the similar records and an error if any occurred.
// If no similar records are found, it returns an empty slice.
// It requires the DMR client to be initialized and the embedding parameters to be set in the Agent.
// The limit parameter specifies the minimum cosine similarity score for a record to be considered similar.
// It returns an error if the embedding creation fails or if the search operation fails.
func (agent *Agent) RAGMemorySearchSimilaritiesWith(embedding openai.EmbeddingNewParamsInputUnion, limit float64) ([]string, error) {
	// Create the embedding from the question
	agent.EmbeddingParams.Input = embedding
	embeddingResponse, err := agent.clientEngine.Embeddings.New(agent.ctx, agent.EmbeddingParams)
	if err != nil {
		return nil, err
	}
	// -------------------------------------------------
	// Create a vector record from the user embedding
	// -------------------------------------------------
	embeddingFromText := rag.VectorRecord{
		Embedding: embeddingResponse.Data[0].Embedding,
	}

	similarities, _ := agent.Store.SearchSimilarities(embeddingFromText, limit)
	var results []string
	for _, similarity := range similarities {
		results = append(results, similarity.Prompt)
	}
	return results, nil
}
