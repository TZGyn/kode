package message

import "github.com/charmbracelet/lipgloss"

var UserStyle = lipgloss.NewStyle().
	MarginTop(1).
	BorderLeft(true).
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("#5C9CF5"))

var AssistantStyle = lipgloss.NewStyle().
	MarginTop(1).
	MarginBottom(1).
	BorderLeft(true).
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("#F57FE0"))
