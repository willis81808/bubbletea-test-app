package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"

	//"github.com/charmbracelet/bubbles/table"
	"main/models/project"
	"main/models/region"
	"main/models/results"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	regionView sessionState = iota
	projectView
	resultsView
)

type sessionState uint

type errMsg error

type stateBag struct {
	projectId string
	region    string
}

type EntryModel struct {
	state         sessionState
	data          stateBag
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
	model.regionSelect = region.InitialRegionSelect()
	model.resultsView = results.InitialResultsPage()
	model.data = stateBag{}
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
	case project.ProjectSelected:
		m.data.projectId = msg.GetProjectId()
		m.state = regionView
	case region.RegionSelected:
		m.data.region = msg.GetRegion()
		m.state = resultsView
	}

	switch m.state {
	case regionView:
		newRegion, newCmd := m.regionSelect.Update(msg)
		m.regionSelect = newRegion.(region.RegionSelect)
		cmds = append(cmds, newCmd)
	case projectView:
		newProject, newCmd := m.projectSelect.Update(msg)
		m.projectSelect = newProject.(project.ProjectSelect)
		cmds = append(cmds, newCmd)
	case resultsView:
		newResults, newCmd := m.resultsView.Update(msg)
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
		str += m.projectSelect.View()
	case regionView:
		str += m.regionSelect.View()
	case resultsView:
		m.resultsView.SetData(m.data.projectId, m.data.region)
		str += m.resultsView.View()
	}

	if m.quitting {
		return "\n"
	}

	return str //+ fmt.Sprintf("\n\nProject: %s\nRegion: %s", m.data.projectId, m.data.region)
}

func main() {
	p := tea.NewProgram(initialModel())
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.ExitAltScreen()
}
