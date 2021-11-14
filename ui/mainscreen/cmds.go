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
