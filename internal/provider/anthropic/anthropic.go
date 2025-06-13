package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/TZGyn/kode/internal/provider/prompt"
	"github.com/TZGyn/kode/internal/tool"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Config struct {
	ANTHROPIC_API_KEY string `json:"ANTHROPIC_API_KEY"`
	Model             string `json:"model"`
}

type AnthropicClient struct {
	context       context.Context
	cancelRequest context.CancelFunc

	client anthropic.Client
	model  anthropic.Model

	Messages []anthropic.MessageParam
}

func DefaultConfig(apiKey string, model string) Config {
	return Config{
		ANTHROPIC_API_KEY: apiKey,
		Model:             model,
	}
}

func Create(config Config) (*AnthropicClient, error) {
	model := anthropic.Model(config.Model)

	if model == "" {
		return nil, errors.New("invalid model")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	client := anthropic.NewClient(
		option.WithAPIKey(config.ANTHROPIC_API_KEY),
	)

	return &AnthropicClient{
		context:       ctx,
		cancelRequest: cancel,

		model: model,

		client: client,
	}, nil
}

func (c *AnthropicClient) SendMessage(messages []anthropic.MessageParam, response *string) error {
	c.Messages = messages

	completion, err := c.client.Messages.New(
		c.context,
		anthropic.MessageNewParams{
			Messages: messages,
			System: []anthropic.TextBlockParam{
				{
					Text: prompt.SystemPrompt(),
				},
			},
			Tools:     tools,
			Model:     c.model,
			MaxTokens: 5000,
		},
	)

	if err != nil {
		return err
	}

	toolResults := []anthropic.ContentBlockParamUnion{}
	for _, block := range completion.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.TextBlock:
			*response += block.Text
		case anthropic.ToolUseBlock:
			var input map[string]any
			err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input)
			if err != nil {
				continue
			}

			result, err := tool.HandleTool(block.Name, input, response)
			if err != nil {
				continue
			}
			toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, result, false))
		}
	}

	if len(toolResults) == 0 {
		return nil
	}

	messages = append(messages, completion.ToParam())
	messages = append(messages, anthropic.NewUserMessage(toolResults...))

	return c.SendMessage(messages, response)
}

func (c *AnthropicClient) CancelRequest() error {
	if c.cancelRequest != nil {
		c.cancelRequest()
	}
	return nil
}
