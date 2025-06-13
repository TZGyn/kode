package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
					Text: fmt.Sprintf(`
						You are a cli code assistant named kode
						Today's Date: %s

						It is a must to generate some text, letting the user knows your thinking process before using a tool.
						Thus providing better user experience, rather than immediately jump to using the tool and generate a conclusion

						Common Order: Tool, Text
						Better order you must follow: Text, Tool, Text

						You have been given tools to fulfill user request, make sure to keep using them until the user request is fulfilled
						Always check the progress to make sure you dont infinite loop
					`, time.Now().Format("2006-01-02 15:04:05")),
				},
			},
			Tools: tools,
			Model: c.model,
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
