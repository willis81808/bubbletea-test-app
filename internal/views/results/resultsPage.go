package results

import (
	"fmt"

	"github.com/willis81808/bubbletea-test-app/internal/data"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResultsPage struct {
	Wrapper   lipgloss.Style
	projectId string
	region    string
	state     *data.StateBag
}

func InitialResultsPage() ResultsPage {
	return ResultsPage{Wrapper: lipgloss.NewStyle()}
}

func (m ResultsPage) Init() tea.Cmd {
	return nil
}

func (m ResultsPage) DoUpdate(msg tea.Msg, state *data.StateBag) (tea.Model, tea.Cmd) {
	m.state = state
	return m.Update(msg)
}

func (m ResultsPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ResultsPage) RenderView(state *data.StateBag) string {
	m.state = state
	return m.View()
}

func (m ResultsPage) View() string {
	str := fmt.Sprintf("\n\nProject ID: %s\nRegion: %s\n\n", m.state.ProjectId, m.state.Region)
	return m.Wrapper.Render(str)
}
