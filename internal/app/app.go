package app

import (
	"github.com/TZGyn/kode/internal/agent"
	"github.com/TZGyn/kode/internal/components/message"
	"github.com/TZGyn/kode/internal/layout"
)

type App struct {
	ID string

	Session Session
	Agent   *agent.Agent

	Messages *[]message.Message

	Layout *layout.Layout
}

type Session struct {
	ID string
}

func NewApp(layout *layout.Layout) *App {
	var messages *[]message.Message = &[]message.Message{}

	agent := agent.New(messages, layout)

	return &App{Agent: agent, Messages: messages, Layout: layout}
}

func (a *App) RerenderMessages() {
	for i, _ := range *a.Messages {
		(*a.Messages)[i].UIString = (*a.Messages)[i].ToUIString()
	}
}
