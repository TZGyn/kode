package message

import (
	"math"

	"charm.land/lipgloss/v2"
)

func container(width int) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(
		lipgloss.Color("#ffffff"),
	).Background(
		lipgloss.Color("#000000"),
	).Width(width - 4)
}

func card(role MessageRole, width int) lipgloss.Style {
	if role == User {
		b := lipgloss.ThickBorder()

		return lipgloss.NewStyle().Foreground(
			lipgloss.White,
		).Background(
			lipgloss.Color("#131313"),
		).Border(
			b,
			false,
			false,
			false,
			true,
		).BorderLeftForeground(lipgloss.Color("#2d8fff")).Width(
			int(math.Min(float64(width-4), 100)),
		)
	}

	if role == Assistant {
		return lipgloss.NewStyle().Foreground(
			lipgloss.Color("#ffffff"),
		).Background(
			lipgloss.Color("#000000"),
		).Width(
			int(math.Min(float64(width-4), 100)),
		)
	}

	return lipgloss.NewStyle()
}
