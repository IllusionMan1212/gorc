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
	"github.com/illusionman1212/gorc/commands"
	"github.com/illusionman1212/gorc/parser"
	"github.com/illusionman1212/gorc/ui"
)

type Window int

const (
	Viewport = iota
	InputBox
)

type State struct {
	ConnReader *bufio.Reader
	Client     *client.Client
	Viewport   *viewport.Model
	Style      lipgloss.Style
	FocusIndex Window

	InputBox  InputState
	SidePanel *SidePanelState
}

func NewMainScreen(client *client.Client) State {
	newViewport := &viewport.Model{
		HighPerformanceRendering: false,
		Wrap:                     true,
	}

	return State{
		ConnReader: nil,
		Client:     client,
		Viewport:   newViewport,
		Style:      MessagesStyle.Copy(),
		FocusIndex: InputBox,
		InputBox:   NewInputBox(),
		SidePanel:  NewSidePanel(client),
	}
}

func (s State) Update(msg tea.Msg) (State, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case SendPrivMsg:
		// for sending slash commands
		// TODO: make this better
		if msg.Msg[0] == '/' {
			substrs := strings.SplitN(msg.Msg[1:], " ", 2)
			command := substrs[0]
			var params []string
			if len(substrs) > 1 {
				params = strings.Split(substrs[1], " ")
			}
			if strings.ToUpper(command) == commands.JOIN {
				s.Client.ActiveChannel = params[0]
				s.Client.Channels[params[0]] = client.Channel{}
				s.Client.SendCommand(command, params...)
				return s, nil
			}
			s.Client.SendCommand(command, params...)
		} else {
			// TODO: dont send command when activechannel == c.host
			fullMsg := s.Client.Nickname + ": " + msg.Msg
			if c, ok := s.Client.Channels[s.Client.ActiveChannel]; ok {
				c.History += fullMsg + client.CRLF

				s.Client.Channels[s.Client.ActiveChannel] = c
			}
			s.Client.SendCommand(commands.PRIVMSG, s.Client.ActiveChannel, msg.Msg)
			s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
		}

		return s, nil
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

	// TODO: maybe we still need this, not sure
	// s.Viewport.SetContent(s.Client.ActiveChannel.History)

	switch s.FocusIndex {
	case Viewport:
		*s.Viewport, cmd = s.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	case InputBox:
		s.InputBox, cmd = s.InputBox.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s *State) SetSize(width, height int) {
	s.InputBox.SetSize(width)
	s.SidePanel.SetSize(width, height, s.InputBox.Style.Copy().GetVerticalFrameSize())

	// -1 is for some extra invisible margin or padding or whatever (idk what it is tbh)
	newViewportWidth := (width * 8 / 10) - s.Style.GetHorizontalBorderSize() - 1
	newViewportHeight := height - s.InputBox.Style.GetVerticalFrameSize() - s.Style.GetVerticalBorderSize() - 1

	// copying the existing style is important here otherwise we'll end up with artifacting
	s.Viewport.Width = newViewportWidth
	s.Style = s.Style.Width(newViewportWidth)
	s.Viewport.Height = newViewportHeight
	s.Style = s.Style.Height(newViewportHeight)
	// we need to re-set the content because words wrap differently on different sizes
	s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
}

func (s *State) HandleCommand(msg parser.IRCMessage) {
	// TODO: handle different commands
	switch msg.Command {
	case commands.PING:
		s.Client.SendCommand(commands.PONG, msg.Parameters[0])
	case commands.PRIVMSG:
		nick := strings.SplitN(msg.Source, "!", 2)[0]
		channel := msg.Parameters[0]
		msgContent := msg.Parameters[1]
		fullMsg := fmt.Sprintf("%s: %s", nick, msgContent)

		if c, ok := s.Client.Channels[channel]; ok {
			c.History += fullMsg + client.CRLF

			s.Client.Channels[channel] = c
		}
		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	case commands.JOIN:
		nick := strings.SplitN(msg.Source, "!", 2)[0]
		channel := msg.Parameters[0]

		fullMsg := fmt.Sprintf("== %s has joined", nick)
		if c, ok := s.Client.Channels[channel]; ok {
			c.History += fullMsg + client.CRLF

			s.Client.Channels[channel] = c
		}
		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	case commands.QUIT:
		nick := strings.SplitN(msg.Source, "!", 2)[0]
		reason := msg.Parameters[0]
		fullMsg := fmt.Sprintf("== %s has quit (%s)", nick, reason)
		if c, ok := s.Client.Channels[s.Client.ActiveChannel]; ok {
			c.History += fullMsg + client.CRLF

			s.Client.Channels[s.Client.ActiveChannel] = c
		}
		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	case commands.RPL_NAMREPLY:
		// TODO: put the users into their respective channels
		// client := msg.Parameters[0]
		// chanSymbol := msg.Parameters[1]
		// channel := msg.Parameters[2]
		nicks := strings.Split(msg.Parameters[3], " ")

		s.SidePanel.Nicks = append(s.SidePanel.Nicks, nicks...)
	case commands.RPL_ENDOFNAMES:
		fullMsg := fmt.Sprintf("%s %s %s %s", msg.Tags, msg.Source, msg.Command, strings.Join(msg.Parameters, " "))
		if c, ok := s.Client.Channels[s.Client.ActiveChannel]; ok {
			c.History += fullMsg + client.CRLF

			s.Client.Channels[s.Client.ActiveChannel] = c
		}

		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	default:
		fullMsg := fmt.Sprintf("%s %s %s %s", msg.Tags, msg.Source, msg.Command, strings.Join(msg.Parameters, " "))
		if c, ok := s.Client.Channels[s.Client.ActiveChannel]; ok {
			c.History += fullMsg + client.CRLF

			s.Client.Channels[s.Client.ActiveChannel] = c
		}

		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	}
}

func (s State) View() string {
	top := lipgloss.JoinHorizontal(lipgloss.Right, s.Style.Render(s.Viewport.View()), s.SidePanel.View())
	screen := lipgloss.JoinVertical(0, top, s.InputBox.View())

	return ui.MainStyle.Render(screen)
}
