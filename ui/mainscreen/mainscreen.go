package mainscreen

import (
	"bufio"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/parser"
)

type State struct {
	content string
	Reader  *bufio.Reader
}

func (s State) Update(msg tea.Msg) (State, tea.Cmd) {
	switch msg := msg.(type) {
	case InitialReadMsg:
		return s, s.readFromServer
	case ReceivedIRCCommandMsg:
		message := parser.ParseIRCMessage(msg.Msg)
		fullMsg := fmt.Sprintf("%s %s %s", message.Source, message.Command, strings.Join(message.Parameters, " "))
		s.content += fullMsg

		return s, s.readFromServer
	}

	return s, nil
}

func (s State) View() string {
	return s.content
}
