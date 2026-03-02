package app

import "github.com/TZGyn/kode/internal/components/message"

type App struct {
	ID string

	Session Session

	Messages []message.MessagePart
}

type Session struct {
	ID string
}

func NewApp() *App {
	return &App{}
}
