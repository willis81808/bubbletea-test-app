package main

import (
	"fmt"
	"os"

	"github.com/willis81808/bubbletea-test-app/internal/data"
	"github.com/willis81808/bubbletea-test-app/internal/utils"
	"github.com/willis81808/bubbletea-test-app/internal/views/compound"
	"github.com/willis81808/bubbletea-test-app/internal/views/project"
	"github.com/willis81808/bubbletea-test-app/internal/views/region"
	"github.com/willis81808/bubbletea-test-app/internal/views/results"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	regionView sessionState = iota
	projectView
	resultsView
	compoundView
)

var (
	FullscreenStyle = lipgloss.NewStyle().Padding(0).Margin(0).Align(lipgloss.Center, lipgloss.Center)
)

type sessionState uint

type errMsg error

type EntryModel struct {
	state     sessionState
	history   *utils.Stack[sessionState]
	stateData data.StateBag
	models    map[sessionState]utils.SubModel
	err       error
}

func initialModel() EntryModel {
	model := EntryModel{state: projectView}

	projectModel := project.InitialProjectSelect()
	projectModel.Wrapper = FullscreenStyle

	regionModel := region.InitialRegionSelect()

	resultsModel := results.InitialResultsPage()
	resultsModel.Wrapper = FullscreenStyle

	compoundModel := compound.InitialCompoundView()

	model.stateData = data.NewStateBag()

	model.history = utils.NewStack[sessionState]()

	model.models = make(map[sessionState]utils.SubModel)
	model.models[projectView] = projectModel
	model.models[regionView] = regionModel
	model.models[resultsView] = resultsModel
	model.models[compoundView] = compoundModel

	return model
}

func (m *EntryModel) lastState() {
	if m.history.Length() == 0 {
		return
	}
	m.state, _ = m.history.Pop()
}

func (m *EntryModel) nextState(newState sessionState) {
	if m.state == newState {
		return
	}
	m.history.Push(m.state)
	m.state = newState
}

func (m EntryModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, submodel := range m.models {
		cmds = append(cmds, submodel.Init())
	}

	return tea.Batch(cmds...)
}

func (m EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		case "esc":
			m.lastState()
			return m, nil
		case "\\":
			m.nextState(compoundView)
		}
	case tea.WindowSizeMsg:
		FullscreenStyle.Width(msg.Width)
		FullscreenStyle.Height(msg.Height)

		// ensure even non-active models are alerted to the layout change
		for key, value := range m.models {
			if key != m.state {
				newModel, newCmd := value.DoUpdate(msg, &m.stateData)
				m.models[key] = newModel
				cmds = append(cmds, newCmd)
			}
		}
	case project.SelectedEvent:
		m.nextState(regionView)
	case region.SelectedEvent:
		m.nextState(resultsView)
	case errMsg:
		m.err = msg
		return m, nil
	}

	// forward command to the currently active model
	activeModel := m.models[m.state]
	newModel, newCmd := activeModel.DoUpdate(msg, &m.stateData)
	m.models[m.state] = newModel
	cmds = append(cmds, newCmd)

	return m, tea.Batch(cmds...)
}

func (m EntryModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	// render currently active model
	str := m.models[m.state].RenderView(&m.stateData)

	// enable mouse zones for rendered model
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
