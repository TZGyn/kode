package model

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/TZGyn/kode/internal/animation"
	"github.com/TZGyn/kode/internal/google"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"google.golang.org/genai"
)

type state int

const (
	startState state = iota
	requestState
	responseState
	doneState
)

type ChatModel struct {
	anim   tea.Model
	state  state
	status string

	GoogleClient *google.GoogleClient
	messages     ChatMessages

	Prompt   string
	Response string

	glam         *glamour.TermRenderer
	glamHeight   int
	glamViewport viewport.Model
	glamOutput   string
	width        int
	height       int

	renderer *lipgloss.Renderer
}

type ChatConfig struct {
	GEMINI_API_KEY string `json:"GEMINI_API_KEY"`
}

type initMsg struct{}
type generatingMsg struct{}
type receivingMsg struct{}

func InitialModel(prompt string, messages ChatMessages, config ChatConfig) *ChatModel {
	gr, _ := glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	vp := viewport.New(0, 0)
	vp.GotoBottom()

	renderer := lipgloss.NewRenderer(os.Stderr, termenv.WithColorCache(true))

	googleConfig := google.DefaultConfig(config.GEMINI_API_KEY)

	client, err := google.CreateGoogle(googleConfig)

	if err != nil {
		log.Fatal(err)
	}

	return &ChatModel{
		state: startState,

		GoogleClient: client,
		messages:     messages,

		Prompt:       prompt,
		status:       "generating",
		glam:         gr,
		glamViewport: vp,
		renderer:     renderer,
	}
}

func (m *ChatModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."

	return func() tea.Msg { return initMsg{} }
}
func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case initMsg:
		m.state = requestState
		m.anim = animation.NewAnim("Generating")
		cmds = append(cmds, m.anim.Init(), func() tea.Msg { return generatingMsg{} })
	case generatingMsg:
		go func(model *ChatModel) {
			googleMessages, err := model.messages.ConvertToGoogleMessages()
			if err == nil {
				model.GoogleClient.Messages = append(model.GoogleClient.Messages, googleMessages...)
			}

			model.GoogleClient.Messages = append(model.GoogleClient.Messages, &genai.Content{Role: "assistant", Parts: []*genai.Part{{Text: model.Prompt}}})
			model.GoogleClient.SendMessage(
				model.GoogleClient.Messages,
				&model.Response,
			)
			model.status = "done"
		}(m)
		cmds = append(cmds, func() tea.Msg { return receivingMsg{} })
	case receivingMsg:
		m.state = responseState

		if m.Response != "" {
			wasAtBottom := m.glamViewport.ScrollPercent() == 1.0
			oldHeight := m.glamHeight

			var err error
			m.glamOutput, err = m.glam.Render(m.Response)
			if err != nil {
				fmt.Println(err)
			}

			m.glamOutput = strings.TrimRightFunc(m.glamOutput, unicode.IsSpace)
			m.glamOutput = strings.ReplaceAll(m.glamOutput, "\t", strings.Repeat(" ", 4))

			m.glamHeight = lipgloss.Height(m.glamOutput)

			truncatedGlamOutput := m.renderer.NewStyle().
				MaxWidth(m.width).
				Render(m.glamOutput)

			m.glamViewport.SetContent(truncatedGlamOutput)

			if oldHeight < m.glamHeight && wasAtBottom {
				// If the viewport's at the bottom and we've received a new
				// line of content, follow the output by auto scrolling to
				// the bottom.
				m.glamViewport.GotoBottom()
			}
		}

		if m.status == "done" {
			m.state = doneState

			return m, m.quit
		}

		cmds = append(cmds, func() tea.Msg { return receivingMsg{} })
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.glamViewport.Width = m.width
		m.glamViewport.Height = m.height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, m.quit
		}
	}

	if m.state == requestState {
		var cmd tea.Cmd
		m.anim, cmd = m.anim.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.viewportNeeded() {
		var cmd tea.Cmd
		m.glamViewport, cmd = m.glamViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *ChatModel) quit() tea.Msg {
	if m.GoogleClient != nil {
		m.GoogleClient.CancelRequest()
	}
	return tea.Quit()
}

func (m *ChatModel) viewportNeeded() bool {
	return m.glamHeight > m.height
}

func (m *ChatModel) View() string {
	switch m.state {
	case requestState:
		return m.anim.View()
	case responseState:
		if m.viewportNeeded() {
			return m.glamViewport.View() + "\n" + strings.Join(m.GoogleClient.FunctionCalls, " ")
		}

		return m.glamOutput + "\n" + strings.Join(m.GoogleClient.FunctionCalls, " ")
	case doneState:
		return ""
	}

	return ""
}
