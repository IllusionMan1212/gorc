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

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/client"
	"github.com/illusionman1212/gorc/parser"
	"github.com/illusionman1212/gorc/ui"
)

type State struct {
	Content  string
	Reader   *bufio.Reader
	Client   *client.Client
	Viewport viewport.Model

	inputBox InputState
}

func NewMainScreen() State {
	newViewport := viewport.Model{Width: ui.MainStyle.GetWidth(), Height: ui.MainStyle.GetHeight() - InputBoxHeight}
	newViewport.HighPerformanceRendering = false
	newViewport.SetContent("")

	state := State{
		Content:  "",
		Reader:   nil,
		Client:   nil,
		Viewport: newViewport,
		inputBox: NewInputBox(),
	}
	return state
}

func (s State) Update(msg tea.Msg) (State, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case InitialReadMsg:
		return s, s.readFromServer
	case ReceivedIRCCommandMsg:
		message := parser.ParseIRCMessage(msg.Msg)
		fullMsg := fmt.Sprintf("%s %s %s %s", message.Tags, message.Source, message.Command, strings.Join(message.Parameters, " "))
		s.Content += fullMsg

		s.Viewport.SetContent(s.Content)

		// TODO: handle different commands
		switch message.Command {
		case "PING":
			s.Client.SendCommand("PONG")
			break
		}

		return s, s.readFromServer
	}

	s.Viewport, cmd = s.Viewport.Update(msg)

	return s, cmd
}

func (s State) View() string {
	mainscreen := s.Viewport.View()

	screen := lipgloss.JoinVertical(0, MessagesStyle.Render(mainscreen), s.inputBox.View())

	return screen
}
