package region

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RegionSelect struct {
	regionList   list.Model
	regionChoice string
}

type item struct {
	title, desc string
}

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	regions  = []string{"asia-east1", "asia-east2", "asia-northeast1", "asia-northeast2", "asia-northeast3", "asia-south1", "asia-south2", "asia-southeast1", "asia-southeast2", "australia-southeast1", "australia-southeast2", "europe-central2", "europe-north1", "europe-west1", "europe-west2", "europe-west3", "europe-west4", "europe-west6", "northamerica-northeast1", "northamerica-northeast2", "southamerica-east1", "southamerica-west1", "us-central1", "us-east1", "us-east4", "us-west1", "us-west2", "us-west3", "us-west4"}
)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type RegionSelected struct {
	region string
}

func (r *RegionSelected) GetRegion() string {
	return r.region
}

func regionSelectedCmd(region string) tea.Cmd {
	return func() tea.Msg {
		return RegionSelected{region}
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

func (m RegionSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.regionList.FilterState() == list.Filtering {
				break
			}
			region := m.regionList.SelectedItem().(item).Title()
			return m, regionSelectedCmd(region)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.regionList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.regionList, cmd = m.regionList.Update(msg)
	return m, cmd
}

func (m RegionSelect) View() string {
	return docStyle.Render(m.regionList.View())
}
