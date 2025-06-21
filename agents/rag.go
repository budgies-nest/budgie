package agents

import (
	"encoding/json"
	"os"

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

// CreateEmbeddingFromText creates an embedding from the provided text using the OpenAI API.
// It returns the embedding and an error if any occurred.
// If the text is empty, it returns an empty embedding and no error.
func (agent *Agent) CreateEmbeddingFromText(text string) (openai.Embedding, error) {
	// Create the embedding from the text
	agent.EmbeddingParams.Input = openai.EmbeddingNewParamsInputUnion{
		OfString: openai.String(text),
	}
	embeddingResponse, err := agent.clientEngine.Embeddings.New(agent.ctx, agent.EmbeddingParams)
	if err != nil {
		return openai.Embedding{}, err
	}
	if len(embeddingResponse.Data) == 0 {
		return openai.Embedding{}, nil // No embedding created
	}
	return embeddingResponse.Data[0], nil
}

// CreateAndSaveEmbeddingFromText creates an embedding from the provided text and saves it to the memory vector store.
// It returns the saved vector record and an error if any occurred.
// If a recordId is provided, it will be used as the ID for the vector record; otherwise, a new UUID will be generated.
// The text is used as the prompt for the vector record.
func (agent *Agent) CreateAndSaveEmbeddingFromText(text string, recordId ...string) (rag.VectorRecord, error) {
	// Create the embedding from the text
	embedding, err := agent.CreateEmbeddingFromText(text)
	if err != nil {
		return rag.VectorRecord{}, err
	}

	// Create a vector record from the embedding
	vectorRecord := rag.VectorRecord{
		Id: recordId[0],
		// If no recordId is provided, it will be an empty string
		// and a new UUID will be generated in the Save method
		Prompt:    text,
		Embedding: embedding.Embedding,
	}

	// Save the vector record to the memory vector store
	savedRecord, err := agent.Store.Save(vectorRecord)
	if err != nil {
		return rag.VectorRecord{}, err
	}
	return savedRecord, nil
}

// CreateAndSaveEmbeddingFromChunks creates embeddings from the provided text chunks and saves them to the memory vector store.
// It returns a slice of saved vector records and an error if any occurred.
// Each chunk is processed individually, creating an embedding and saving it as a vector record.
// This method is useful when you have multiple text chunks and want to create and save embeddings for each of them.
// It does not require a recordId; each vector record will be saved with a new UUID generated by the Store.Save method.
// If any chunk fails to create an embedding or save, the method returns an error immediately.
func (agent *Agent) CreateAndSaveEmbeddingFromChunks(chunks []string) ([]rag.VectorRecord, error) {
	var savedRecords []rag.VectorRecord
	for _, chunk := range chunks {
		// Create the embedding from the chunk
		embedding, err := agent.CreateEmbeddingFromText(chunk)
		if err != nil {
			return nil, err
		}

		// Create a vector record from the embedding
		vectorRecord := rag.VectorRecord{
			//Id: fmt.Sprintf("chunk-%d", idx+1),
			// If no recordId is provided, it will be an empty string
			// and a new UUID will be generated in the Save method
			Prompt:    chunk,
			Embedding: embedding.Embedding,
		}

		// Save the vector record to the memory vector store
		savedRecord, err := agent.Store.Save(vectorRecord)
		if err != nil {
			return nil, err
		}
		savedRecords = append(savedRecords, savedRecord)
	}
	return savedRecords, nil
}

// SaveEmbedding saves the provided embedding to the memory vector store.
// It returns the saved vector record and an error if any occurred.
// If a recordId is provided, it will be used as the ID for the vector record; otherwise, a new UUID will be generated.
// The text is used as the prompt for the vector record.
// This method is useful when you already have an embedding and want to save it without creating a new one.
func (agent *Agent) SaveEmbedding(text string, embedding openai.Embedding, recordId ...string) (rag.VectorRecord, error) {
	// Create a vector record from the embedding
	vectorRecord := rag.VectorRecord{
		Id: recordId[0],
		// If no recordId is provided, it will be an empty string
		// and a new UUID will be generated in the Save method
		Prompt:    text,
		Embedding: embedding.Embedding,
	}

	// Save the vector record to the memory vector store
	savedRecord, err := agent.Store.Save(vectorRecord)
	if err != nil {
		return rag.VectorRecord{}, err
	}
	return savedRecord, nil
}

// PersistMemoryVectorStore persists the memory vector store to a JSON file.
// It marshals the store to JSON and writes it to the specified file path.
func (agent *Agent) PersistMemoryVectorStore() error {
	// Marshal the store to JSON
	storeJSON, err := json.MarshalIndent(agent.Store, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON to a file

	err = os.WriteFile(agent.storeFilePath, storeJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadMemoryVectorStore loads the memory vector store from a JSON file.
// It reads the file, unmarshals the JSON into a MemoryVectorStore, and assigns it to the agent's Store.
func (agent *Agent) LoadMemoryVectorStore() error {
	// Check if the store file exists
	if _, err := os.Stat(agent.storeFilePath); os.IsNotExist(err) {
		return nil // No store file to load
	}

	// Read the store file
	file, err := os.ReadFile(agent.storeFilePath)
	if err != nil {
		return err
	}

	// Unmarshal the JSON into the vector store
	var vectorStore rag.MemoryVectorStore
	if err := json.Unmarshal(file, &vectorStore); err != nil {
		return err
	}

	// Assign the loaded store to the agent
	agent.Store = &vectorStore
	return nil
}

// ResetMemoryVectorStore resets the agent's vector store to a new empty MemoryVectorStore.
// It clears the existing records and persists the empty store to the file.
// This is useful for clearing the memory and starting fresh without any records.
func (agent *Agent) ResetMemoryVectorStore() error {
	// Reset the vector store to a new empty MemoryVectorStore
	agent.Store = &rag.MemoryVectorStore{
		Records: make(map[string]rag.VectorRecord),
	}

	// Persist the empty store to the file
	return agent.PersistMemoryVectorStore()
}
