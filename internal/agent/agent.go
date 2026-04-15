package agent

import (
	"context"
	"fmt"

	"charm.land/fantasy"
	"github.com/TZGyn/kode/internal/agent/tools"
	"github.com/TZGyn/kode/internal/components/message"
	"github.com/TZGyn/kode/internal/components/part"
	"github.com/TZGyn/kode/internal/layout"
)

type Agent struct {
	Agent   fantasy.Agent
	Context context.Context

	RequestRefresh bool

	Generating    bool
	HasPressedEsc bool

	messages *[]message.Message

	Layout *layout.Layout
}

func New(messages *[]message.Message, layout *layout.Layout, config Config) *Agent {
	agent := fantasy.NewAgent(
		config.model,
		fantasy.WithSystemPrompt(""),
		fantasy.WithTools(tools.GetTools()...),
	)

	return &Agent{Agent: agent, Context: config.ctx, messages: messages, Layout: layout}
}

func (a *Agent) SwitchProvider(config Config) {
	agent := fantasy.NewAgent(
		config.model,
		fantasy.WithSystemPrompt(""),
		fantasy.WithTools(tools.GetTools()...),
	)

	a.Agent = agent
	a.Context = config.ctx
}

func (a *Agent) Generate(prompt string) {
	messages := a.GetFantasyMessages()

	a.Generating = true

	defer func() {
		a.Generating = false
	}()

	streamCall := fantasy.AgentStreamCall{
		Messages: messages[:len(messages)-1],
		// The prompt.
		Prompt: prompt,

		// When we receive a chunk of streaming data.
		OnTextDelta: func(id, text string) error {
			a.AppendTextDelta(text)
			return nil
		},

		OnReasoningDelta: func(id, text string) error {
			messages := *a.messages

			part := messages[len(messages)-1]
			part.Content.AppendReasonDelta(text)
			part.UIString = part.ToUIString()
			messages[len(messages)-1] = part

			*a.messages = messages

			a.RequestRefresh = true
			return nil
		},

		// When tool calls are invoked.
		OnToolCall: func(toolCall fantasy.ToolCallContent) error {
			a.AppendToolCall(
				toolCall.ToolCallID,
				toolCall.ToolName,
				toolCall.Input,
			)
			return nil
		},

		// When a tool call completes.
		OnToolResult: func(res fantasy.ToolResultContent) error {
			a.AppendToolResult(res)

			return nil
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
		messages := *a.messages

		p := messages[len(messages)-1]
		p.Content.AppendFinish(
			part.FinishReasonError,
			"Error generating response",
			fmt.Sprintf("%v", err),
		)
		p.UIString = p.ToUIString()
		messages[len(messages)-1] = p

		*a.messages = messages
		a.RequestRefresh = true
		// fmt.Fprintf(os.Stderr, "Error generating response: %v\n", err)
		// os.Exit(1)
	}
}

func (a *Agent) GetFantasyMessages() []fantasy.Message {
	messages := []fantasy.Message{}

	for _, m := range *a.messages {
		var role fantasy.MessageRole

		if m.Role == message.Assistant {
			role = fantasy.MessageRoleAssistant
		} else if m.Role == message.User {
			role = fantasy.MessageRoleUser
		}

		parts := []fantasy.MessagePart{}

		for _, p := range *m.Content {
			switch t := p.(type) {
			case part.TextPart:
				parts = append(parts, fantasy.TextPart{
					Text: t.Content,
				})
			}
		}

		messages = append(messages, fantasy.Message{
			Role:    role,
			Content: parts,
		})
	}

	return messages
}

func (a *Agent) AppendTextDelta(delta string) {
	messages := *a.messages

	part := messages[len(messages)-1]
	part.Content.AppendTextDelta(delta)
	part.UIString = part.ToUIString()

	messages[len(messages)-1] = part

	*a.messages = messages

	a.RequestRefresh = true
}

func (a *Agent) AppendToolCall(ID string, name string, input string) {
	messages := *a.messages

	part := messages[len(messages)-1]
	part.Content.AppendToolCall(ID, name, input)
	part.UIString = part.ToUIString()

	messages[len(messages)-1] = part

	*a.messages = messages

	a.RequestRefresh = true
}

func (a *Agent) AppendToolResult(result fantasy.ToolResultContent) {
	messages := *a.messages

	part := messages[len(messages)-1]
	part.Content.AppendToolResult(result)
	part.UIString = part.ToUIString()

	messages[len(messages)-1] = part

	*a.messages = messages

	a.RequestRefresh = true
}
