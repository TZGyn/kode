package model

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TZGyn/kode/internal/app"
	"github.com/TZGyn/kode/internal/components/commandlist"
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
	commandList commandlist.CommandList

	showCommandModal bool

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

	vp := viewport.CreateViewport("", app, &prompt, &layout)

	cl := commandlist.New()

	return Model{
		spinner:     s,
		App:         app,
		Layout:      &layout,
		viewport:    vp,
		homeinput:   ta,
		promptInput: &prompt,
		commandList: cl,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Init(), tea.RequestWindowSize)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updateInput := true
	var cmds []tea.Cmd

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.App.Agent.HasPressedEsc {
			} else {
				m.App.Agent.HasPressedEsc = true
			}
		case "enter":
			if m.App.Agent.Generating {
				updateInput = false
				break
			}

			if m.promptInput.Value() == "/exit" {
				return m, tea.Quit
			}
			if m.App.Session.ID != "" {
				if m.promptInput.Value() != "" {
					userMessage := message.NewMessageWithText(message.User, m.promptInput.Value(), m.Layout)
					assistantMessage := message.NewEmptyMessage(message.Assistant, m.Layout)

					userMessage.UIString = userMessage.ToUIString()
					*m.App.Messages = append(
						*m.App.Messages,
						userMessage,
						assistantMessage,
					)

					go func() {
						m.App.Agent.Generate(m.promptInput.Value())
					}()
				}
			} else {
				m.App.Session = app.NewSession()

				if m.homeinput.Value() != "" {
					userMessage := message.NewMessageWithText(message.User, m.homeinput.Value(), m.Layout)
					assistantMessage := message.NewEmptyMessage(message.Assistant, m.Layout)

					userMessage.UIString = userMessage.ToUIString()
					*m.App.Messages = append(
						*m.App.Messages,
						userMessage,
						assistantMessage,
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
		var requireReload = false
		if m.Layout.Width != msg.Width || m.Layout.Height != msg.Height {
			requireReload = true
		}
		if requireReload {
			// msg.Height -= 2 // Make space for the status bar
			m.Layout.Width, m.Layout.Height = msg.Width, msg.Height

			m.App.RerenderMessages()

			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Reload(msg)

			cmds = append(cmds, cmd)
		}
	case reloadMsg:
		var cmd tea.Cmd

		if m.viewport.AtBottom() {
			m.viewport, cmd = m.viewport.ReloadAndScrollDown(msg)
		} else {
			m.viewport, cmd = m.viewport.Reload(msg)
		}

		cmds = append(cmds, cmd)
	case errMsg:
		m.err = msg
		return m, nil

	default:
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	shouldUpdatePrompt := true
	if m.commandList.IsVisible() {
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.String() == "enter" {
				shouldUpdatePrompt = false
			}
		}
	}

	if updateInput {
		m.homeinput, cmd = m.homeinput.Update(msg)
		cmds = append(cmds, cmd)

		if shouldUpdatePrompt {
			m.promptInput, cmd = m.promptInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	prompt := m.promptInput.Value()
	if strings.HasPrefix(prompt, "/") {
		if !m.commandList.IsVisible() {
			m.commandList.Show(prompt)
		} else {
			m.commandList.Filter(prompt)
		}
	} else {
		if m.commandList.IsVisible() {
			m.commandList.Hide()
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	m.commandList, cmd = m.commandList.Update(msg)
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
		chatLayer := lipgloss.NewLayer(m.chat())

		layers := []*lipgloss.Layer{
			chatLayer,
		}
		if m.commandList.IsVisible() {
			modalY := 5
			layers = append(layers, lipgloss.NewLayer(m.modal()).X(0).Y(modalY))
		}

		layers = append(layers, m.promptLayer())

		comp := lipgloss.NewCompositor(layers...)

		view = tea.NewView(comp.Render())
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
	mainLayout := lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Render(
		lipgloss.Place(
			m.Layout.Width,
			m.Layout.Height,
			lipgloss.Left,
			lipgloss.Bottom,
			m.viewport.View().Content+"\n\n"+"\n\n\n\n\n"+"\n"+m.statusBar()+"\n",
			lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Background(lipgloss.Color("#000000"))),
		),
	)

	return mainLayout
}

func (m Model) prompt() string {
	return lipgloss.NewStyle().
		Padding(1).
		Background(lipgloss.Color("#000000")).
		Render(m.promptInput.View().Content)
}

func (m Model) promptLayer() *lipgloss.Layer {
	y := m.Layout.Height - 9

	if m.promptInput.Height > prompt.MinHeight {
		y -= (m.promptInput.Height - prompt.MinHeight)
	}

	return lipgloss.NewLayer(m.prompt()).X(0).Y(y)
}

func (m Model) statusBar() string {
	content := ""

	if m.App.Agent.Generating {
		content += m.spinner.View().Content + " "

		if m.App.Agent.HasPressedEsc {
			content += "esc again to confirm"
		} else {
			content += "esc interrupt"
		}
	}

	return lipgloss.NewStyle().Padding(0, 2).Render(content)
}

func (m Model) modal() string {
	return m.commandList.View()
}
