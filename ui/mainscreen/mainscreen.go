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
	Style    lipgloss.Style

	InputBox InputState
}

func NewMainScreen() State {
	newViewport := viewport.Model{}
	newViewport.HighPerformanceRendering = false
	newViewport.SetContent("")

	state := State{
		Content:  "",
		Reader:   nil,
		Client:   nil,
		Viewport: newViewport,
		Style:    MessagesStyle.Copy(),
		InputBox: NewInputBox(),
	}
	return state
}

func (s State) Update(msg tea.Msg) (State, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

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
	cmds = append(cmds, cmd)

	s.InputBox, cmd = s.InputBox.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s *State) SetSize(width, height int) {
	s.InputBox.SetSize(width)
	s.Viewport.Width = width - s.Style.GetVerticalFrameSize()
	// copying the existing style is important here otherwise we'll end up with artifacting
	s.Style = s.Style.Copy().Width(width - s.Style.GetVerticalFrameSize())
	s.Viewport.Height = height - s.InputBox.Style.GetHorizontalFrameSize() - s.Style.GetHorizontalFrameSize()
	s.Style = s.Style.Copy().Height(height - s.InputBox.Style.GetHorizontalFrameSize() - s.Style.GetHorizontalFrameSize())
	// we need to re-set the content because words wrap differently on different sizes
	s.Viewport.SetContent(s.Content)
}

func (s State) View() string {
	screen := lipgloss.JoinVertical(0, s.Style.Render(s.Viewport.View()), s.InputBox.View())

	return ui.MainStyle.Render(screen)
}
