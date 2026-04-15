package message

type MessageRole string

const (
	Assistant MessageRole = "assistant"
	User      MessageRole = "user"
	System    MessageRole = "system"
	Tool      MessageRole = "tool"
)
