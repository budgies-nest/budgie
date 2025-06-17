package main

import (
	"context"
	"fmt"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/enums/base"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
)

var chunks = []string{
	`Michael Burnham is the main character on the Star Trek series, Discovery.  
		She's a human raised on the logical planet Vulcan by Spock's father.  
		Burnham is intelligent and struggles to balance her human emotions with Vulcan logic.  
		She's become a Starfleet captain known for her determination and problem-solving skills.
		Originally played by actress Sonequa Martin-Green`,

	`James T. Kirk, also known as Captain Kirk, is a fictional character from the Star Trek franchise.  
		He's the iconic captain of the starship USS Enterprise, 
		boldly exploring the galaxy with his crew.  
		Originally played by actor William Shatner, 
		Kirk has appeared in TV series, movies, and other media.`,

	`Jean-Luc Picard is a fictional character in the Star Trek franchise.
		He's most famous for being the captain of the USS Enterprise-D,
		a starship exploring the galaxy in the 24th century.
		Picard is known for his diplomacy, intelligence, and strong moral compass.
		He's been portrayed by actor Patrick Stewart.`,
}

func main() {
	bob, err := agents.NewAgent("Bob",
		agents.WithDMR(context.Background(), base.DockerModelRunnerContainerURL),
		agents.WithEmbeddingParams(
			openai.EmbeddingNewParams{
				Model: "ai/mxbai-embed-large",
			},
		),
		agents.WithRAGMemory(chunks),
		agents.WithMCPStreamableHttpServer(agents.MCPServerConfig{
			Port:     "9090",
			Version:  "v1",
			Name:     "mcp-bob",
			Endpoint: "/mcp",
		}),
	)

	if err != nil {
		panic(err)
	}

	searchInDoc := mcp.NewTool("question_about_something",
		mcp.WithDescription(`Find an answer in the internal database.`),
		mcp.WithString("question",
			mcp.Required(),
			mcp.Description("Search question"),
		),
	)

	bob.AddToolToMCPServer(searchInDoc, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		question := args["question"].(string)
		similarities, err := bob.RAGMemorySearchSimilaritiesWithText(question, 0.7)
		if err != nil {
			return nil, err
		}
		content := ""
		for _, similarity := range similarities {
			content += similarity + "\n"
		}
		return mcp.NewToolResultText(content), nil
	})

	fmt.Println("MCP Streamable HTTP server Agent is running on port 9090")
	bob.StartMCPHttpServer()

}
