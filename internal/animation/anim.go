package animation

import (
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	charCyclingFPS  = time.Second / 22
	colorCycleFPS   = time.Second / 5
	maxCyclingChars = 120
)

var charRunes = []rune("0123456789abcdefABCDEF~!@#$£€%^&*()+=_")

type charState int

const (
	charInitialState charState = iota
	charCyclingState
	charEndOfLifeState
)

// cyclingChar is a single animated character.
type cyclingChar struct {
	finalValue   rune // if < 0 cycle forever
	currentValue rune
	initialDelay time.Duration
	lifetime     time.Duration
}

func (c cyclingChar) randomRune() rune {
	return (charRunes)[rand.Intn(len(charRunes))] //nolint:gosec
}

func (c cyclingChar) state(start time.Time) charState {
	now := time.Now()
	if now.Before(start.Add(c.initialDelay)) {
		return charInitialState
	}
	if c.finalValue > 0 && now.After(start.Add(c.initialDelay)) {
		return charEndOfLifeState
	}
	return charCyclingState
}

type stepCharsMsg struct{}

func stepChars() tea.Cmd {
	return tea.Tick(charCyclingFPS, func(time.Time) tea.Msg {
		return stepCharsMsg{}
	})
}

type colorCycleMsg struct{}

func cycleColors() tea.Cmd {
	return tea.Tick(colorCycleFPS, func(time.Time) tea.Msg {
		return colorCycleMsg{}
	})
}

// anim is the model that manages the animation that displays while the
// output is being generated.
type anim struct {
	start           time.Time
	cyclingChars    []cyclingChar
	labelChars      []cyclingChar
	ramp            []lipgloss.Style
	label           []rune
	ellipsis        spinner.Model
	ellipsisStarted bool
}

func NewAnim(label string) anim {
	// #nosec G115
	n := int(12)
	// if n > maxCyclingChars {
	// 	n = maxCyclingChars
	// }

	gap := " "
	if n == 0 {
		gap = ""
	}

	c := anim{
		start:    time.Now(),
		label:    []rune(gap + label),
		ellipsis: spinner.New(spinner.WithSpinner(spinner.Ellipsis)),
	}

	// If we're in truecolor mode (and there are enough cycling characters)
	// color the cycling characters with a gradient ramp.
	const minRampSize = 3
	if n >= minRampSize {
		// Note: double capacity for color cycling as we'll need to reverse and
		// append the ramp for seamless transitions.
		c.ramp = make([]lipgloss.Style, n, n*2) //nolint:mnd
		ramp := makeGradientRamp(n)
		for i, color := range ramp {
			c.ramp[i] = lipgloss.NewStyle().Foreground(color)
		}
		c.ramp = append(c.ramp, reverse(c.ramp)...) // reverse and append for color cycling
	}

	makeDelay := func(a int32, b time.Duration) time.Duration {
		return time.Duration(rand.Int31n(a)) * (time.Millisecond * b) //nolint:gosec
	}

	makeInitialDelay := func() time.Duration {
		return makeDelay(8, 60) //nolint:mnd
	}

	// Initial characters that cycle forever.
	c.cyclingChars = make([]cyclingChar, n)

	for i := 0; i < n; i++ {
		c.cyclingChars[i] = cyclingChar{
			finalValue:   -1, // cycle forever
			initialDelay: makeInitialDelay(),
		}
	}

	// Label text that only cycles for a little while.
	c.labelChars = make([]cyclingChar, len(c.label))

	for i, r := range c.label {
		c.labelChars[i] = cyclingChar{
			finalValue:   r,
			initialDelay: makeInitialDelay(),
			lifetime:     makeDelay(5, 180), //nolint:mnd
		}
	}

	return c
}

// Init initializes the animation.
func (anim) Init() tea.Cmd {
	return tea.Batch(stepChars(), cycleColors())
}

// Update handles messages.
func (a anim) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case stepCharsMsg:
		a.updateChars(&a.cyclingChars)
		a.updateChars(&a.labelChars)

		if !a.ellipsisStarted {
			var eol int
			for _, c := range a.labelChars {
				if c.state(a.start) == charEndOfLifeState {
					eol++
				}
			}
			if eol == len(a.label) {
				// If our entire label has reached end of life, start the
				// ellipsis "spinner" after a short pause.
				a.ellipsisStarted = true
				cmd = tea.Tick(time.Millisecond*220, func(time.Time) tea.Msg { //nolint:mnd
					return a.ellipsis.Tick()
				})
			}
		}

		return a, tea.Batch(stepChars(), cmd)
	case colorCycleMsg:
		const minColorCycleSize = 2
		if len(a.ramp) < minColorCycleSize {
			return a, nil
		}
		a.ramp = append(a.ramp[1:], a.ramp[0])
		return a, cycleColors()
	case spinner.TickMsg:
		var cmd tea.Cmd
		a.ellipsis, cmd = a.ellipsis.Update(msg)
		return a, cmd
	default:
		return a, nil
	}
}

func (a *anim) updateChars(chars *[]cyclingChar) {
	for i, c := range *chars {
		switch c.state(a.start) {
		case charInitialState:
			(*chars)[i].currentValue = '.'
		case charCyclingState:
			(*chars)[i].currentValue = c.randomRune()
		case charEndOfLifeState:
			(*chars)[i].currentValue = c.finalValue
		}
	}
}

// View renders the animation.
func (a anim) View() string {
	var b strings.Builder

	for i, c := range a.cyclingChars {
		if len(a.ramp) > i {
			b.WriteString(a.ramp[i].Render(string(c.currentValue)))
			continue
		}
		b.WriteRune(c.currentValue)
	}

	for _, c := range a.labelChars {
		b.WriteRune(c.currentValue)
	}

	return b.String() + a.ellipsis.View()
}

func makeGradientRamp(length int) []lipgloss.Color {
	const startColor = "#93c5fd"

	const endColor = "#2563eb"

	var (
		c        = make([]lipgloss.Color, length)
		start, _ = colorful.Hex(startColor)
		end, _   = colorful.Hex(endColor)
	)
	for i := range length {
		step := start.BlendLuv(end, float64(i)/float64(length))
		c[i] = lipgloss.Color(step.Hex())
	}
	return c
}

func reverse[T any](in []T) []T {
	out := make([]T, len(in))
	copy(out, in[:])
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}
