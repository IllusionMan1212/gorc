// gorc project
// Copyright (C) 2022 IllusionMan1212
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
	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/irc/commands"
	"github.com/illusionman1212/gorc/irc/handler"
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
	newViewport.MouseWheelEnabled = false

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
	var cmdsToProcess []tea.Cmd

	switch msg := msg.(type) {
	case cmds.ReceivedIRCMsgMsg:
		wasAtBottom := s.Viewport.AtBottom()
		s.Viewport.SetContent(strings.Join(s.Client.Channels[s.Client.ActiveChannelIndex].ScrollBackBuf, "\n"))

		if wasAtBottom {
			s.Client.Channels[s.Client.ActiveChannelIndex].AtBottom = true
			s.Viewport.GotoBottom()
		} else {
			s.Client.Channels[s.Client.ActiveChannelIndex].AtBottom = false
		}
		return s, nil
	case cmds.ReadHistoryFromDBMsg:
		s.Viewport.SetContent(strings.Join(s.Client.Channels[s.Client.ActiveChannelIndex].ScrollBackBuf, "\n"))
		switch msg.Direction {
		case "top":
			s.Viewport.GotoTop()
		case "up":
			s.Viewport.SetYOffset(s.Viewport.YOffset + int(msg.LinesRead))
		case "down":
			s.Viewport.SetYOffset(s.Viewport.YOffset - int(msg.LinesRead))
		case "bottom":
			s.Viewport.GotoBottom()
		}
		return s, nil
	case cmds.SendPrivMsgMsg:
		if msg.Msg[0] == '/' {
			cmd = handler.HandleSlashCommand(msg.Msg, s.Client)
			return s, cmd
		} else {
			if s.Client.ActiveChannel != s.Client.Host {
				fullMsg := s.Client.Nickname + ": " + msg.Msg
				msgOpts := irc.MsgFmtOpts{
					WithTimestamp: true,
				}

				s.Client.Channels[s.Client.ActiveChannelIndex].AppendMsg(s.Client.Host, s.Client.DB, msg.Timestamp, fullMsg, msgOpts)
				// TODO: make sure to only append the message to the history if server sends back no errors
				s.Client.SendCommand(commands.PRIVMSG, s.Client.ActiveChannel, msg.Msg)
				s.Viewport.SetContent(strings.Join(s.Client.Channels[s.Client.ActiveChannelIndex].ScrollBackBuf, "\n"))
				s.Viewport.GotoBottom()
			}
		}

		return s, nil
	case cmds.SwitchChannelsMsg:
		s.Viewport.SetContent(strings.Join(s.Client.Channels[s.Client.ActiveChannelIndex].ScrollBackBuf, "\n"))
		s.Viewport.GotoBottom()

		*s.SidePanel, cmd = s.SidePanel.Update(msg)
		return s, cmd
	case cmds.UpdateTabBarMsg:
		s.Client.LastTabIndexInTabBar = len(s.Client.Channels) - 1
		s.TabRenderingDirection = Left
		return s, nil
	case cmds.TickMsg:
		s.Client.Ticks++
		return s, cmds.Tick(s.Client)
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
				s.Blur()
				s.InputBox.Focus()
				s.SidePanel.Blur()
			case SidePanel:
				s.Blur()
				s.InputBox.Blur()
				s.SidePanel.Focus()
			}

			return s, textinput.Blink
		case "right", "left":
			// only cycle between tabs if focus isn't on the inputbox
			if s.FocusIndex == InputBox {
				break
			}

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

			if s.Client.ActiveChannelIndex > s.Client.LastTabIndexInTabBar {
				s.TabRenderingDirection = Left
				s.Client.LastTabIndexInTabBar = s.Client.ActiveChannelIndex
			} else if s.Client.ActiveChannelIndex < s.Client.FirstTabIndexInTabBar {
				s.Client.FirstTabIndexInTabBar = s.Client.ActiveChannelIndex
				s.TabRenderingDirection = Right
			}

			s.Client.ActiveChannel = s.Client.Channels[s.Client.ActiveChannelIndex].Name
			return s, cmds.SwitchChannels
		case "g":
			if s.FocusIndex == Viewport {
				s.Client.Channels[s.Client.ActiveChannelIndex].AtBottom = false
				if !s.Viewport.AtTop() {
					cmdsToProcess = append(cmdsToProcess, cmds.ReadHistoryFromDB(s.Client, irc.SCROLLBACK_MAX, "top"))
				}
			}
		case "G":
			if s.FocusIndex == Viewport {
				if !s.Viewport.AtBottom() {
					cmdsToProcess = append(cmdsToProcess, cmds.ReadHistoryFromDB(s.Client, irc.SCROLLBACK_MAX, "bottom"))
				}
				s.Client.Channels[s.Client.ActiveChannelIndex].AtBottom = true
			}
		case "j", "k", "d", "u", "f", "b", "up", "down":
			// TODO: what???
			// if s.Viewport.AtBottom() {
			// 	break
			// }

			if key == "j" || key == "d" || key == "f" || key == "down" {
				if s.Viewport.ScrollPercent() > 0.5 {
					excess := s.Viewport.ScrollPercent() - 0.5 + 0.1
					linesToRead := int(math.Floor(excess * float64(irc.SCROLLBACK_MAX)))
					cmdsToProcess = append(cmdsToProcess, cmds.ReadHistoryFromDB(s.Client, linesToRead, "down"))
				}
			} else {
				s.Client.Channels[s.Client.ActiveChannelIndex].AtBottom = false
				if s.Viewport.ScrollPercent() < 0.5 && s.Viewport.ScrollPercent() > 0.0 {
					excess := 0.5 - s.Viewport.ScrollPercent() + 0.1
					linesToRead := int(math.Floor(excess * float64(irc.SCROLLBACK_MAX)))
					cmdsToProcess = append(cmdsToProcess, cmds.ReadHistoryFromDB(s.Client, linesToRead, "up"))
				}
			}
		}
	}

	switch s.FocusIndex {
	case Viewport:
		*s.Viewport, cmd = s.Viewport.Update(msg)
		cmdsToProcess = append(cmdsToProcess, cmd)
	case InputBox:
		s.InputBox, cmd = s.InputBox.Update(msg)
		cmdsToProcess = append(cmdsToProcess, cmd)
	}

	*s.SidePanel, cmd = s.SidePanel.Update(msg)
	cmdsToProcess = append(cmdsToProcess, cmd)

	return s, tea.Batch(cmdsToProcess...)
}

