package app

import (
	"fmt"
	"os"

	"github.com/TZGyn/kode/internal/agent"
	"github.com/TZGyn/kode/internal/components/message"
	"github.com/TZGyn/kode/internal/layout"
	"github.com/joho/godotenv"
)

type App struct {
	ID string

	Session Session
	Agent   *agent.Agent

	Messages *[]message.Message

	Layout *layout.Layout
}

func NewApp(layout *layout.Layout) *App {
	var messages *[]message.Message = &[]message.Message{}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env found")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set OPENAI_API_KEY environment variable")
		os.Exit(1)
	}

	openrouterApiKey := os.Getenv("OPENROUTER_API_KEY")
	if openrouterApiKey == "" {
		fmt.Println("Please set OPENROUTER_API_KEY environment variable")
		os.Exit(1)
	}

	config := agent.NewConfig("openai", "gpt-5-mini", apiKey)
	config = agent.NewConfig("openrouter", "openai/gpt-oss-120b:free", openrouterApiKey)

	agent := agent.New(messages, layout, config)

	return &App{Agent: agent, Messages: messages, Layout: layout}
}

func (a *App) RerenderMessages() {
	for i, _ := range *a.Messages {
		(*a.Messages)[i].UIString = (*a.Messages)[i].ToUIString()
	}
}
