package mainscreen

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type ReceivedIRCCommandMsg struct {
	Msg string
}

type InitialReadMsg struct{}

func InitialRead() tea.Msg {
	return InitialReadMsg{}
}

func (s State) readFromServer() tea.Msg {
	msg, err := s.Reader.ReadString('\n')
	if err != nil {
		log.Print(err)
	}

	return ReceivedIRCCommandMsg{Msg: msg}
}
