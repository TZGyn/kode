package part

type ToolResultPart struct {
	ToolCallID string `json:"tool_call_id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	Data       string `json:"data"`
	MIMEType   string `json:"mime_type"`
	Metadata   string `json:"metadata"`
	IsError    bool   `json:"is_error"`
	Time       int64  `json:"time"`
}

func (tr ToolResultPart) String() string {
	return tr.Content
}
