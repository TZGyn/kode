package commandlist

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TZGyn/kode/internal/color"
	"github.com/TZGyn/kode/internal/command"
	"github.com/charmbracelet/x/exp/charmtone"
)

const listHeight = 14
const defaultWidth = 80

type CommandList struct {
	list   list.Model
	choice string

	Width, Height       int
	minWidth, minHeight int
}

type item string

func (i item) FilterValue() string { return string(i) }

func New() CommandList {

	items := []list.Item{}

	for _, command := range command.Commands {
		items = append(items, item(command.CommandID))
	}

	l := list.New(items, list.DefaultDelegate{}, defaultWidth, listHeight)

	return CommandList{
		list: l,
	}
}

func (c CommandList) Init() tea.Cmd {
	return nil
}

func (c CommandList) Update(msg tea.Msg) (CommandList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		c.Width = msg.Width
		c.list.SetWidth(msg.Width)
		return c, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := c.list.SelectedItem().(item)
			if ok {
				c.choice = string(i)
			}
			return c, nil
		case "esc":
			return c, nil
		}
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c CommandList) View() tea.View {
	return tea.NewView(newCard(true, c.list.View(), c.Width))
}

func (c CommandList) List() list.Model {
	return c.list
}

func newCard(darkMode bool, text string, width int) string {
	lightDark := lipgloss.LightDark(darkMode)

	fg := color.NewRGBFloat(0.98, 0.57, 0.24)

	fg.A = color.MaxColorValue

	bg := color.NewHex("#000000")

	content := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		// BorderForegroundBlend(
		// 	charmtone.Cherry,
		// 	charmtone.Charple,
		// 	charmtone.Guac,
		// 	charmtone.Charple,
		// 	charmtone.Sriracha,
		// ).
		BorderForeground(fg.ToColor(bg)).
		Foreground(lightDark(charmtone.Iron, charmtone.Butter)).
		Background(lipgloss.Color("#000000")).
		// Height(9).
		Width(width).
		PaddingTop(3).
		Align(lipgloss.Center)

	return content.
		Render(text)
}
