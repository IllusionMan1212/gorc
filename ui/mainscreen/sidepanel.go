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
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/client"
)

type SidePanelState struct {
	// TODO: make a viewport and use that to display the names
	// OR: we could use a list
	Style    lipgloss.Style
	Nicks    []string
	Viewport viewport.Model
	Content  string
}

func NewSidePanel(client *client.Client) *SidePanelState {
	return &SidePanelState{
		Style: SidePanelStyle.Copy(),
	}
}

func (s SidePanelState) Update(msg tea.Msg) (SidePanelState, tea.Cmd) {
	return s, nil
}

func (s *SidePanelState) SetSize(width, height, inputboxHeight int) {
	s.Style = s.Style.Copy().Width(width * 2 / 10)
	// -1 is for some extra invisible margin or padding or whatever (idk what it is tbh)
	s.Style = s.Style.Copy().Height(height - inputboxHeight - s.Style.GetVerticalBorderSize() - 1)
}

func (s SidePanelState) View() string {
	separator := strings.Repeat("-", s.Style.GetWidth()) + "\n"
	header := fmt.Sprintf("Users in this channel (%d)\n", len(s.Nicks)) + separator
	nicks := strings.Join(s.Nicks, "\n")

	return s.Style.Render(header + nicks)
}
