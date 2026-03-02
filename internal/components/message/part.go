package message

type MessageRole string

const (
	Assistant MessageRole = "assistant"
	User      MessageRole = "user"
	System    MessageRole = "system"
	Tool      MessageRole = "tool"
)

type ContentPart interface {
	String() string
}

type TextPart struct {
	Content string
}

func (p TextPart) String() string {
	return p.Content
}

type ReasoningPart struct {
	Reason string
}

func (p ReasoningPart) String() string {
	return p.Reason
}
