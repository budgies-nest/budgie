package agents

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/budgies-nest/budgie/enums/base"
	"github.com/openai/openai-go"
)

// go test -v -run TestHawaiianPizzaExpert
// Test the chat Completion for a Hawaiian Pizza Expert agent
func TestHawaiianPizzaExpert(t *testing.T) {

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(`
			You are a Hawaiian pizza expert. Your name is Bob.
			Provide accurate, enthusiastic information about Hawaiian pizza's history 
			(invented in Canada in 1962 by Sam Panopoulos), 
			ingredients (ham, pineapple, cheese on tomato sauce), preparation methods, and cultural impact.
			Use a friendly tone with occasional pizza puns. 
			Defend pineapple on pizza good-naturedly while respecting differing opinions. 
			If asked about other pizzas, briefly answer but return focus to Hawaiian pizza. 
			Emphasize the sweet-savory flavor combination that makes Hawaiian pizza special.
			USE ONLY THE INFORMATION PROVIDED IN THE KNOWLEDGE BASE.		
		`),
		openai.SystemMessage(`
			KNOWLEDGE BASE: 
			## Traditional Ingredients
			- Base: Traditional pizza dough
			- Sauce: Tomato-based pizza sauce
			- Cheese: Mozzarella cheese
			- Key toppings: Ham (or Canadian bacon) and pineapple
			- Optional additional toppings: Bacon, mushrooms, bell peppers, jalapeños

			## Regional Variations
			- Australia: "Hawaiian and bacon" adds extra bacon to the traditional recipe
			- Brazil: "Portuguesa com abacaxi" combines the traditional Portuguese pizza (with ham, onions, hard-boiled eggs, olives) with pineapple
			- Japan: Sometimes includes teriyaki chicken instead of ham
			- Germany: "Hawaii-Toast" is a related open-faced sandwich with ham, pineapple, and cheese
			- Sweden: "Flying Jacob" pizza includes banana, pineapple, curry powder, and chicken		
		`),
		openai.UserMessage("give me the main ingredients of the Hawaiian pizza"),
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
	response, err := bob.ChatCompletion(context.Background())

	if err != nil {
		t.Fatalf("😡 Failed to get chat completion: %v", err)
	}

	fmt.Println("🐳🤖 Chat completion result:", response)

	response = strings.ToLower(response)

	expectedWords := []string{"cheese", "bacon", "pineapple"}
	for _, word := range expectedWords {
		if !strings.Contains(response, word) {
			t.Errorf("😡 Expected response to contain word '%s', but it was not found", word)
		}
	}
}
