package app

type App struct {
	ID string

	Session Session

	Messages []string
}

type Session struct {
	ID string
}

func NewApp() *App {
	return &App{}
}
