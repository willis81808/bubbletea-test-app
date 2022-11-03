package main

import (
	"fmt"
	"os"

	"main/internal/data"
	"main/internal/views/project"
	"main/internal/views/region"
	"main/internal/views/results"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	regionView sessionState = iota
	projectView
	resultsView
)

var (
	FullscreenStyle = lipgloss.NewStyle().Padding(0).Margin(0).Align(lipgloss.Center, lipgloss.Center)
)

type sessionState uint

type errMsg error

type EntryModel struct {
	state         sessionState
	stateData     data.StateBag
	regionSelect  region.RegionSelect
	projectSelect project.ProjectSelect
	resultsView   results.ResultsPage
	quitting      bool
	err           error
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func initialModel() EntryModel {
	model := EntryModel{state: projectView}

	model.projectSelect = project.InitialProjectSelect()
	model.projectSelect.Wrapper = FullscreenStyle

	model.regionSelect = region.InitialRegionSelect()

	model.resultsView = results.InitialResultsPage()
	model.resultsView.Wrapper = FullscreenStyle

	model.stateData = data.NewStateBag()
	return model
}

func (m EntryModel) Init() tea.Cmd {
	return nil
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			m.state = projectView
		}
	case errMsg:
		m.err = msg
		return m, nil
	case tea.WindowSizeMsg:
		FullscreenStyle.Width(msg.Width)
		FullscreenStyle.Height(msg.Height)

		newRegion, newCmd := m.regionSelect.DoUpdate(msg, &m.stateData)
		m.regionSelect = newRegion.(region.RegionSelect)
		cmds = append(cmds, newCmd)

		newProject, newCmd := m.projectSelect.DoUpdate(msg, &m.stateData)
		m.projectSelect = newProject.(project.ProjectSelect)
		cmds = append(cmds, newCmd)

		newResults, newCmd := m.resultsView.DoUpdate(msg, &m.stateData)
		m.resultsView = newResults.(results.ResultsPage)
		cmds = append(cmds, newCmd)
	case project.SelectedEvent:
		m.state = regionView
	case region.SelectedEvent:
		m.state = resultsView
	}

	switch m.state {
	case regionView:
		newRegion, newCmd := m.regionSelect.DoUpdate(msg, &m.stateData)
		m.regionSelect = newRegion.(region.RegionSelect)
		cmds = append(cmds, newCmd)
	case projectView:
		newProject, newCmd := m.projectSelect.DoUpdate(msg, &m.stateData)
		m.projectSelect = newProject.(project.ProjectSelect)
		cmds = append(cmds, newCmd)
	case resultsView:
		newResults, newCmd := m.resultsView.DoUpdate(msg, &m.stateData)
		m.resultsView = newResults.(results.ResultsPage)
		cmds = append(cmds, newCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m EntryModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var str string

	switch m.state {
	case projectView:
		str += m.projectSelect.RenderView(&m.stateData)
	case regionView:
		str += m.regionSelect.RenderView(&m.stateData)
	case resultsView:
		str += m.resultsView.RenderView(&m.stateData)
	}

	if m.quitting {
		return "\n"
	}

	return zone.Scan(str)
}

func main() {
	zone.NewGlobal()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())
	zone.SetEnabled(true)
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.ExitAltScreen()
}
