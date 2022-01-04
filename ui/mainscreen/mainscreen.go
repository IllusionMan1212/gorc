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

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/client"
	"github.com/illusionman1212/gorc/parser"
	"github.com/illusionman1212/gorc/ui"
)

type Window int

const (
	Viewport = iota
	InputBox
)

type State struct {
	Content    string
	Reader     *bufio.Reader
	Client     *client.Client
	Viewport   viewport.Model
	Style      lipgloss.Style
	FocusIndex Window

	// a channel could be either a server channel
	// or a username of a user
	CurrentChannel string

	InputBox InputState
}

func NewMainScreen(client *client.Client) State {
	newViewport := viewport.Model{
		HighPerformanceRendering: false,
		Wrap:                     true,
	}

	state := State{
		Content:    "",
		Reader:     nil,
		Client:     client,
		Viewport:   newViewport,
		Style:      MessagesStyle.Copy(),
		FocusIndex: InputBox,
		InputBox:   NewInputBox(client),
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
		s.HandleCommand(message)

		return s, s.readFromServer
	case SendPrivMsg:
		s.Content += msg.Msg + client.CRLF
		s.Client.SendCommand("PRIVMSG", s.CurrentChannel, msg.Msg)

		s.Viewport.SetContent(s.Content)

		return s, s.readFromServer
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "tab", "shift+tab":
			if key == "tab" {
				s.FocusIndex++
			} else {
				s.FocusIndex--
			}

			if s.FocusIndex > 1 {
				s.FocusIndex = 0
			} else if s.FocusIndex < 0 {
				s.FocusIndex = 1
			}

			switch s.FocusIndex {
			case Viewport:
				s.Style = s.Style.Copy().BorderForeground(lipgloss.Color("105"))
				s.InputBox.Style = s.InputBox.Style.Copy().BorderForeground(lipgloss.Color("#EEE"))

				s.InputBox.Input.Blur()
			case InputBox:
				s.Style = s.Style.Copy().BorderForeground(lipgloss.Color("#EEE"))
				s.InputBox.Style = s.InputBox.Style.Copy().BorderForeground(lipgloss.Color("105"))

				s.InputBox, cmd = s.InputBox.Update(msg)
				cmds = append(cmds, cmd)
				cmds = append(cmds, textinput.Blink)
				s.InputBox.Input.Focus()
			}

			return s, tea.Batch(cmds...)
		}
	}

	switch s.FocusIndex {
	case Viewport:
		s.Viewport, cmd = s.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	case InputBox:
		s.InputBox, cmd = s.InputBox.Update(msg)
		cmds = append(cmds, cmd)
	}

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

func (s *State) HandleCommand(msg parser.IRCMessage) {
	fullMsg := fmt.Sprintf("%s %s %s %s", msg.Tags, msg.Source, msg.Command, strings.Join(msg.Parameters, " "))
	s.Content += fullMsg

	s.Viewport.SetContent(s.Content)

	// TODO: handle different commands
	switch msg.Command {
	case "PING":
		s.Client.SendCommand("PONG", msg.Parameters[0])
		break
	}
}

func (s State) View() string {
	screen := lipgloss.JoinVertical(0, s.Style.Render(s.Viewport.View()), s.InputBox.View())

	return ui.MainStyle.Render(screen)
}
