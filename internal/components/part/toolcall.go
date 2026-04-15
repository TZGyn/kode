package part

type ToolCallPart struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Input            string `json:"input"`
	ProviderExecuted bool   `json:"provider_executed"`
	Finished         bool   `json:"finished"`
	Time             int64  `json:"time"`
}

func (tc ToolCallPart) String() string {
	return tc.Name + "\n" + tc.Input
}
