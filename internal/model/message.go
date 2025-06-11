package model

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"google.golang.org/genai"
)

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

func createKeyValuePairs(m map[string]any) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func (c *ChatMessages) ConvertToOpenAIMessages() ([]openai.ChatCompletionMessageParamUnion, error) {
	openAIMessages := []openai.ChatCompletionMessageParamUnion{}

	for _, content := range *c {

		for _, part := range content.Parts {
			if part.Type == "text" {
				openAIMessages = append(openAIMessages, openai.UserMessage(part.Text))
			}
			if part.Type == "tool-call" {
			}
			if part.Type == "tool-result" {
				openAIMessages = append(openAIMessages, openai.ToolMessage(createKeyValuePairs(part.ToolCallResult), part.ToolCallID))
			}
		}
	}

	return openAIMessages, nil
}

func (c *ChatMessages) AddOpenAIMessages(messages []openai.ChatCompletionMessageParamUnion) error {
	for _, message := range messages {
		parts := []*ChatPart{}
		if message.OfUser != nil {
			parts = append(parts, &ChatPart{
				Type: "text",
				Text: message.OfUser.Content.OfString.String(),
			})
		}
		if message.OfTool != nil {
			data, err := message.OfTool.Content.OfString.MarshalJSON()
			if err != nil {
				continue
			}
			var result map[string]any
			json.Unmarshal(data, &result)

			parts = append(parts, &ChatPart{
				Type:           "tool-result",
				ToolCallID:     message.OfTool.ToolCallID,
				ToolCallResult: result,
			})
		}
		if message.OfAssistant != nil {
			parts = append(parts, &ChatPart{
				Type: "text",
				Text: message.OfAssistant.Content.OfString.String(),
			})
			toolCalls := message.OfAssistant.ToolCalls
			if len(toolCalls) > 0 {

				for _, toolCall := range toolCalls {
					var args map[string]any
					json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

					parts = append(parts, &ChatPart{
						Type:         "tool-call",
						ToolCallID:   toolCall.ID,
						ToolCallName: toolCall.Function.Name,
						ToolCallArgs: args,
					})
				}
			}
		}

		*c = append(*c, &ChatMessage{Role: *message.GetRole(), Parts: parts})
	}
	return nil
}

func (c *ChatMessages) Print() {
	result := ""
	for _, message := range *c {
		result += "Role: " + message.Role + "\n"
		for _, part := range message.Parts {
			result += "Parts:"
			result += "\tType: " + part.Type + "\n"
			result += "\tReasoning: " + part.Reasoning + "\n"
			result += "\tText: " + part.Text + "\n"
			result += "\tToolCallID: " + part.ToolCallID + "\n"
			result += "\tToolCallName: " + part.ToolCallName + "\n"
			if len(part.ToolCallArgs) != 0 {
				result += "\tToolCallArgs:\n"
				for key, value := range part.ToolCallArgs {
					result += "\t\t" + key + ":" + fmt.Sprintf(" %v\n", value)
				}
			}
			if len(part.ToolCallResult) != 0 {
				result += "\tToolCallResult:\n"
				for key, value := range part.ToolCallResult {
					result += "\t\t" + key + ":" + fmt.Sprintf(" %v\n", value)
				}
			}
			result += "\n"
		}
	}

	fmt.Println("Result", result)
}
