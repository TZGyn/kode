package message

import (
	"math"

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

func Render(role MessageRole, width int, content string) string {
	t := theme.GetTheme()

	if role == User {
		return card(role, width).Render(lipgloss.NewStyle().Foreground(
			lipgloss.White,
		).Background(
			lipgloss.Color("#131313"),
		).Padding(
			1,
		).Width(int(math.Min(float64(width-4), 100))).Render(content))
	}

	if role == Assistant {
		card := card(role, width)
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

		return card.Render(message.Render(content))
	}

	return card(role, width).Render(lipgloss.NewStyle().Render(content))
}
