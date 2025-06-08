package model

import "google.golang.org/genai"

type ChatMessages []*ChatMessage

type ChatMessage struct {
	Role  string
	Parts []*ChatPart
}

type ChatPart struct {
	Type           string
	Text           string
	Reasoning      string
	ToolCallName   string
	ToolCallID     string
	ToolCallArgs   map[string]any
	ToolCallResult map[string]any
}

func (c *ChatMessages) ConvertToGoogleMessages() ([]*genai.Content, error) {
	googleMessages := []*genai.Content{}

	for _, content := range *c {
		parts := []*genai.Part{}

		for _, part := range content.Parts {
			if part.Type == "text" {
				parts = append(parts, &genai.Part{
					Text: part.Text,
				})
			}
			if part.Type == "tool-call" {
				parts = append(parts, &genai.Part{
					FunctionCall: &genai.FunctionCall{
						Name: part.ToolCallName,
						ID:   part.ToolCallID,
						Args: part.ToolCallArgs,
					},
				})
			}
			if part.Type == "tool-result" {
				parts = append(parts, &genai.Part{
					FunctionResponse: &genai.FunctionResponse{
						Name:     part.ToolCallName,
						ID:       part.ToolCallID,
						Response: part.ToolCallResult,
					},
				})
			}
		}

		googleMessages = append(googleMessages, &genai.Content{
			Role:  content.Role,
			Parts: parts,
		})
	}

	return googleMessages, nil
}

func (c *ChatMessages) AddGoogleMessages(messages []*genai.Content) error {
	for _, content := range messages {
		role := ""
		if content.Role == "user" {
			role = "user"
		} else {
			role = "assistant"
		}

		parts := []*ChatPart{}

		for _, part := range content.Parts {
			if len(part.Text) > 0 {
				parts = append(parts, &ChatPart{
					Type: "text",
					Text: part.Text,
				})
			}
			if part.FunctionCall != nil {
				parts = append(parts, &ChatPart{
					Type:         "tool-call",
					ToolCallName: part.FunctionCall.Name,
					ToolCallID:   part.FunctionCall.ID,
					ToolCallArgs: part.FunctionCall.Args,
				})
			}
			if part.FunctionResponse != nil {
				parts = append(parts, &ChatPart{
					Type:           "tool-result",
					ToolCallName:   part.FunctionResponse.Name,
					ToolCallID:     part.FunctionResponse.ID,
					ToolCallResult: part.FunctionResponse.Response,
				})
			}
		}

		*c = append(*c, &ChatMessage{Role: role, Parts: parts})
	}

	return nil
}
