package compound

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
	"github.com/willis81808/bubbletea-test-app/internal/data"
	"github.com/willis81808/bubbletea-test-app/internal/utils"
)

const (
	leftAddr   = "left"
	middleAddr = "middle"
	rightAddr  = "right"
	lowerAddr  = "lower"
)

type stringer string

func (s stringer) String() string {
	return string(s)
}

// satisfy the tea.Model interface
func (s stringer) Init() tea.Cmd                           { return nil }
func (s stringer) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s stringer) View() string                            { return s.String() }

type CompoundView struct {
	tui   boxer.Boxer
	state *data.StateBag
}

func InitialCompoundView() CompoundView {
	left := stringer(leftAddr)
	middle := stringer(middleAddr)
	right := stringer(rightAddr)

	lower := stringer(fmt.Sprintf("%s: user ctrl+c to quit", lowerAddr))

	m := CompoundView{tui: boxer.Boxer{}}

	m.tui.LayoutTree = boxer.Node{
		VerticalStacked: true,
		SizeFunc: func(_ boxer.Node, widthOrHeight int) []int {
			return []int{
				widthOrHeight - 1,
				1,
			}
		},
		Children: []boxer.Node{
			{
				Children: []boxer.Node{
					m.tui.CreateLeaf(leftAddr, left),
					m.tui.CreateLeaf(middleAddr, middle),
					m.tui.CreateLeaf(rightAddr, right),
				},
			},
			m.tui.CreateLeaf(lowerAddr, lower),
		},
	}

	return m
}

func (c CompoundView) Init() tea.Cmd {
	return nil
}

func (m CompoundView) DoUpdate(msg tea.Msg, state *data.StateBag) (utils.SubModel, tea.Cmd) {
	m.state = state
	model, cmd := m.Update(msg)
	return model.(CompoundView), cmd
}

func (c CompoundView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.tui.UpdateSize(msg)
	}
	return c, nil
}

func (m CompoundView) RenderView(state *data.StateBag) string {
	m.state = state
	return m.View()
}

func (c CompoundView) View() string {
	return c.tui.View()
}
