package google

import (
	"context"
	"strings"

	"github.com/TZGyn/kode/internal/tool"
	"google.golang.org/genai"
)

type Config struct {
	GEMINI_API_KEY string `json:"GEMINI_API_KEY"`
}

func DefaultConfig(apiKey string) Config {
	return Config{GEMINI_API_KEY: apiKey}
}

func CreateGoogle(ctx context.Context, config Config) (*genai.Client, *genai.Chat, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.GEMINI_API_KEY,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, nil, err
	}

	chat, err := client.Chats.Create(
		ctx,
		"gemini-2.0-flash",
		&genai.GenerateContentConfig{
			// SystemInstruction: &genai.Content{},
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
					},
				},
			},
		},
		nil,
	)

	if err != nil {
		return nil, nil, err
	}

	return client, chat, nil
}

func SendMessage(ctx context.Context, chat *genai.Chat, part genai.Part, status *string, response *string) {
	stream := chat.SendMessageStream(
		ctx,
		part,
	)
	for result, err := range stream {
		if err != nil {
			*status = "done"
			return
		}

		var texts []string
		for _, part := range result.Candidates[0].Content.Parts {
			if part.Text != "" {
				if part.Thought {
					continue
				}
				texts = append(texts, part.Text)
			}
		}

		*response = *response + strings.Join(texts, "")

		if len(result.FunctionCalls()) == 1 {
			functionCall := result.FunctionCalls()[0]
			if functionCall.Name == "list-directory" {
				result := []string{}

				directory, ok := functionCall.Args["directory"].(string)
				if ok {
					entires, err := tool.ListDirectory(directory)
					if err == nil {
						result = entires
					}
				}

				SendMessage(
					ctx,
					chat,
					genai.Part{
						FunctionResponse: &genai.FunctionResponse{
							ID:   functionCall.ID,
							Name: functionCall.Name,
							Response: map[string]any{
								"children": result,
							},
						},
					},
					status,
					response,
				)
			}
		}
	}
	*status = "done"
}
