package command

import "strings"

type Command struct {
	CommandID   string
	Title       string
	Description string
}

var Commands = []Command{
	{
		CommandID:   "/exit",
		Title:       "Exit",
		Description: "Save session and quit",
	},
	{
		CommandID:   "/new",
		Title:       "New Session",
		Description: "Start a new chat session",
	},
	{
		CommandID:   "/clear",
		Title:       "Clear",
		Description: "Clear the chat history",
	},
	{
		CommandID:   "/model",
		Title:       "Model",
		Description: "Switch AI model",
	},
	{
		CommandID:   "/help",
		Title:       "Help",
		Description: "Show available commands",
	},
}

func FilterCommands(query string) []Command {
	if query == "" {
		return Commands
	}
	var filtered []Command
	for _, cmd := range Commands {
		if strings.Contains(strings.ToLower(cmd.CommandID), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(cmd.Title), strings.ToLower(query)) {
			filtered = append(filtered, cmd)
		}
	}
	return filtered
}
