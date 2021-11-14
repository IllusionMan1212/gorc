package login

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ConnectingMsg struct{}

func connect() tea.Msg {
	return ConnectingMsg{}
}
