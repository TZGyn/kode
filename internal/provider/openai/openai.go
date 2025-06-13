package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TZGyn/kode/internal/tool"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Config struct {
	OPENAI_API_KEY string `json:"OPENAI_API_KEY"`
	Model          string `json:"model"`
}

type OpenAIClient struct {
	context       context.Context
	cancelRequest context.CancelFunc

	client openai.Client
	model  string

	Messages []openai.ChatCompletionMessageParamUnion
}

func DefaultConfig(apiKey string, model string) Config {
	return Config{
		OPENAI_API_KEY: apiKey,
		Model:          model,
	}
}

func Create(config Config) (*OpenAIClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	client := openai.NewClient(
		option.WithAPIKey(config.OPENAI_API_KEY),
	)

	return &OpenAIClient{
		context:       ctx,
		cancelRequest: cancel,

		model: config.Model,

		client: client,
	}, nil
}

func (c *OpenAIClient) SendMessage(messages []openai.ChatCompletionMessageParamUnion, response *string) error {
	c.Messages = messages
	withSystemMessage := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(
			fmt.Sprintf(`
				You are a cli code assistant named kode
				Today's Date: %s

				It is a must to generate some text, letting the user knows your thinking process before using a tool.
				Thus providing better user experience, rather than immediately jump to using the tool and generate a conclusion

				Common Order: Tool, Text
				Better order you must follow: Text, Tool, Text

				You have been given tools to fulfill user request, they are optional to use but make sure to use them if needed to fulfill the user request
				Always check the progress to make sure you dont infinite loop
			`, time.Now().Format("2006-01-02 15:04:05")),
		),
	}
	withSystemMessage = append(withSystemMessage, messages...)
	params := openai.ChatCompletionNewParams{
		Messages: withSystemMessage,
		Tools:    tools,
		Model:    c.model,
	}

	completion, err := c.client.Chat.Completions.New(c.context, params)
	if err != nil {
		return err
	}
	*response += completion.Choices[0].Message.Content + "\n"

	toolCalls := completion.Choices[0].Message.ToolCalls
	if len(toolCalls) == 0 {
		return nil
	}

	params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
	for _, toolCall := range toolCalls {
		var args map[string]any
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			panic(err)
		}
		result, err := tool.HandleTool(toolCall.Function.Name, args, response)
		if err != nil {
			continue
		}
		params.Messages = append(params.Messages, openai.ToolMessage(result, toolCall.ID))
	}

	return c.SendMessage(params.Messages, response)
}

func (c *OpenAIClient) CancelRequest() error {
	if c.cancelRequest != nil {
		c.cancelRequest()
	}
	return nil
}
