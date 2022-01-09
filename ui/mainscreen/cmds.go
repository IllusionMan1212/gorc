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
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/parser"
)

func (s *State) ReadLoop() {
	for {
		select {
		case <-s.Client.Quit:
			return
		default:
			msg, err := s.ConnReader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}

			msg = strings.Replace(msg, "\r\n", "", 1)
			if ircMessage, valid := parser.ParseIRCMessage(msg); valid {
				s.HandleCommand(ircMessage)
			}
		}
	}
}

type SendPrivMsg struct {
	Msg string
}

func (s InputState) SendingPrivMsg(msg string) tea.Cmd {
	return func() tea.Msg {
		return SendPrivMsg{Msg: msg}
	}
}
