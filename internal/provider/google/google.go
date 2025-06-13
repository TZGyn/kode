package google

import (
	"context"
	"time"

	"github.com/TZGyn/kode/internal/tool"
	"google.golang.org/genai"
)

type Config struct {
	GEMINI_API_KEY string `json:"GEMINI_API_KEY"`
	Model          string `json:"model"`
}

func DefaultConfig(apiKey string, model string) Config {
	return Config{
		GEMINI_API_KEY: apiKey,
		Model:          model,
	}
}

type GoogleClient struct {
	context       context.Context
	cancelRequest context.CancelFunc

	model string

	client        *genai.Client
	Messages      []*genai.Content
	FunctionCalls []string
}

func CreateGoogle(config Config) (*GoogleClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.GEMINI_API_KEY,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		cancel()
		return nil, err
	}

	return &GoogleClient{
		context:       ctx,
		cancelRequest: cancel,

		model: config.Model,

		client:        client,
		Messages:      []*genai.Content{},
		FunctionCalls: []string{},
	}, nil
}

func (c *GoogleClient) SendMessage(messages []*genai.Content, response *string) error {
	content, err := c.client.Models.GenerateContent(
		c.context,
		c.model,
		messages,
		googleConfig,
	)
	if err != nil {
		return err
	}

	text := ""

	for _, part := range content.Candidates[0].Content.Parts {
		if part.Text != "" {
			if part.Thought {
				continue
			}
			*response += part.Text
			text += part.Text
		}
	}

	c.Messages = append(c.Messages, &genai.Content{Role: "assistant", Parts: []*genai.Part{{Text: text}}})

	if len(content.FunctionCalls()) == 0 {
		return nil
	}

	for _, functionCall := range content.FunctionCalls() {
		result, err := tool.HandleTool(functionCall.Name, functionCall.Args, response)
		if err != nil {
			continue
		}

		messages = append(messages, &genai.Content{
			Role: "tool",
			Parts: []*genai.Part{
				{
					FunctionResponse: &genai.FunctionResponse{
						ID:   functionCall.ID,
						Name: functionCall.Name,
						Response: map[string]any{
							"result": result,
						},
					},
				},
			}})
	}

	return c.SendMessage(c.Messages, response)
}

func (c *GoogleClient) CancelRequest() error {
	if c.cancelRequest != nil {
		c.cancelRequest()
	}
	return nil
}
