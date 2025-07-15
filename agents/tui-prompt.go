package agents

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/charmbracelet/huh"
)

// PromptConfig defines configuration options for the TUI prompt
type PromptConfig struct {
	UseStreamCompletion        bool
	StartingMessage            string // Optional starting message (default: "🤖 Starting TUI for agent: {name}")
	ExplanationMessage         string // Optional explanation message (default: "Type your questions below. Use '/bye' to quit or Ctrl+C to interrupt completions.")
	PromptTitle                string // Optional prompt title (default: "💬 Chat with {name}")
	ThinkingPrompt             string // Optional thinking prompt (default: "🤔 ")
	InterruptInstructions      string // Optional interrupt instructions (default: "(Press Ctrl+C to interrupt)")
	CompletionInterruptMessage string // Optional completion interrupt message (default: "🚫 Completion was interrupted\n")
	GoodbyeMessage             string // Optional goodbye message (default: "👋 Goodbye!")
}

// Prompt starts an interactive TUI that allows users to chat with the agent
// config allows configuration of the prompt behavior (streaming vs non-streaming)
// Returns an error if the TUI cannot be started
func (agent *Agent) Prompt(config PromptConfig) error {
	// Set default messages if not provided
	startingMessage := config.StartingMessage
	if startingMessage == "" {
		startingMessage = fmt.Sprintf("🤖 Starting TUI for agent: %s", agent.Name)
	}
	
	explanationMessage := config.ExplanationMessage
	if explanationMessage == "" {
		explanationMessage = "Type your questions below. Use '/bye' to quit or Ctrl+C to interrupt completions."
	}
	
	promptTitle := config.PromptTitle
	if promptTitle == "" {
		promptTitle = fmt.Sprintf("💬 Chat with %s", agent.Name)
	}
	
	thinkingPrompt := config.ThinkingPrompt
	if thinkingPrompt == "" {
		thinkingPrompt = "🤔 "
	}
	
	interruptInstructions := config.InterruptInstructions
	if interruptInstructions == "" {
		interruptInstructions = "(Press Ctrl+C to interrupt)"
	}
	
	completionInterruptMessage := config.CompletionInterruptMessage
	if completionInterruptMessage == "" {
		completionInterruptMessage = "🚫 Completion was interrupted\n"
	}
	
	goodbyeMessage := config.GoodbyeMessage
	if goodbyeMessage == "" {
		goodbyeMessage = "👋 Goodbye!"
	}

	fmt.Println(startingMessage)
	fmt.Println(explanationMessage)
	fmt.Println()

	for {
		var userInput string

		// Create the input form
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title(promptTitle).
					Placeholder("Type your question here...").
					Value(&userInput).
					ExternalEditor(false),
			),
		)

		// Run the form
		err := form.Run()
		if err != nil {
			return fmt.Errorf("form error: %w", err)
		}

		// Trim whitespace
		userInput = strings.TrimSpace(userInput)

		// Check for empty input
		if userInput == "" {
			continue
		}

		// Check for /bye command
		if userInput == "/bye" {
			fmt.Println(goodbyeMessage)
			break
		}

		// Add user message to the agent's parameters
		agent.AddUserMessage(userInput)

		// Create a context that can be cancelled with Ctrl+C or ESC
		ctx, cancel := context.WithCancel(context.Background())
		
		// Set up signal handling for interrupting completion
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		
		go func() {
			<-sigChan
			fmt.Print(completionInterruptMessage)
			cancel()
		}()

		if config.UseStreamCompletion {
			// Handle streaming completion
			fmt.Print(thinkingPrompt)
			fmt.Println(interruptInstructions)
			
			response, err := agent.ChatCompletionStream(ctx, func(self *Agent, content string, err error) error {
				if err != nil {
					return err
				}
				fmt.Print(content)
				return nil
			})
			
			fmt.Println() // New line after completion

			if err != nil {
				if ctx.Err() == context.Canceled {
					fmt.Print(completionInterruptMessage)
				} else {
					fmt.Printf("❌ Error: %v\n", err)
				}
			} else {
				// Add assistant response to conversation
				agent.AddAssistantMessage(response)
			}
		} else {
			// Handle regular completion
			fmt.Print(thinkingPrompt + "Thinking...")
			fmt.Println(" " + interruptInstructions)
			
			response, err := agent.ChatCompletion(ctx)
			if err != nil {
				if ctx.Err() == context.Canceled {
					fmt.Print(completionInterruptMessage)
				} else {
					fmt.Printf("❌ Error: %v\n", err)
				}
			} else {
				fmt.Printf("🤖 %s\n", response)
				// Add assistant response to conversation
				agent.AddAssistantMessage(response)
			}
		}

		// Clean up signal handling
		signal.Stop(sigChan)
		cancel()

		fmt.Println() // Add spacing between interactions
	}

	return nil
}

