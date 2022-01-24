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
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/irc/commands"
	"github.com/illusionman1212/gorc/ui"
)

type Window int
type TabDirection int

const (
	Viewport Window = iota
	SidePanel
	InputBox
)

const (
	Left TabDirection = iota
	Right
)

// NOTE: I ABSOLUTELY HATE THIS
var FirstTabIndexInTabBar = 0
var LastTabIndexInTabBar = 0

type State struct {
	Client                *irc.Client
	Viewport              *viewport.Model
	FocusIndex            Window
	TabRenderingDirection TabDirection

	InputBox  InputState
	SidePanel *SidePanelState
}

func NewMainScreen(client *irc.Client) State {
	newViewport := viewport.New(0, 0)
	newViewport.Wrap = viewport.Wrap
	newViewport.Style = MessagesStyle.Copy()

	return State{
		Client:                client,
		Viewport:              &newViewport,
		FocusIndex:            InputBox,
		InputBox:              NewInputBox(),
		SidePanel:             NewSidePanel(client),
		TabRenderingDirection: Right,
	}
}

func (s State) Update(msg tea.Msg) (State, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ReceivedIRCMsgMsg:
		wasAtBottom := s.Viewport.AtBottom()
		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannelIndex].History)

		if wasAtBottom {
			s.Viewport.GotoBottom()
		}
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
				// TODO: handle comma separated channels
				channel := params[0]
				s.Client.ActiveChannel = channel
				for i, c := range s.Client.Channels {
					if c.Name == channel {
						s.Client.ActiveChannelIndex = i
						break
					} else if i == len(s.Client.Channels)-1 {
						s.Client.Channels = append(s.Client.Channels, irc.Channel{
							Name:  channel,
							Users: make(map[string]irc.User),
						})
						s.Client.ActiveChannelIndex = len(s.Client.Channels) - 1
						LastTabIndexInTabBar = len(s.Client.Channels) - 1
						s.TabRenderingDirection = Left
					}
				}
				s.Client.SendCommand(commands.JOIN, params...)
				return s, SwitchChannels
			}
			s.Client.SendCommand(command, params...)
		} else {
			if s.Client.ActiveChannel != s.Client.Host {
				fullMsg := s.Client.Nickname + ": " + msg.Msg
				s.Client.Channels[s.Client.ActiveChannelIndex].History += fullMsg + irc.CRLF
				// TODO: make sure to only append the message to the history if server sends back no errors
				s.Client.SendCommand(commands.PRIVMSG, s.Client.ActiveChannel, msg.Msg)
				s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannelIndex].History)
				s.Viewport.GotoBottom()
			}
		}

		return s, nil
	case SwitchChannelsMsg:
		s.Viewport.SetContent(s.Client.Channels[s.Client.ActiveChannelIndex].History)
		s.Viewport.GotoBottom()

		*s.SidePanel, cmd = s.SidePanel.Update(msg)
		cmds = append(cmds, cmd)

		return s, tea.Batch(cmds...)
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
		case "right", "left":
			// we need this so that the viewport doesnt scroll to the bottom
			if len(s.Client.Channels) == 1 {
				return s, nil
			}

			if key == "right" {
				s.Client.ActiveChannelIndex++
			} else {
				s.Client.ActiveChannelIndex--
			}

			if s.Client.ActiveChannelIndex >= len(s.Client.Channels) {
				s.Client.ActiveChannelIndex = 0
			} else if s.Client.ActiveChannelIndex < 0 {
				s.Client.ActiveChannelIndex = len(s.Client.Channels) - 1
			}

			if s.Client.ActiveChannelIndex > LastTabIndexInTabBar {
				s.TabRenderingDirection = Left
				LastTabIndexInTabBar = s.Client.ActiveChannelIndex
			} else if s.Client.ActiveChannelIndex < FirstTabIndexInTabBar {
				FirstTabIndexInTabBar = s.Client.ActiveChannelIndex
				s.TabRenderingDirection = Right
			}

			s.Client.ActiveChannel = s.Client.Channels[s.Client.ActiveChannelIndex].Name
			return s, SwitchChannels
		case "g":
			if s.FocusIndex == Viewport {
				s.Viewport.GotoTop()
			}
		case "G":
			if s.FocusIndex == Viewport {
				s.Viewport.GotoBottom()
			}
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

	tab = tab.Copy().BorderForeground(lipgloss.Color("105"))
	leftArrowDim = leftArrowDim.Copy().BorderForeground(lipgloss.Color("105"))
	rightArrowDim = rightArrowDim.Copy().BorderForeground(lipgloss.Color("105"))
	leftArrowLit = leftArrowLit.Copy().BorderForeground(lipgloss.Color("105"))
	rightArrowLit = rightArrowLit.Copy().BorderForeground(lipgloss.Color("105"))
	tabLine = tabLine.Copy().Foreground(lipgloss.Color("105"))
}
func (s *State) Blur() {
	s.Viewport.Style = s.Viewport.Style.Copy().BorderForeground(lipgloss.Color("#EEE"))

	tab = tab.Copy().UnsetBorderForeground()
	leftArrowDim = leftArrowDim.Copy().UnsetBorderForeground()
	rightArrowDim = rightArrowDim.Copy().UnsetBorderForeground()
	leftArrowLit = leftArrowLit.Copy().UnsetBorderForeground()
	rightArrowLit = rightArrowLit.Copy().UnsetBorderForeground()
	tabLine = tabLine.Copy().UnsetForeground()
}

