package google

import (
	"context"
	"fmt"
	"time"

	"github.com/TZGyn/kode/internal/tool"
	"google.golang.org/genai"
)

type Config struct {
	GEMINI_API_KEY string `json:"GEMINI_API_KEY"`
}

func DefaultConfig(apiKey string) Config {
	return Config{GEMINI_API_KEY: apiKey}
}

type GoogleClient struct {
	context       context.Context
	cancelRequest context.CancelFunc

	client   *genai.Client
	chat     *genai.Chat
	Messages []*genai.Content
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

	chat, err := client.Chats.Create(
		ctx,
		"gemini-2.0-flash",
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Role: "system",
				Parts: []*genai.Part{
					{
						Text: fmt.Sprintf(`
							You are a chat assistant
							Today's Date: %s

							It is a must to generate some text, letting the user knows your thinking process before using a tool.
							Thus providing better user experience, rather than immediately jump to using the tool and generate a conclusion

							Common Order: Tool, Text
							Better order you must follow: Text, Tool, Text

							You have been given a tool which will take a directory as input and return its direct children
							You may call this tool repeatedly until you fulfill the user request
						`, time.Now().Format("2006-01-02 15:04:05")),
					},
				},
			},
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						{
							Name:        "list-directory",
							Description: "Given a directory, return all the children of it",
							Parameters: &genai.Schema{
								Type: "object",
								Properties: map[string]*genai.Schema{
									"directory": {
										Type:        "string",
										Description: "the directory to output, use . for root",
									},
								},
							},
							Response: &genai.Schema{
								Type: "object",
								Properties: map[string]*genai.Schema{
									"children": {
										Type:        "array",
										Description: "children as list",
										Items: &genai.Schema{
											Type:        "string",
											Description: "children, folder or file",
										},
									},
								},
							},
						},
						{
							Name:        "cat-file",
							Description: "Given a file path, return all its content as string",
							Parameters: &genai.Schema{
								Type: "object",
								Properties: map[string]*genai.Schema{
									"filePath": {
										Type:        "string",
										Description: "the file path to output relative to root",
									},
								},
							},
							Response: &genai.Schema{
								Type: "object",
								Properties: map[string]*genai.Schema{
									"content": {
										Type:        "string",
										Description: "file content",
									},
								},
							},
						},
					},
				},
			},
		},
		nil,
	)

	if err != nil {
		cancel()
		return nil, err
	}

	return &GoogleClient{
		context:       ctx,
		cancelRequest: cancel,

		client:   client,
		chat:     chat,
		Messages: []*genai.Content{},
	}, nil
}

func (c *GoogleClient) SendMessage(messages []*genai.Content, response *string) {
	var parts []genai.Part

	for _, content := range messages {
		for _, part := range content.Parts {
			parts = append(parts, *part)
		}
	}

	stream := c.chat.SendMessageStream(
		c.context,
		parts...,
	)
	text := ""
	for result, err := range stream {
		if err != nil {
			return
		}

		for _, part := range result.Candidates[0].Content.Parts {
			if part.Text != "" {
				if part.Thought {
					continue
				}
				*response = *response + part.Text
				text += part.Text
			}
			if part.FunctionCall != nil {
				c.Messages = append(c.Messages, &genai.Content{Role: "assistant", Parts: []*genai.Part{{Text: text}}})

				text = ""

				functionCall := part.FunctionCall
				if functionCall.Name == "list-directory" {
					result := []string{}

					directory, ok := functionCall.Args["directory"].(string)
					if ok {
						entires, err := tool.ListDirectory(directory)
						if err == nil {
							result = entires
						}

						toolResult := ""
						toolResult += "## Files Start\n"
						for _, entry := range result {
							toolResult += "- " + entry + "\n"
						}
						toolResult += "## Files End\n"

						*response = *response + toolResult
					}

					messages = append(messages, &genai.Content{
						Role: "tool",
						Parts: []*genai.Part{
							{
								FunctionResponse: &genai.FunctionResponse{
									ID:   functionCall.ID,
									Name: functionCall.Name,
									Response: map[string]any{
										"children": result,
									},
								},
							},
						}})

					c.SendMessage(
						messages,
						response,
					)
				}
				if functionCall.Name == "cat-file" {
					result := ""

					filePath, ok := functionCall.Args["filePath"].(string)
					if ok {
						content, err := tool.CatFile(filePath)
						if err == nil {
							result = content
						}

						toolResult := ""
						toolResult += "## File content " + filePath + "\n"
						toolResult += result + "\n"
						toolResult += "## File content\n"

						*response = *response + toolResult
					}

					messages = append(messages, &genai.Content{
						Role: "tool",
						Parts: []*genai.Part{
							{
								FunctionResponse: &genai.FunctionResponse{
									ID:   functionCall.ID,
									Name: functionCall.Name,
									Response: map[string]any{
										"content": result,
									},
								},
							},
						}})

					c.SendMessage(
						messages,
						response,
					)
				}
			}
		}
	}

	c.Messages = append(c.Messages, &genai.Content{Role: "assistant", Parts: []*genai.Part{{Text: text}}})
}

func (c *GoogleClient) CancelRequest() error {
	if c.cancelRequest != nil {
		c.cancelRequest()
	}
	return nil
}
