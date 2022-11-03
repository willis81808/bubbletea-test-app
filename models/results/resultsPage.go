package results

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ResultsPage struct {
	projectId string
	region    string
}

func InitialResultsPage() ResultsPage {
	return ResultsPage{}
}

func (m *ResultsPage) SetData(projectId string, region string) ResultsPage {
	m.projectId = projectId
	m.region = region
	return *m
}

func (m ResultsPage) Init() tea.Cmd {
	return nil
}

func (m ResultsPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ResultsPage) View() string {
	str := fmt.Sprintf("\n\nProject ID: %s\nRegion: %s\n\n", m.projectId, m.region)
	return str
}
