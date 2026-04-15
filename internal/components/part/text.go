package part

type TextPart struct {
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

func (p TextPart) String() string {
	return p.Content
}

func (t TextPart) AppendText(delta string) TextPart {
	t.Content = t.Content + delta
	return t
}
