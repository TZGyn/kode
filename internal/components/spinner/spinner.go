package spinner

import (
	"math"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/TZGyn/kode/internal/color"
)

const (
	big   = "■"
	small = "⬝"
)

func generateFrames(width int) []string {
	// fg, err := color.New("#2d8fff")
	fg := color.NewRGBFloat(0.98, 0.57, 0.24)

	bg := color.NewHex("#000000")

	loaderLength := width - 2
	frames := [][]string{}

	var currentFrame []string

	// forward direction
	for frameIndex := range width * 3 {
		currentFrame = []string{}
		for blockIndex := range width {
			if blockIndex > frameIndex {
				fg.A = 50
				currentFrame = append(currentFrame, getStyle().Foreground(lipgloss.Color(fg.ToHex(bg))).Render(small))
			} else if frameIndex-blockIndex >= loaderLength {
				fg.A = 50
				currentFrame = append(currentFrame, getStyle().Foreground(lipgloss.Color(fg.ToHex(bg))).Render(small))
			} else {
				fg.A = uint8(math.Round(math.Max(255-math.Abs(float64(frameIndex-blockIndex)*50), 50)))
				currentFrame = append(currentFrame, getStyle().Foreground(lipgloss.Color(fg.ToHex(bg))).Render(big))
			}
		}

		frames = append(frames, currentFrame)
	}

	// backward direction
	frames = append(frames, reverseFrameArray(frames)...)

	var res []string

	for _, frame := range frames {
		res = append(res, strings.Join(frame, ""))
	}

	return res
}

func reverseFrameArray(arr [][]string) [][]string {
	reversed := make([][]string, len(arr))
	for i := range arr {
		reversed[i] = reverseFrame(arr[i])
	}
	return reversed
}

func reverseFrame(arr []string) []string {
	reversed := make([]string, len(arr))
	for i := range arr {
		reversed[len(arr)-1-i] = arr[i]
	}
	return reversed
}

type Spinner struct {
	spinner spinner.Model
}

func getStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("#000000"))
}

func New() *Spinner {
	loaderSpinner := spinner.Spinner{
		Frames: generateFrames(8),

		FPS: time.Second / 18,
	}

	s := spinner.New()
	s.Spinner = loaderSpinner

	return &Spinner{
		spinner: s,
	}
}

func (s *Spinner) Init() tea.Cmd {
	var cmd []tea.Cmd

	cmd = append(cmd, s.spinner.Tick)

	return tea.Batch(cmd...)
}

func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s *Spinner) View() tea.View {
	view := tea.View{}
	view.Content = s.spinner.View() + getStyle().Render("  "+"Generating")
	return view
}
