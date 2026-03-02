package message

import "github.com/TZGyn/kode/internal/layout"

type MessagePart struct {
	ID        string
	MessageID string
	SessionID string

	Layout *layout.Layout

	Role MessageRole

	Content ContentPart
}

func (m MessagePart) ToUIString() string {
	return Render(m.Role, m.Layout.Width, m.Content.String())
}
