package message

import (
	"math"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/TZGyn/kode/internal/markdown"
)

func RenderParts(role MessageRole, width int, c *Content) string {
	card := card(role, width)

	var content strings.Builder

	for i, part := range *c {
		switch p := part.(type) {
		case TextPart:
			content.WriteString(RenderTextPart(role, width, p.String()))
		case ReasoningPart:
			content.WriteString(RenderReasonPart(role, width, p.String()))
		case FinishPart:
			content.WriteString(RenderFinishPart(role, width, p.Reason, p.Message, p.Details))
		}
		if i != len((*c))-1 {
			content.WriteString("\n\n")
		}
	}

	return container(width).Render(card.Render(content.String()))
}

func RenderReasonPart(role MessageRole, width int, reasoning string) string {
	b := lipgloss.ThickBorder()

	reasoningPrefix := lipgloss.NewStyle().Foreground(lipgloss.Yellow).Italic(true).Render("Thinking:")

	return lipgloss.NewStyle().Foreground(
		lipgloss.Color("#444444"),
	).Background(
		lipgloss.Color("#000000"),
	).Padding(
		1,
	).Border(
		b,
		false,
		false,
		false,
		true,
	).BorderLeftForeground(
		lipgloss.Color("#444444"),
	).Width(
		int(math.Min(float64(width-4), 100)),
	).Render(reasoningPrefix + " " + strings.TrimSpace(reasoning))
}

func RenderTextPart(role MessageRole, width int, content string) string {
	// t := theme.GetTheme()

	if role == User {
		return lipgloss.NewStyle().Foreground(
			lipgloss.White,
		).Background(
			lipgloss.Color("#131313"),
		).Padding(
			1,
		).Width(int(math.Min(float64(width-4), 100))).Render(content)
	}

	if role == Assistant {
		message := lipgloss.NewStyle().Foreground(
			lipgloss.White,
		).Background(
			lipgloss.Color("#000000"),
			// t.Background(),
		).Padding(
			1,
		).Width(
			int(math.Min(float64(width-4), 100)),
		)

		content := markdown.ToMarkdown(
			content,
			message.GetWidth(),
		)

		return message.Render(content)
	}

	return lipgloss.NewStyle().Render(content)
}

func RenderFinishPart(role MessageRole, width int, reason FinishReason, title string, details string) string {
	message := lipgloss.NewStyle().Foreground(
		lipgloss.Red,
	).Background(
		lipgloss.Color("#000000"),
		// t.Background(),
	).Padding(
		1,
	).Width(
		int(math.Min(float64(width-4), 100)),
	)

	return message.Render(title + "\n" + details)
}