func (s *State) SetSize(width, height int) {
	s.InputBox.SetSize(width)
	s.SidePanel.SetSize(width, height, s.InputBox.Style.GetVerticalFrameSize())

	// We floor here because width is an int and some fractions are lost
	// and also because we ceil the sidepanel's width
	// -3 for the tab bar height
	newWidth := int(math.Floor(float64(width)*8.0/10.0) - float64(s.Viewport.Style.GetHorizontalFrameSize()))
	newHeight := height - s.InputBox.Style.GetVerticalFrameSize() - s.Viewport.Style.GetVerticalFrameSize() - 3

	s.Viewport.Width = newWidth
	s.Viewport.Height = newHeight

	s.Viewport.Style = s.Viewport.Style.Width(newWidth)
	s.Viewport.Style = s.Viewport.Style.Height(newHeight)

	// we need this to render an empty viewport
	history := ""
	if len(s.Client.Channels) != 0 {
		history = s.Client.Channels[s.Client.ActiveChannelIndex].History
	}

	// we need to re-set the content because words wrap differently on different sizes
	s.Viewport.SetContent(history)
	s.Viewport.SetYOffset(s.Viewport.YOffset)
}

func (s State) buildTabBar(rightArrow string, leftArrow string) string {
	var renderedTabs []string
	tabs := ""

	switch s.TabRenderingDirection {
	case Left:
		for i := LastTabIndexInTabBar; i >= 0; i-- {
			if s.Client.Channels[i].Name == s.Client.ActiveChannel {
				renderedTabs = append([]string{activeTab.Render(s.Client.Channels[i].Name)}, renderedTabs...)
			} else {
				renderedTabs = append([]string{tab.Render(s.Client.Channels[i].Name)}, renderedTabs...)
			}

			tabs = lipgloss.JoinHorizontal(
				lipgloss.Top,
				renderedTabs...,
			)

			if lipgloss.Width(tabs) > lipgloss.Width(s.Viewport.View())-lipgloss.Width(leftArrow)-lipgloss.Width(rightArrow) {
				// set the first tab to be displayed to the index of the previous tab in the loop
				FirstTabIndexInTabBar = i + 1
				// dont render the newly added tab
				renderedTabs = renderedTabs[1:]
				break
			}
		}
	case Right:
		for i := FirstTabIndexInTabBar; i < len(s.Client.Channels); i++ {
			if s.Client.Channels[i].Name == s.Client.ActiveChannel {
				renderedTabs = append(renderedTabs, activeTab.Render(s.Client.Channels[i].Name))
			} else {
				renderedTabs = append(renderedTabs, tab.Render(s.Client.Channels[i].Name))
			}

			tabs = lipgloss.JoinHorizontal(
				lipgloss.Top,
				renderedTabs...,
			)

			if lipgloss.Width(tabs) > lipgloss.Width(s.Viewport.View())-lipgloss.Width(leftArrow)-lipgloss.Width(rightArrow) {
				// set the last tab to be displayed to the index of the previous tab in the loop
				LastTabIndexInTabBar = i - 1
				// dont render the newly added tab
				renderedTabs = renderedTabs[:len(renderedTabs)-1]
				break
			}
		}
	}

	tabs = lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderedTabs...,
	)

	return tabs
}

func (s State) View() string {
	leftArrow := leftArrowDim.Render("❰")
	rightArrow := rightArrowDim.Render("❱")

	tabs := s.buildTabBar(leftArrow, rightArrow)
	tabBar := strings.Builder{}

	if FirstTabIndexInTabBar != 0 {
		leftArrow = leftArrowLit.Render("❰")
	}

	if LastTabIndexInTabBar != len(s.Client.Channels)-1 {
		rightArrow = rightArrowLit.Render("❱")
	}

	tabBarLine := tabLine.Render(
		strings.Repeat(
			"─",
			max(0, lipgloss.Width(s.Viewport.View())-lipgloss.Width(tabs)-lipgloss.Width(rightArrow)-lipgloss.Width(leftArrow)),
		),
	)
	tabs = lipgloss.JoinHorizontal(lipgloss.Bottom, leftArrow, tabs, tabBarLine, rightArrow)
	tabBar.WriteString(tabs)

	leftSide := lipgloss.JoinVertical(0, tabBar.String(), s.Viewport.View())
	top := lipgloss.JoinHorizontal(lipgloss.Right, leftSide, s.SidePanel.View())
	screen := lipgloss.JoinVertical(0, top, s.InputBox.View())

	return ui.MainStyle.Render(screen)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
