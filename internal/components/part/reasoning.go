package part

type ReasoningPart struct {
	Reason string `json:"reason"`
	Time   int64  `json:"time"`
}

func (p ReasoningPart) String() string {
	return p.Reason
}
