package project

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
)

type ProjectSelect struct {
	inputs     []textinput.Model
	focusIndex int
}

type ProjectSelected struct {
	projectid string
}

func (p *ProjectSelected) GetProjectId() string {
	return p.projectid
}

func projectSelectedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ProjectSelected{projectid: id}
	}
}

func InitialProjectSelect() ProjectSelect {
	projectid := textinput.New()
	projectid.Placeholder = "GCP Project ID"
	projectid.CharLimit = 156
	projectid.Width = 35
	projectid.TextStyle = focusedStyle
	projectid.Focus()

	region := textinput.New()
	region.Placeholder = "GCP Region"
	region.CharLimit = 156
	region.Width = 35

	return ProjectSelect{
		inputs:     []textinput.Model{projectid},
		focusIndex: 0,
	}
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

func (m ProjectSelect) Init() tea.Cmd {
	return tea.Batch(spinner.Tick, textinput.Blink)
}

func (m ProjectSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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
			projectid := m.inputs[0].Value()
			return m, projectSelectedCmd(projectid)
		}
	}

	// update inputs
	cmds = append(cmds, m.updateInputs(msg))

	return m, tea.Batch(cmds...)
}

func (m ProjectSelect) View() string {
	var str string
	for _, text := range m.inputs {
		var placeholderStyle lipgloss.Style
		if text.Focused() {
			placeholderStyle = noStyle
		} else {
			placeholderStyle = blurredStyle
		}
		str += fmt.Sprintf("\n%s\n%s\n", placeholderStyle.Render(text.Placeholder), text.View())
	}
	str += "\n"
	return str
}
