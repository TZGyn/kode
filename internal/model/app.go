package model

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TZGyn/kode/internal/app"
	"github.com/TZGyn/kode/internal/components/message"
	"github.com/TZGyn/kode/internal/components/prompt"
	"github.com/TZGyn/kode/internal/components/spinner"
	"github.com/TZGyn/kode/internal/components/viewport"
	"github.com/TZGyn/kode/internal/layout"
)

type errMsg error

type Model struct {
	App    *app.App
	Layout *layout.Layout

	spinner     *spinner.Spinner
	homeinput   textarea.Model
	promptInput *prompt.PromptComponent
	viewport    viewport.Model

	// use to request window size every 5 frames in windows
	// due to windows not having terminal resize message
	frame int

	err error
}

type reloadMsg struct{}

func InitAppModel() Model {
	s := spinner.New()

	ta := textarea.New()

	ta.KeyMap.InsertNewline.SetEnabled(false)

	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.ShowLineNumbers = false

	ta.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#2d8fff")).Render("┃ ")

	style := ta.Styles()

	style.Focused.Base = lipgloss.NewStyle()
	style.Focused.CursorLine = lipgloss.NewStyle()

	ta.SetStyles(style)

	ta.CharLimit = 280

	prompt := prompt.NewPrompt()

	layout := layout.Layout{}

	app := app.NewApp(&layout)

	viewport := viewport.CreateViewport("", app, &prompt, &layout)

	return Model{
		spinner:     s,
		App:         app,
		Layout:      &layout,
		viewport:    viewport,
		homeinput:   ta,
		promptInput: &prompt,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.promptInput.Value() == "/exit" {
				return m, tea.Quit
			}
			if m.App.Session.ID != "" {
				if m.promptInput.Value() != "" {
					*m.App.Messages = append(
						*m.App.Messages,
						message.Message{
							ID:        "hello",
							MessageID: "hello",
							SessionID: "hello",
							Role:      message.User,
							Content: &message.Content{
								message.TextPart{
									Content: m.promptInput.Value(),
								},
							},
							Layout: m.Layout,
						},
						message.Message{
							ID:        "hello",
							MessageID: "hello",
							SessionID: "hello",
							Role:      message.Assistant,
							Content:   &message.Content{},
							Layout:    m.Layout,
						},
					)

					go func() {
						m.App.Agent.Generate(m.promptInput.Value())
					}()
				}
			} else {
				m.App.Session.ID = "hello"

				if m.homeinput.Value() != "" {
					*m.App.Messages = append(
						*m.App.Messages,
						message.Message{
							ID:        "hello",
							MessageID: "hello",
							SessionID: "hello",
							Role:      message.User,
							Content: &message.Content{
								message.TextPart{
									Content: m.homeinput.Value(),
								},
							},
							Layout: m.Layout,
						},
						message.Message{
							ID:        "hello",
							MessageID: "hello",
							SessionID: "hello",
							Role:      message.Assistant,
							Content:   &message.Content{},
							Layout:    m.Layout,
						},
					)

					go func() {
						m.App.Agent.Generate(m.homeinput.Value())
					}()
				}
			}

			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.ReloadAndScrollDown(msg)

			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		// msg.Height -= 2 // Make space for the status bar
		m.Layout.Width, m.Layout.Height = msg.Width, msg.Height

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Reload(msg)

		cmds = append(cmds, cmd)
	case reloadMsg:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.ReloadAndScrollDown(msg)

		cmds = append(cmds, cmd)
	case errMsg:
		m.err = msg
		return m, nil

	default:
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.homeinput, cmd = m.homeinput.Update(msg)
	cmds = append(cmds, cmd)

	m.promptInput, cmd = m.promptInput.Update(msg)
	cmds = append(cmds, cmd)

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	m.frame += 1
	m.frame %= 5

	if m.frame == 0 {
		cmds = append(cmds, tea.RequestWindowSize)
	}

	if m.App.Agent.RequestRefresh {
		cmds = append(cmds, func() tea.Msg { return reloadMsg{} })
		m.App.Agent.RequestRefresh = false
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var view tea.View

	if m.err != nil {
		view = tea.NewView(m.err.Error())
	}
	if m.App.Session.ID == "" {
		view = tea.NewView(m.home())
	} else {
		view = tea.NewView(m.chat())
	}

	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}

func (m Model) home() string {

	kode := `
░██     ░██   ░██████   ░███████   ░██████████ 
░██    ░██   ░██   ░██  ░██   ░██  ░██         
░██   ░██   ░██     ░██ ░██    ░██ ░██         
░███████    ░██     ░██ ░██    ░██ ░█████████  
░██   ░██   ░██     ░██ ░██    ░██ ░██         
░██    ░██   ░██   ░██  ░██   ░██  ░██         
░██     ░██   ░██████   ░███████   ░██████████ 
`

	mainLayout := lipgloss.Place(
		m.Layout.Width,
		m.Layout.Height+2,
		lipgloss.Center,
		lipgloss.Center,
		kode+"\n"+fmt.Sprintf("%s", m.homeinput.View()),
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Background(lipgloss.Color("#000000"))),
	)

	return mainLayout
}

func (m Model) chat() string {
	mainLayout := lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Render(lipgloss.Place(
		m.Layout.Width,
		m.Layout.Height,
		lipgloss.Left,
		lipgloss.Bottom,
		m.viewport.View().Content+"\n\n"+m.prompt()+"\n"+m.statusBar()+"\n",
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Background(lipgloss.Color("#000000"))),
	))

	return mainLayout
}

func (m Model) prompt() string {
	return lipgloss.NewStyle().Padding(1).Background(lipgloss.Color("#000000")).Render(m.promptInput.View().Content)
}

func (m Model) statusBar() string {
	return lipgloss.NewStyle().Padding(0, 2).Render(m.spinner.View().Content)
}
