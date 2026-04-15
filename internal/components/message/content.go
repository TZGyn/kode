package message

import (
	"strings"

	"charm.land/fantasy"
	"github.com/TZGyn/kode/internal/components/part"
)

type Content []part.ContentPart

func (c *Content) AppendTextDelta(delta string) {
	if len(*c) == 0 {
		*c = append(*c, part.NewTextPart(delta))
		return
	}

	lastPart := (*c)[len(*c)-1]

	if p, ok := lastPart.(part.TextPart); ok {
		(*c)[len(*c)-1] = p.AppendText(delta)
	} else {
		*c = append(*c, part.NewTextPart(delta))
	}
}

func (c *Content) AppendReasonDelta(delta string) {
	if len(*c) == 0 {
		*c = append(*c, part.ReasoningPart{Reason: delta})
		return
	}

	lastPart := (*c)[len(*c)-1]

	if p, ok := lastPart.(part.ReasoningPart); ok {
		(*c)[len(*c)-1] = part.ReasoningPart{Reason: p.Reason + delta}
	} else {
		*c = append(*c, part.ReasoningPart{Reason: delta})
	}
}

func (c *Content) AppendFinish(reason part.FinishReason, message, details string) {
	*c = append(*c, part.NewFinishPart(
		reason,
		message,
		details,
	))
}

func (c *Content) AppendToolCall(ID string, name string, input string) {
	*c = append(*c, part.NewToolCallPart(
		ID,
		name,
		input,
	))
}

func (c *Content) AppendToolResult(result fantasy.ToolResultContent) {
	*c = append(*c, part.NewToolResultPart(result))
}

func (c *Content) String() string {
	var content strings.Builder

	for _, part := range *c {
		content.WriteString(part.String())
	}

	return content.String()
}
