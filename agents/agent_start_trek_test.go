package agents

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/budgies-nest/budgie/enums/base"
	"github.com/openai/openai-go"
)

// go test -v -run TestStarTrekExpert
func TestStarTrekExpert(t *testing.T) {

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
			You are a Star Trek expert. Your name is Seven of Nine.
			USE ONLY THE INFORMATION PROVIDED IN THE KNOWLEDGE BASE.		
		`),
		openai.SystemMessage(`
			KNOWLEDGE BASE: 
			Star Trek is a science fiction media franchise that includes television series, films, books, and more.
			James T. Kirk is a fictional character in the Star Trek franchise, known for being the captain of the USS Enterprise.
			USS Enterprise is a starship in the Star Trek universe, known for its missions in space exploration.
			Spock is a fictional character in the Star Trek franchise, known for his Vulcan heritage and logical thinking.
			Leonard McCoy, also known as "Bones," is a fictional character in the Star Trek franchise, serving as the chief medical officer of the USS Enterprise.
			The best friend of James T. Kirk is Spock, who is known for his logical thinking and Vulcan heritage.
		`),
		openai.UserMessage("Who is James T. Kirk?"),
	}

	bob, err := NewAgent("Bob",
		WithDMR(base.DockerModelRunnerContainerURL),
		WithParams(openai.ChatCompletionNewParams{
			Model:       "k33g/qwen2.5:0.5b-instruct-q8_0",
			Temperature: openai.Opt(0.8),
			Messages:    messages,
		}),
	)
	if err != nil {
		t.Fatalf("😡 Failed to create agent: %v", err)
	}

	fmt.Println("🐳🤖 First Chat completion result:")

	response, err := bob.ChatCompletionStream(context.Background(), func(self *Agent, content string, err error) error {
		fmt.Print(content)
		return nil
	})

	if err != nil {
		t.Fatalf("😡 Failed to get chat completion: %v", err)
	}

	// Update the messages with the new user message
	bob.Params.Messages = append(bob.Params.Messages,
		openai.AssistantMessage(response), // NOTE: save the response as an assistant message (conversational memory)
		openai.UserMessage("Who is his best friend?"),
	)

	fmt.Println("\n\n🐳🤖 Second Chat completion result:")

	response, err = bob.ChatCompletionStream(context.Background(), func(self *Agent, content string, err error) error {
		fmt.Print(content)
		return nil
	})

	if err != nil {
		t.Fatalf("😡 Failed to get chat completion: %v", err)
	}

	response = strings.ToLower(response)

	expectedWords := []string{"spock"}
	for _, word := range expectedWords {
		if !strings.Contains(response, word) {
			t.Errorf("😡 Expected response to contain word '%s', but it was not found", word)
		}
	}
}
