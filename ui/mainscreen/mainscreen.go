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
	"math"
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
	Viewport Window = iota
	SidePanel
	InputBox
)

type State struct {
	ConnReader *bufio.Reader
	Client     *client.Client
	Viewport   *viewport.Model
	FocusIndex Window

	InputBox  InputState
	SidePanel *SidePanelState
}

func NewMainScreen(client *client.Client) State {
	newViewport := &viewport.Model{
		HighPerformanceRendering: false,
		Wrap:                     viewport.Wrap,
		Style:                    MessagesStyle.Copy(),
	}

	return State{
		ConnReader: nil,
		Client:     client,
		Viewport:   newViewport,
		FocusIndex: InputBox,
		InputBox:   NewInputBox(),
		SidePanel:  NewSidePanel(client),
	}
}

func (s State) Update(msg tea.Msg) (State, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ReceivedIRCMsgMsg:
		s.Viewport.GotoBottom()
		return s, nil
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
				if _, ok := s.Client.Channels[params[0]]; !ok {
					s.Client.Channels[params[0]] = client.Channel{
						Users: make(map[string]client.User),
					}
				}
				s.Client.SendCommand(commands.JOIN, params...)
				return s, SwitchChannels
			}
			s.Client.SendCommand(command, params...)
		} else {
			if s.Client.ActiveChannel != s.Client.Host {
				fullMsg := s.Client.Nickname + ": " + msg.Msg
				if c, ok := s.Client.Channels[s.Client.ActiveChannel]; ok {
					c.History += fullMsg + client.CRLF

					s.Client.Channels[s.Client.ActiveChannel] = c
				}
				// TODO: make sure to only append the message to the history if server sends back no errors
				s.Client.SendCommand(commands.PRIVMSG, s.Client.ActiveChannel, msg.Msg)
				s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
				s.Viewport.GotoBottom()
			}
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

			if s.FocusIndex > 2 {
				s.FocusIndex = 0
			} else if s.FocusIndex < 0 {
				s.FocusIndex = 2
			}

			switch s.FocusIndex {
			case Viewport:
				s.Focus()
				s.InputBox.Blur()
				s.SidePanel.Blur()
			case InputBox:
				cmds = append(cmds, textinput.Blink)
				s.Blur()
				s.InputBox.Focus()
				s.SidePanel.Blur()
			case SidePanel:
				s.Blur()
				s.InputBox.Blur()
				s.SidePanel.Focus()
			}

			return s, tea.Batch(cmds...)
		}
	}

	switch s.FocusIndex {
	case Viewport:
		*s.Viewport, cmd = s.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	case InputBox:
		s.InputBox, cmd = s.InputBox.Update(msg)
		cmds = append(cmds, cmd)
	}

	*s.SidePanel, cmd = s.SidePanel.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s *State) Focus() {
	s.Viewport.Style = s.Viewport.Style.Copy().BorderForeground(lipgloss.Color("105"))
}
func (s *State) Blur() {
	s.Viewport.Style = s.Viewport.Style.Copy().BorderForeground(lipgloss.Color("#EEE"))
}

func (s *State) SetSize(width, height int) {
	s.InputBox.SetSize(width)
	s.SidePanel.SetSize(width, height, s.InputBox.Style.GetVerticalFrameSize())

	// We floor here because width is an int and some fractions are lost
	// and also because we ceil the sidepanel's width
	// -1 is for some extra invisible margin or padding between the main viewport and the inputbox
	newWidth := int(math.Floor(float64(width)*8.0/10.0) - float64(s.Viewport.Style.GetHorizontalFrameSize()))
	newHeight := height - s.InputBox.Style.GetVerticalFrameSize() - s.Viewport.Style.GetVerticalFrameSize() - 1

	s.Viewport.Width = newWidth
	s.Viewport.Height = newHeight

	s.Viewport.Style = s.Viewport.Style.Width(newWidth)
	s.Viewport.Style = s.Viewport.Style.Height(newHeight)

	// we need to re-set the content because words wrap differently on different sizes
	s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
}

