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
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/ui"
)

type SidePanelState struct {
	Client   *irc.Client
	Viewport viewport.Model
	Focused  bool
}

func (s *SidePanelState) getHeader() string {
	usersCount := 0
	if len(s.Client.Channels) != 0 {
		usersCount = len(s.Client.Channels[s.Client.ActiveChannelIndex].Users)
	}
	separator := strings.Repeat("—", s.Viewport.Width) + "\n"
	header := fmt.Sprintf("Users in this channel (%d)\n", usersCount) + separator

	return header
}

func (s *SidePanelState) getLatestNicks() []string {
	nicks := make([]string, 0)

	if len(s.Client.Channels) != 0 {
		for nick, user := range s.Client.Channels[s.Client.ActiveChannelIndex].Users {
			_nick := user.Prefix + nick
			nicks = append(nicks, _nick)
		}
	}

	return nicks
}

func NewSidePanel(client *irc.Client) *SidePanelState {
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
	var cmdsToProcess []tea.Cmd

	switch msg := msg.(type) {
	case cmds.SwitchChannelsMsg:
		s.UpdateNicks()
		return s, nil
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "g":
			if s.Focused {
				s.Viewport.GotoTop()
			}
		case "G":
			if s.Focused {
				s.Viewport.GotoBottom()
			}
		}
	}

	if s.Focused {
		s.Viewport, cmd = s.Viewport.Update(msg)
		cmdsToProcess = append(cmdsToProcess, cmd)
	}

	return s, tea.Batch(cmdsToProcess...)
}

func (s *SidePanelState) Focus() {
	s.Focused = true
	s.Viewport.Style = s.Viewport.Style.BorderForeground(ui.AccentColor)
}

func (s *SidePanelState) Blur() {
	s.Focused = false
	s.Viewport.Style = s.Viewport.Style.BorderForeground(ui.PrimaryColor)
}

func (s *SidePanelState) UpdateNicks() {
	header := s.getHeader()
	nicks := s.getLatestNicks()

	s.Viewport.SetContent(header + strings.Join(nicks, "\n"))
}

func (s *SidePanelState) SetSize(width, height, inputboxHeight int) {
	// We ceil here because width is an int and some fractions are lost
	// -1 is for some extra invisible margin or padding between the sidepanel and the inputbox
	newWidth := int(math.Ceil(math.Ceil((float64(width) * 2 / 10)) - float64(s.Viewport.Style.GetHorizontalFrameSize())))
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
