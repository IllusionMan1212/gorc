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
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/client"
)

type SidePanelState struct {
	Client   *client.Client
	Viewport viewport.Model
	Focused  bool
}

func (s *SidePanelState) getHeader() string {
	usersCount := 0
	if len(s.Client.Channels) != 0 {
		usersCount = len(s.Client.Channels[s.Client.ChannelIndex].Users)
	}
	separator := strings.Repeat("—", s.Viewport.Width) + "\n"
	header := fmt.Sprintf("Users in this channel (%d)\n", usersCount) + separator

	return header
}

func (s *SidePanelState) getLatestNicks() []string {
	nicks := make([]string, 0)

	if len(s.Client.Channels) != 0 {
		for nick, user := range s.Client.Channels[s.Client.ChannelIndex].Users {
			_nick := user.Prefix + nick
			nicks = append(nicks, _nick)
		}
	}

	return nicks
}

func NewSidePanel(client *client.Client) *SidePanelState {
	newViewport := viewport.New(0, 0)
	newViewport.Style = SidePanelStyle.Copy()
	newViewport.Wrap = viewport.Wrap

	return &SidePanelState{
		Client:   client,
		Viewport: newViewport,
	}
}

func (s SidePanelState) Update(msg tea.Msg) (SidePanelState, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.(type) {
	case SwitchChannelsMsg:
		s.UpdateNicks()
		return s, nil
	}

	if s.Focused {
		s.Viewport, cmd = s.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s *SidePanelState) Focus() {
	s.Focused = true
	s.Viewport.Style = s.Viewport.Style.BorderForeground(lipgloss.Color("105"))
}

func (s *SidePanelState) Blur() {
	s.Focused = false
	s.Viewport.Style = s.Viewport.Style.BorderForeground(lipgloss.Color("#EEE"))
}

func (s *SidePanelState) UpdateNicks() {
	header := s.getHeader()
	nicks := s.getLatestNicks()

	s.Viewport.SetContent(header + strings.Join(nicks, "\n"))
}

func (s *SidePanelState) SetSize(width, height, inputboxHeight int) {
	// We ceil here because width is an int and some fractions are lost
	// -1 is for some extra invisible margin or padding between the sidepanel and the inputbox
	newWidth := int(math.Ceil(math.Ceil((float64(width) * 2.0 / 10.0)) - float64(s.Viewport.Style.GetHorizontalFrameSize())))
	newHeight := height - inputboxHeight - s.Viewport.Style.GetVerticalFrameSize() - 1

	s.Viewport.Style = s.Viewport.Style.Width(newWidth)
	s.Viewport.Style = s.Viewport.Style.Height(newHeight)

	s.Viewport.Width = newWidth
	s.Viewport.Height = newHeight

	// on resize recalculate the width and re-set the viewport contents
	s.UpdateNicks()
}

func (s SidePanelState) View() string {
	return s.Viewport.View()
}