func (s *State) HandleCommand(msg parser.IRCMessage) {
	// TODO: handle different commands
	switch msg.Command {
	case commands.PING:
		token := msg.Parameters[0]
		s.Client.SendCommand(commands.PONG, token)
	case commands.PRIVMSG:
		nick := strings.SplitN(msg.Source, "!", 2)[0]
		channel := msg.Parameters[0]
		msgContent := msg.Parameters[1]
		privMsg := fmt.Sprintf("%s: %s", nick, msgContent)

		if c, ok := s.Client.Channels[channel]; ok {
			c.History += privMsg + client.CRLF

			s.Client.Channels[channel] = c
		}
		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	case commands.JOIN:
		nick := strings.SplitN(msg.Source, "!", 2)[0]
		channel := msg.Parameters[0]

		joinMsg := fmt.Sprintf("== %s has joined", nick)
		if c, ok := s.Client.Channels[channel]; ok {
			c.History += joinMsg + client.CRLF
			if _, exists := c.Users[nick]; !exists {
				c.Users[nick] = client.User{}
			}

			s.Client.Channels[channel] = c
		}

		if channel == s.Client.ActiveChannel {
			s.SidePanel.UpdateNicks()
		}

		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	case commands.QUIT:
		nick := strings.SplitN(msg.Source, "!", 2)[0]
		reason := msg.Parameters[0]
		quitMsg := fmt.Sprintf("== %s has quit (%s)", nick, reason)

		for chanName, channel := range s.Client.Channels {
			channel.History += quitMsg + client.CRLF
			delete(channel.Users, nick)

			s.Client.Channels[chanName] = channel
		}

		s.SidePanel.UpdateNicks()
		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	case commands.PART:
		nick := strings.SplitN(msg.Source, "!", 2)[0]
		channel := msg.Parameters[0]
		reason := ""
		if len(msg.Parameters) > 1 {
			reason = msg.Parameters[1]
		}

		partMsg := fmt.Sprintf("== %s has left %s (%s)", nick, channel, reason)
		if c, ok := s.Client.Channels[channel]; ok {
			c.History += partMsg + client.CRLF
			delete(c.Users, nick)

			s.Client.Channels[channel] = c
		}

		s.SidePanel.UpdateNicks()
		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	case commands.RPL_NAMREPLY:
		// TODO: do i need these
		// client := msg.Parameters[0]
		// chanSymbol := msg.Parameters[1]
		channel := msg.Parameters[2]
		nicks := strings.Split(msg.Parameters[3], " ")

		if c, ok := s.Client.Channels[channel]; ok {
			for _, nick := range nicks {
				prefix := ""
				_nick := nick

				if commands.UserPrefixes[string(nick[0])] {
					prefix = string(nick[0])
					_nick = nick[1:]
				}

				c.Users[_nick] = client.User{
					Prefix: prefix,
				}
			}
		}

		if channel == s.Client.ActiveChannel {
			s.SidePanel.UpdateNicks()
		}
	case commands.RPL_ENDOFNAMES:
		// TODO: what do i do here lol
	default:
		fullMsg := fmt.Sprintf("%s %s %s %s", msg.Tags, msg.Source, msg.Command, strings.Join(msg.Parameters, " "))
		if c, ok := s.Client.Channels[s.Client.ActiveChannel]; ok {
			c.History += fullMsg + client.CRLF

			s.Client.Channels[s.Client.ActiveChannel] = c
		}

		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannel].History)
	}

	// send a receivedIRCmsg tea message so the ui can update
	// we also use this tea message to scroll the viewport down
	s.Client.Tea.Send(ReceivedIRCMsg())
}

func (s State) View() string {
	top := lipgloss.JoinHorizontal(lipgloss.Right, s.Viewport.View(), s.SidePanel.View())
	screen := lipgloss.JoinVertical(0, top, s.InputBox.View())

	return ui.MainStyle.Render(screen)
}
