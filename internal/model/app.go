package model

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/TZGyn/kode/internal/app"
	"github.com/TZGyn/kode/internal/components/viewport"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Model struct {
	App *app.App

	spinner spinner.Model

	homeinput textarea.Model
	chatinput textarea.Model

	viewport viewport.Model

	width    int
	height   int
	quitting bool
	err      error
}

func InitAppModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	app := app.NewApp()

	ta := textarea.New()

	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.ShowLineNumbers = false

	ta.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#2d8fff")).Render("┃ ")
	ta.CharLimit = 280

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	chatinput := textarea.New()

	chatinput.Placeholder = "Send a message..."
	chatinput.ShowLineNumbers = false

	chatinput.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#2d8fff")).Render("┃ ")
	chatinput.CharLimit = 3000

	// Remove cursor line styling
	chatinput.FocusedStyle.CursorLine = lipgloss.NewStyle()

	chatinput.KeyMap.InsertNewline.SetEnabled(false)
	chatinput.Focus()

	viewport := viewport.CreateViewport("", app)

	return Model{
		spinner:   s,
		App:       app,
		viewport:  viewport,
		homeinput: ta,
		chatinput: chatinput,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var (
		tiCmd tea.Cmd
	)

	m.homeinput, tiCmd = m.homeinput.Update(msg)
	cmds = append(cmds, tiCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.App.Session.ID = "hello"
			m.App.Messages = append(m.App.Messages, m.chatinput.Value())

			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Reload(msg)

			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		// msg.Height -= 2 // Make space for the status bar
		m.width, m.height = msg.Width, msg.Height

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	var cmd tea.Cmd
	m.chatinput, cmd = m.chatinput.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	if m.App.Session.ID == "" {
		return m.home()
	} else {
		return m.chat()
	}
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
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		kode+"\n"+fmt.Sprintf("%s", m.homeinput.View()),
		lipgloss.WithWhitespaceBackground(lipgloss.Color("#000000")),
	)

	return mainLayout
}

func (m Model) chat() string {
	mainLayout := lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Padding(1).Render(lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Bottom,
		m.viewport.View()+"\n"+m.chatinput.View(),
		lipgloss.WithWhitespaceBackground(lipgloss.Color("#000000")),
	))

	return mainLayout
}
