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
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/ui"
)

type InputState struct {
	Input textinput.Model
	Style lipgloss.Style
}

func NewInputBox() InputState {
	input := textinput.New()
	input.Placeholder = "Send a message..."
	input.Focus()

	return InputState{
		Input: input,
		Style: InputboxStyle,
	}
}

func (s InputState) Update(msg tea.Msg) (InputState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "enter":
			value := s.Input.Value()
			if len(value) == 0 {
				return s, nil
			}

			s.Input.Reset()

			return s, cmds.SendPrivMsg(value)
		}
	}

	var cmd tea.Cmd
	s.Input, cmd = s.Input.Update(msg)
	return s, cmd
}

func (s *InputState) Focus() {
	s.Input.Focus()
	s.Style = s.Style.BorderForeground(ui.AccentColor)
}

func (s *InputState) Blur() {
	s.Input.Blur()
	s.Style = s.Style.BorderForeground(ui.PrimaryColor)
}

func (s *InputState) SetSize(width int) {
	// set a max width for the input field so it scrolls horizontally instead of wrapping to a newline
	// -4 for input prompt char, cursor, and some extra magical padding
	s.Input.Width = width - s.Style.GetHorizontalFrameSize() - 4
	s.Style = s.Style.Width(width - s.Style.GetHorizontalBorderSize())
}

func (s InputState) View() string {
	return s.Style.Render(s.Input.View())
}
