package part

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
	Message string       `json:"message,omitempty"`
	Details string       `json:"details,omitempty"`
	Time    int64        `json:"time"`
}

func (f FinishPart) String() string {
	return string(f.Reason) + " " + f.Details
}
