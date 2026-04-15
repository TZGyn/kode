package part

import (
	"fmt"
	"time"

	"charm.land/fantasy"
)

type ContentPart interface {
	String() string
}

func NewTextPart(content string) TextPart {
	return TextPart{
		Content: content,
		Time:    time.Now().UTC().UnixMilli(),
	}
}

func NewReasoningPart(reason string) ReasoningPart {
	return ReasoningPart{
		Reason: reason,
		Time:   time.Now().UTC().UnixMilli(),
	}
}

func NewFinishPart(reason FinishReason, message, details string) FinishPart {
	return FinishPart{
		Reason:  reason,
		Message: message,
		Details: details,
		Time:    time.Now().UTC().UnixMilli(),
	}
}

func NewToolCallPart(ID string, name string, input string) ToolCallPart {
	return ToolCallPart{
		ID:               ID,
		Name:             name,
		Input:            input,
		ProviderExecuted: true,
		Finished:         true,
		Time:             time.Now().UTC().UnixMilli(),
	}
}

func NewToolResultPart(result fantasy.ToolResultContent) ToolResultPart {
	part := ToolResultPart{
		ToolCallID: result.ToolCallID,
		Name:       result.ToolName,
		Metadata:   result.ClientMetadata,
		Time:       time.Now().UTC().UnixMilli(),
	}

	switch result.Result.GetType() {
	case fantasy.ToolResultContentTypeText:
		if r, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentText](result.Result); ok {
			part.Content = r.Text
		}
	case fantasy.ToolResultContentTypeError:
		if r, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentError](result.Result); ok {
			part.Content = r.Error.Error()
			part.IsError = true
		}
	case fantasy.ToolResultContentTypeMedia:
		if r, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentMedia](result.Result); ok {
			content := r.Text
			if content == "" {
				content = fmt.Sprintf("Loaded %s content", r.MediaType)
			}
			part.Content = content
			part.Data = r.Data
			part.MIMEType = r.MediaType
		}
	}

	return part
}
