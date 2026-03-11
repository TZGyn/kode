package message

import (
	"math"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/TZGyn/kode/internal/markdown"
	"github.com/TZGyn/kode/internal/theme"
)

func card(role MessageRole, width int) lipgloss.Style {
	if role == User {
		b := lipgloss.ThickBorder()

		return lipgloss.NewStyle().Foreground(
			lipgloss.White,
		).Background(
			lipgloss.Color("#000000"),
		).Border(
			b,
			false,
			false,
			false,
			true,
		).BorderLeftForeground(lipgloss.Color("#2d8fff")).Width(width - 4)
	}

	if role == Assistant {
		return lipgloss.NewStyle().Foreground(
			lipgloss.Color("#ffffff"),
		).Background(
			lipgloss.Color("#000000"),
		).Width(width - 4)
	}

	return lipgloss.NewStyle()
}

func RenderParts(role MessageRole, width int, c *Content) string {
	card := card(role, width)

	var content strings.Builder

	for _, part := range *c {
		switch p := part.(type) {
		case TextPart:
			content.WriteString(RenderTextPart(role, width, p.String()))
		case ReasoningPart:
			content.WriteString(RenderReasonPart(role, width, p.String()))
		}
		content.WriteString("\n")
	}

	return card.Render(content.String())
}

func RenderReasonPart(role MessageRole, width int, reasoning string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Italic(true).Render(reasoning)
}

func RenderTextPart(role MessageRole, width int, content string) string {
	t := theme.GetTheme()

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
			t.Background(),
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
