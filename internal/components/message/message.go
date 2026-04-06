package message

import (
	"github.com/TZGyn/kode/internal/layout"
	"github.com/google/uuid"
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

func NewMessageWithText(role MessageRole, text string, layout *layout.Layout) Message {
	return Message{
		ID:        uuid.NewString(),
		MessageID: uuid.NewString(),
		SessionID: uuid.NewString(),
		Role:      role,
		Content: &Content{
			TextPart{
				Content: text,
			},
		},
		Layout: layout,
	}
}

func NewEmptyMessage(role MessageRole, layout *layout.Layout) Message {
	return Message{
		ID:        uuid.NewString(),
		MessageID: uuid.NewString(),
		SessionID: uuid.NewString(),
		Role:      role,
		Content:   &Content{},
		Layout:    layout,
	}
}

func (m Message) ToUIString() string {
	return RenderParts(m.Role, m.Layout.Width, m.Content)
}
