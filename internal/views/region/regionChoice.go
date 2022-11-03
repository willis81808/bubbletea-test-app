package region

import (
	"main/internal/data"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RegionSelect struct {
	regionList   list.Model
	regionChoice string
	state        *data.StateBag
}

type item struct {
	title, desc string
}

var (
	docStyle = lipgloss.NewStyle().AlignHorizontal(lipgloss.Left)
	regions  = []string{"asia-east1", "asia-east2", "asia-northeast1", "asia-northeast2", "asia-northeast3", "asia-south1", "asia-south2", "asia-southeast1", "asia-southeast2", "australia-southeast1", "australia-southeast2", "europe-central2", "europe-north1", "europe-west1", "europe-west2", "europe-west3", "europe-west4", "europe-west6", "northamerica-northeast1", "northamerica-northeast2", "southamerica-east1", "southamerica-west1", "us-central1", "us-east1", "us-east4", "us-west1", "us-west2", "us-west3", "us-west4"}
)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type SelectedEvent struct{}

func regionSelectedCmd(region string) tea.Cmd {
	return func() tea.Msg {
		return SelectedEvent{}
	}
}

func InitialRegionSelect() RegionSelect {
	regionItems := make([]list.Item, len(regions))
	for i, region := range regions {
		regionItems[i] = item{title: region}
	}

	list := list.New(regionItems, list.NewDefaultDelegate(), 40, 30)
	list.Title = "Regions"

	return RegionSelect{
		regionList: list,
	}
}

func (m RegionSelect) Init() tea.Cmd {
	return nil
}

func (m RegionSelect) DoUpdate(msg tea.Msg, state *data.StateBag) (tea.Model, tea.Cmd) {
	m.state = state
	return m.Update(msg)
}

func (m RegionSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.regionList.FilterState() == list.Filtering {
				break
			}
			region := m.regionList.SelectedItem().(item).Title()
			m.state.Region = region
			return m, regionSelectedCmd(region)
		}
	case tea.WindowSizeMsg:
		l, v := docStyle.GetFrameSize()
		m.regionList.SetSize(msg.Width-l, msg.Height-v)
	}

	var cmd tea.Cmd
	m.regionList, cmd = m.regionList.Update(msg)
	return m, cmd
}

func (m RegionSelect) RenderView(state *data.StateBag) string {
	m.state = state
	return m.View()
}

func (m RegionSelect) View() string {
	return docStyle.Render(m.regionList.View())
}
