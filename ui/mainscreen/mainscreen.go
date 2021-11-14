// gorc project
// Copyright (C) 2021 IllusionMan1212
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see https://www.gnu.org/licenses.

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
