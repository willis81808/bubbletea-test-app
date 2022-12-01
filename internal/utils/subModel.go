package utils

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/willis81808/bubbletea-test-app/internal/data"
)

type SubModel interface {
	tea.Model

	// Init is the first function that will be called. It returns an optional
	// initial command. To not perform an initial command return nil.
	Init() tea.Cmd

	// Update is called when a message is received. Use it to inspect messages
	// and, in response, update the model and/or send a command.
	DoUpdate(tea.Msg, *data.StateBag) (SubModel, tea.Cmd)

	// RenderView renders the program's UI, which is just a string. The view is
	// rendered after every Update.
	RenderView(*data.StateBag) string
}
