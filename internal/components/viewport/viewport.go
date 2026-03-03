package viewport

// An example program demonstrating the pager component from the Bubbles
// component library.

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TZGyn/kode/internal/app"
	"github.com/TZGyn/kode/internal/components/prompt"
	"github.com/TZGyn/kode/internal/layout"
)

type Model struct {
	App         *app.App
	promptInput *prompt.PromptComponent

	content  string
	viewport viewport.Model

	layout *layout.Layout
}

func CreateViewport(content string, app *app.App, promptInput *prompt.PromptComponent, layout *layout.Layout) Model {
	vp := viewport.New()

	keymap := viewport.DefaultKeyMap()
	keymap.Up.SetKeys("up")
	keymap.Down.SetKeys("down")
	keymap.HalfPageUp.SetKeys("ctrl+u")
	keymap.HalfPageDown.SetKeys("ctrl+d")

	vp.KeyMap = keymap
	vp.SoftWrap = true
	vp.Style = vp.Style.Background(lipgloss.Color("#000000")).Padding(1).PaddingLeft(2).PaddingRight(2)

	return Model{content: content, App: app, promptInput: promptInput, viewport: vp, layout: layout}
}

func (m Model) SetContent(content string) {
	m.viewport.SetContent(content)
}

func (m Model) Reload(msg tea.Msg) (Model, tea.Cmd) {
	var content []string

	for _, m := range *m.App.Messages {
		content = append(content, m.ToUIString())
	}

	m.viewport.SetContent(strings.Join(content, "\n\n"))
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) ReloadAndScrollDown(msg tea.Msg) (Model, tea.Cmd) {

	m, cmd := m.Reload(msg)
	m.viewport.GotoBottom()

	return m, cmd
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg.(type) {
	case tea.KeyPressMsg:
	case tea.WindowSizeMsg:
		// -6 for statusbar
		height := m.layout.Height - m.promptInput.Height - 6
		width := m.layout.Width

		// -1 line for padding
		height -= 1

		m.viewport.SetWidth(width)
		m.viewport.SetHeight(height)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	var v tea.View
	// v.AltScreen = true                    // use the full size of the terminal in its "alternate screen buffer"
	// v.MouseMode = tea.MouseModeCellMotion // turn on mouse support so we can track the mouse wheel

	v.SetContent(fmt.Sprintf("%s", m.viewport.View()))
	return v
}
