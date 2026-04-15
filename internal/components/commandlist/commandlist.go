package commandlist

import (
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TZGyn/kode/internal/command"
)

type item struct {
	cmd      command.Command
	selected bool
}

func (i item) FilterValue() string { return i.cmd.CommandID }

type CommandList struct {
	list     list.Model
	items    []item
	choice   string
	visible  bool
	search   string
	selected int

	Width, Height int
}

func New() CommandList {
	return CommandList{
		list:     list.New(nil, list.NewDefaultDelegate(), 0, 0),
		items:    []item{},
		visible:  false,
		selected: 0,
	}
}

func (c *CommandList) Show(search string) {
	c.search = strings.TrimPrefix(search, "/")
	c.visible = true
	c.selected = 0
	c.updateItems()
}

func (c *CommandList) Hide() {
	c.visible = false
	c.search = ""
	c.choice = ""
	c.selected = 0
}

func (c *CommandList) Filter(search string) {
	if !c.visible {
		return
	}
	oldSelected := ""
	if len(c.items) > 0 && c.selected < len(c.items) {
		oldSelected = c.items[c.selected].cmd.CommandID
	}

	c.search = strings.TrimPrefix(search, "/")
	c.selected = 0
	c.updateItems()

	for i, item := range c.items {
		if item.cmd.CommandID == oldSelected {
			c.selected = i
			break
		}
	}
}

func (c *CommandList) updateItems() {
	cmds := command.FilterCommands(c.search)
	c.items = make([]item, len(cmds))
	for i, cmd := range cmds {
		c.items[i] = item{cmd: cmd, selected: i == c.selected}
	}
}

func (c CommandList) IsVisible() bool {
	return c.visible
}

func (c CommandList) Init() tea.Cmd {
	return nil
}

func (c *CommandList) Update(msg tea.Msg) (CommandList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.Width = msg.Width
		return *c, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if len(c.items) > 0 && c.selected < len(c.items) {
				c.choice = c.items[c.selected].cmd.CommandID
			}
			return *c, nil
		case "esc":
			c.Hide()
			return *c, nil
		}
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)
	return *c, cmd
}

func (c *CommandList) updateSelection() {
	for i := range c.items {
		c.items[i].selected = (i == c.selected)
	}
}

func (c CommandList) View() string {
	if !c.visible || len(c.items) == 0 {
		return ""
	}

	var sb strings.Builder

	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Render("Commands")

	sb.WriteString(header)
	sb.WriteString("\n\n")

	for i, item := range c.items {
		prefix := "  "
		cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))

		if item.selected {
			prefix = "▸ "
			cmdStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#2d8fff")).Bold(true)
			descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
		}

		if i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(prefix)
		sb.WriteString(cmdStyle.Render(item.cmd.CommandID))
		sb.WriteString("  ")
		sb.WriteString(descStyle.Render(item.cmd.Description))
	}

	borderFg := lipgloss.Color("#333333")
	contentBg := lipgloss.Color("#1a1a1a")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderFg).
		Background(contentBg).
		Width(c.Width).
		Padding(1, 2).
		Render(sb.String())

	return box
}

func (c CommandList) List() list.Model {
	return c.list
}

func (c CommandList) Choice() string {
	return c.choice
}

func (c CommandList) ItemCount() int {
	return len(c.items)
}
