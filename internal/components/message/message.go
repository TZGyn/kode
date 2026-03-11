package message

import (
	"github.com/TZGyn/kode/internal/layout"
)

type Message struct {
	ID        string
	MessageID string
	SessionID string

	Layout *layout.Layout

	Role MessageRole

	Content *Content

	UIString string
}

func (m Message) ToUIString() string {
	return RenderParts(m.Role, m.Layout.Width, m.Content)
}
