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

type FinishReason string

const (
	FinishReasonEndTurn          FinishReason = "end_turn"
	FinishReasonMaxTokens        FinishReason = "max_tokens"
	FinishReasonToolUse          FinishReason = "tool_use"
	FinishReasonCanceled         FinishReason = "canceled"
	FinishReasonError            FinishReason = "error"
	FinishReasonPermissionDenied FinishReason = "permission_denied"

	// Should never happen
	FinishReasonUnknown FinishReason = "unknown"
)

type FinishPart struct {
	Reason  FinishReason `json:"reason"`
	Time    int64        `json:"time"`
	Message string       `json:"message,omitempty"`
	Details string       `json:"details,omitempty"`
}

func (f FinishPart) String() string {
	return string(f.Reason) + " " + f.Details
}
