package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

var chunks = []string{
	`# The Avengers
	"The Avengers" is a classic British spy-fi television series that aired from 1961 to 1969. 
	The show exemplifies the unique style of 1960s British television with its blend of espionage,
	 science fiction, and quintessentially British humor. 
	The series follows secret agents working for a specialized branch of British intelligence, 
	battling eccentric villains and foiling bizarre plots to undermine national security.`,

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

	`# Mother
    Mother, portrayed by Patrick Newell, is Steed's wheelchair-bound superior who appears in later seasons. 
	Operating from various unusual locations, this eccentric spymaster directs operations with a mix of authority 
	and peculiarity that fits perfectly within the show's offbeat universe.`,
}

func main() {

	modelRunnerBaseUrl := helpers.GetModelRunnerBaseUrl()

	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: "ai/mxbai-embed-large",
			},
		),
		agents.WithRAGMemory(context.Background(), chunks),
	)

	if err != nil {
		panic(err)
	}
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
