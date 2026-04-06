package message

import (
	"strings"
	"time"
)

type Content []ContentPart

func (c *Content) AppendTextDelta(delta string) {
	if len(*c) == 0 {
		*c = append(*c, TextPart{Content: delta})
		return
	}

	lastPart := (*c)[len(*c)-1]

	if p, ok := lastPart.(TextPart); ok {
		(*c)[len(*c)-1] = TextPart{Content: p.Content + delta}
	} else {
		*c = append(*c, TextPart{Content: delta})
	}
}

func (c *Content) AppendReasonDelta(delta string) {
	if len(*c) == 0 {
		*c = append(*c, ReasoningPart{Reason: delta})
		return
	}

	lastPart := (*c)[len(*c)-1]

	if p, ok := lastPart.(ReasoningPart); ok {
		(*c)[len(*c)-1] = ReasoningPart{Reason: p.Reason + delta}
	} else {
		*c = append(*c, ReasoningPart{Reason: delta})
	}
}

func (c *Content) AppendFinish(reason FinishReason, message, details string) {
	*c = append(*c, FinishPart{
		Reason:  reason,
		Time:    time.Now().UTC().UnixMilli(),
		Message: message,
		Details: details,
	})
}

func (c *Content) String() string {
	var content strings.Builder

	for _, part := range *c {
		content.WriteString(part.String())
	}

	return content.String()
}
