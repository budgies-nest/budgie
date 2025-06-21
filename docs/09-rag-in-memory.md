# Rag (in memory) Agent
> `agents.WithRAGMemory(chunks)` allows you to create an agent that can use a set of chunks in memory for RAG (Retrieval-Augmented Generation) tasks. The embeddings are created on the fly when the agent is initialized, and the chunks are stored in memory. This is useful for small datasets or when you want to avoid the overhead of a database.

## Create some chunks

```golang

var chunks = []string{
	`# John Steed
    John Steed, portrayed by Patrick Macnee, is the quintessential English gentleman spy 
	who never leaves home without his trademark bowler hat and umbrella (which conceals various weapons). 
	Charming, witty, and deceptively dangerous, Steed approaches even the most perilous situations 
	with impeccable manners and a dry sense of humor. 
	His refined demeanor masks his exceptional combat skills and razor-sharp intelligence.`,

	`# Emma Peel
     Emma Peel, played by Diana Rigg, is perhaps the most iconic of Steed's partners. 
	 A brilliant scientist, martial arts expert, and fashion icon, Mrs. Peel combines beauty, brains, 
	 and remarkable fighting skills. Clad in her signature leather catsuits, she represents the modern, 
	 liberated woman of the 1960s. Her name is a play on "M-appeal" (man appeal), 
	 but her character transcended this origin to become a feminist icon.`,

	`# Tara King
     Tara King, played by Linda Thorson, was Steed's final regular partner in the original series. 
	 Younger and somewhat less experienced than her predecessors, King was nevertheless a trained agent 
	 who continued the tradition of strong female characters. 
	 Her relationship with Steed had more romantic undertones than previous partnerships, 
	 and she brought a fresh, youthful energy to the series.`,
}
```

## Initialize the agent

```golang
bob, err := agents.NewAgent("Bob",
	agents.WithDMR(context.Background(), base.DockerModelRunnerContainerURL),
	agents.WithEmbeddingParams(
		openai.EmbeddingNewParams{
			Model: "ai/mxbai-embed-large",
		},
	),
	agents.WithRAGMemory(chunks),
)
```

## Search for similarity

```golang
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
