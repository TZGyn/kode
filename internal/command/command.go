package command

type Command struct {
	CommandID   string
	Title       string
	Description string
}

var Commands = []Command{
	{
		CommandID:   "/exit",
		Title:       "Exit Application",
		Description: "Save session and quit the application",
	},
}
