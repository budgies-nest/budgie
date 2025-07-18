package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func main() {
	modelRunnerBaseUrl := helpers.GetModelRunnerBaseUrl()

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: "ai/mxbai-embed-large",
			},
		),
		agents.WithMemoryVectorStore("bob.json"),
	)
	if err != nil {
		panic(err)
	}
	bob.LoadMemoryVectorStore()

	similarities, err := bob.RAGMemorySearchSimilaritiesWithText(context.Background(), "Who is Emma Peel?", 0.6)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Similarities found:")
	for _, similarity := range similarities {
		fmt.Println("-", similarity)
	}

	similarities, err = bob.RAGMemorySearchSimilaritiesWithText(context.Background(), "Who is John Steed?", 0.6)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Similarities found:")
	for _, similarity := range similarities {
		fmt.Println("-", similarity)
	}

}
