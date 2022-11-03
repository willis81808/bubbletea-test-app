package project

import (
	"main/internal/data"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("240")).
			Align(lipgloss.Center, lipgloss.Center).
			Padding(1, 6)

	buttonDefaultStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true).
				BorderForeground(lipgloss.Color("#874BFD")).
				Margin(0, 0).
				Padding(0, 1)

	buttonHighlightedStyle = buttonDefaultStyle.Copy().
				BorderForeground(lipgloss.Color("205"))
)

type ProjectSelect struct {
	Wrapper      lipgloss.Style
	inputs       []textinput.Model
	focusIndex   int
	submitActive bool
	exitActive   bool
	state        *data.StateBag
}

type SelectedEvent struct{}

func projectSelectedCmd() tea.Cmd {
	return func() tea.Msg {
		return SelectedEvent{}
	}
}

func InitialProjectSelect() ProjectSelect {
	projectid := textinput.New()
	projectid.CharLimit = 156
	projectid.Width = 25
	projectid.TextStyle = focusedStyle
	projectid.Focus()

	return ProjectSelect{
		inputs:     []textinput.Model{projectid},
		focusIndex: 0,
		Wrapper:    lipgloss.NewStyle(),
	}
}

func renderButton(text string, active bool) string {
	var style lipgloss.Style
	if active {
		style = buttonHighlightedStyle
	} else {
		style = buttonDefaultStyle
	}
	return style.Render(text)
}

func (m ProjectSelect) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].TextStyle = blurredStyle
		}
		m.inputs[i].SetCursorMode(textinput.CursorBlink)
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m ProjectSelect) makeSelection() (tea.Model, tea.Cmd) {
	projectid := m.inputs[m.focusIndex].Value()
	m.state.ProjectId = projectid
	return m, projectSelectedCmd()
}

func (m ProjectSelect) Init() tea.Cmd {
	return tea.Batch(spinner.Tick, textinput.Blink)
}

func (m ProjectSelect) DoUpdate(msg tea.Msg, state *data.StateBag) (tea.Model, tea.Cmd) {
	m.state = state
	return m.Update(msg)
}

func (m ProjectSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.submitActive = false

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focusIndex++
			if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			}
		case "shift+tab":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}
		case "enter":
			return m.makeSelection()
		}
	case tea.MouseMsg:
		inSubmitBounds := zone.Get("submit").InBounds(msg)
		inExitBounds := zone.Get("exit").InBounds(msg)
		switch msg.Type {
		case tea.MouseLeft:
			if inSubmitBounds {
				return m.makeSelection()
			} else if inExitBounds {
				return m, tea.Quit
			}
		case tea.MouseMotion:
			m.submitActive = inSubmitBounds
			m.exitActive = inExitBounds
		}
	}

	// update inputs
	cmds = append(cmds, m.updateInputs(msg))

	return m, tea.Batch(cmds...)
}

func (m ProjectSelect) RenderView(state *data.StateBag) string {
	m.state = state
	return m.View()
}

func (m ProjectSelect) View() string {
	var inputs string
	for _, text := range m.inputs {
		var placeholderStyle lipgloss.Style
		if text.Focused() {
			placeholderStyle = noStyle
		} else {
			placeholderStyle = blurredStyle
		}
		inputGroup := lipgloss.JoinVertical(lipgloss.Center, placeholderStyle.Render("GCP Project ID"), text.View())
		inputs = lipgloss.JoinVertical(lipgloss.Center, inputs, inputGroup)
	}

	submitButton := renderButton("Submit", m.submitActive)
	exitButton := renderButton("Exit", m.exitActive)

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		zone.Mark("submit", submitButton),
		zone.Mark("exit", exitButton),
	)

	dialog := dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, inputs+"\n", buttons))
	return m.Wrapper.Render(dialog)
}
