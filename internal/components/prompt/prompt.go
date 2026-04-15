package prompt

import (
	"math"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type PromptComponent struct {
	textarea textarea.Model

	Width, Height       int
	minWidth, minHeight int
}

const MinHeight = 3

func NewPrompt() PromptComponent {

	initialWidth := 0
	initialHeight := 3

	ta := textarea.New()
	// ta.SetWidth(initialWidth)
	ta.SetHeight(initialHeight)
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.MaxWidth = 100

	keybinds := textarea.DefaultKeyMap()
	keybinds.InsertNewline.SetEnabled(false)
	ta.KeyMap = keybinds

	b := lipgloss.ThickBorder()

	style := ta.Styles()

	style.Focused.Base = lipgloss.NewStyle().Background(lipgloss.Color("#131313")).Padding(1).PaddingLeft(2).PaddingRight(2).Border(b, false, false, false, true).BorderLeftForeground(lipgloss.Color("#2d8fff"))

	style.Focused.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color("#131313"))
	style.Focused.Text = lipgloss.NewStyle().Background(lipgloss.Color("#131313"))
	style.Focused.EndOfBuffer = lipgloss.NewStyle().Background(lipgloss.Color("#131313"))
	style.Focused.Placeholder = lipgloss.NewStyle().Background(lipgloss.Color("#131313"))

	ta.Prompt = ""

	// ta.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#2d8fff")).Render("┃ ")

	ta.SetStyles(style)

	prompt := &PromptComponent{
		Width:     initialWidth,
		Height:    initialHeight,
		minWidth:  0,
		minHeight: initialHeight,
		textarea:  ta,
	}

	return *prompt
}

func (c *PromptComponent) Init() tea.Cmd {
	return nil
}

func (c *PromptComponent) View() tea.View {
	return tea.NewView(c.textarea.View())
}

func (c *PromptComponent) Update(msg tea.Msg) (*PromptComponent, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			c.textarea.SetValue("")
		}
		c.Resize()
	case tea.WindowSizeMsg:
		// msg.Height -= 2 // Make space for the status bar
		c.Width = msg.Width - 2
		c.textarea.SetWidth(c.Width)
		c.Resize()
	}

	var cmd tea.Cmd
	c.textarea, cmd = c.textarea.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c *PromptComponent) Value() string {
	return c.textarea.Value()
}

func (c *PromptComponent) Resize() {
	c.Height = int(math.Max(float64(c.Lines()+1), float64(c.minHeight)))
	c.textarea.SetHeight(c.Height)
}

func (c *PromptComponent) LineCount() int {
	return c.textarea.LineCount()
}

func (c *PromptComponent) InputWidth() int {
	return c.textarea.Width()
}

func (c *PromptComponent) Lines() int {
	lines := strings.Split(c.textarea.Value(), "\n")

	num_of_lines := 0

	for _, line := range lines {
		divide := float64(len(line) / c.textarea.Width())
		num_of_lines += int(math.Ceil(divide))
	}

	return num_of_lines + 1
}

func (c *PromptComponent) SetValue(value string) {
	c.textarea.SetValue(value)
	c.Resize()
}
