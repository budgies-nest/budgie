# Rag (in memory) Agent with persistable Memory Vector Store

## Initialize the agent
```golang
bob, err := agents.NewAgent("Bob",
    agents.WithDMR(base.DockerModelRunnerContainerURL),
    agents.WithEmbeddingParams(
        openai.EmbeddingNewParams{
            Model: "ai/mxbai-embed-large",
        },
    ),
    agents.WithMemoryVectorStore("bob.json"),
)
```

`agents.WithMemoryVectorStore("bob.json")` allows you to create an agent that can use a persistable memory vector store for RAG (Retrieval-Augmented Generation) tasks, into a JSON file.


## Generate the embeddings and store them in the vector store

```golang
var chunks = []string{...}

for idx, chunk := range chunks {
    _, err = bob.CreateAndSaveEmbeddingFromText(chunk, fmt.Sprintf("chunk-%d", idx+1))
    if condition := err != nil; condition {
        fmt.Println("😡 Error creating embedding:", err)
        return
    }
}

bob.PersistMemoryVectorStore()
```

## Load the vector store from the file and search for similarity

```golang
bob.LoadMemoryVectorStore()
similarities, err := bob.RAGMemorySearchSimilaritiesWithText("Who is Emma Peel?", 0.6)

if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Similarities found:")
for _, similarity := range similarities {
    fmt.Println("-", similarity)
}
```

## Reset the vector store

```golang
bob.ResetMemoryVectorStore()
```