func (s *State) Focus() {
	s.Viewport.Style = s.Viewport.Style.Copy().BorderForeground(ui.AccentColor)

	tab = tab.Copy().BorderForeground(ui.AccentColor)
	leftArrowDim = leftArrowDim.Copy().BorderForeground(ui.AccentColor)
	rightArrowDim = rightArrowDim.Copy().BorderForeground(ui.AccentColor)
	leftArrowLit = leftArrowLit.Copy().BorderForeground(ui.AccentColor)
	rightArrowLit = rightArrowLit.Copy().BorderForeground(ui.AccentColor)
	tabLine = tabLine.Copy().Foreground(ui.AccentColor)
}
func (s *State) Blur() {
	s.Viewport.Style = s.Viewport.Style.Copy().BorderForeground(ui.PrimaryColor)

	tab = tab.Copy().BorderForeground(ui.PrimaryColor)
	leftArrowDim = leftArrowDim.Copy().BorderForeground(ui.PrimaryColor)
	rightArrowDim = rightArrowDim.Copy().BorderForeground(ui.PrimaryColor)
	leftArrowLit = leftArrowLit.Copy().BorderForeground(ui.PrimaryColor)
	rightArrowLit = rightArrowLit.Copy().BorderForeground(ui.PrimaryColor)
	tabLine = tabLine.Copy().Foreground(ui.PrimaryColor)
}

func (s *State) SetSize(width, height int) {
	s.InputBox.SetSize(width)
	s.SidePanel.SetSize(width, height, s.InputBox.Style.GetVerticalFrameSize())

	// We floor because width is an int and some fractions are lost when casting
	// and also because we ceil the sidepanel's width
	// -3 for the tab bar height
	newWidth := int(math.Floor(float64(width)*8/10) - float64(s.Viewport.Style.GetHorizontalFrameSize()))
	newHeight := height - s.InputBox.Style.GetVerticalFrameSize() - s.Viewport.Style.GetVerticalFrameSize() - 3

	s.Viewport.Width = newWidth
	s.Viewport.Height = newHeight

	s.Viewport.Style = s.Viewport.Style.Width(newWidth)
	s.Viewport.Style = s.Viewport.Style.Height(newHeight)

	// we need this to render an empty viewport
	history := ""
	if len(s.Client.Channels) != 0 {
		history = strings.Join(s.Client.Channels[s.Client.ActiveChannelIndex].ScrollBackBuf, "\n")
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
		for i := s.Client.LastTabIndexInTabBar; i >= 0; i-- {
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
				s.Client.FirstTabIndexInTabBar = i + 1
				// dont render the newly added tab
				renderedTabs = renderedTabs[1:]
				break
			}
		}
	case Right:
		for i := s.Client.FirstTabIndexInTabBar; i < len(s.Client.Channels); i++ {
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
				s.Client.LastTabIndexInTabBar = i - 1
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

	if s.Client.FirstTabIndexInTabBar != 0 {
		leftArrow = leftArrowLit.Render("❰")
	}

	if s.Client.LastTabIndexInTabBar != len(s.Client.Channels)-1 {
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
