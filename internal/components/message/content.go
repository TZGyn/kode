package message

import "strings"

type Content []ContentPart

func (c *Content) AppendTextDelta(delta string) {
	if len(*c) == 0 {
		*c = append(*c, TextPart{Content: delta})
		return
	}

	lastPart := (*c)[len(*c)-1]

	if p, ok := lastPart.(TextPart); ok {
		(*c)[len(*c)-1] = TextPart{Content: p.Content + delta}
	} else {
		*c = append(*c, TextPart{Content: delta})
	}
}

func (c *Content) String() string {
	var content strings.Builder

	for _, part := range *c {
		content.WriteString(part.String())
	}

	return content.String()
}
