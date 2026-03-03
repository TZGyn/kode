package agent

import (
	"context"
	"fmt"
	"os"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/openai"
	"github.com/TZGyn/kode/internal/components/message"
	"github.com/TZGyn/kode/internal/layout"
	"github.com/joho/godotenv"
)

type Agent struct {
	Agent   fantasy.Agent
	Context context.Context

	RequestRefresh bool

	messages *[]message.Message

	Layout *layout.Layout
}

func New(messages *[]message.Message, layout *layout.Layout) *Agent {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env found")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY environment variable")
		os.Exit(1)
	}

	provider, err := openai.New(openai.WithAPIKey(apiKey))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating OpenAI provider: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	model, err := provider.LanguageModel(ctx, "gpt-4o")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	agent := fantasy.NewAgent(
		model,
		fantasy.WithSystemPrompt(""),
	)

	return &Agent{Agent: agent, Context: ctx, messages: messages, Layout: layout}
}

func (a *Agent) Generate(prompt string) {
	// Alright, let's setup a streaming request!
	streamCall := fantasy.AgentStreamCall{
		// The prompt.
		Prompt: prompt,

		// When we receive a chunk of streaming data.
		OnTextDelta: func(id, text string) error {
			// _, fmtErr := fmt.Print(text)
			// a.RequestRefresh = true
			messages := *a.messages

			part := messages[len(messages)-1]
			part.Content.AppendTextDelta(text)
			messages[len(messages)-1] = part

			*a.messages = messages

			a.RequestRefresh = true
			return nil
		},

		// When tool calls are invoked.
		OnToolCall: func(toolCall fantasy.ToolCallContent) error {
			fmt.Printf("-> Invoking the %s tool with input %s", toolCall.ToolName, toolCall.Input)
			return nil
		},

		// When a tool call completes.
		OnToolResult: func(res fantasy.ToolResultContent) error {
			text, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentText](res.Result)
			if !ok {
				return fmt.Errorf("failed to cast result to text")
			}
			_, fmtErr := fmt.Printf("\n-> Using the %s tool: %s", res.ToolName, text.Text)
			return fmtErr
		},

		// When a step finishes, such as a tool call or a response from the
		// LLM.
		OnStepFinish: func(_ fantasy.StepResult) error {
			// fmt.Print("\n-> Step completed\n")
			return nil
		},
	}
	// Finally, let's stream everything!
	_, err := a.Agent.Stream(a.Context, streamCall)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating response: %v\n", err)
		os.Exit(1)
	}
}
