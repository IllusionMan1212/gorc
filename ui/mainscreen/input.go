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
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputState struct {
	Input textinput.Model
	Style lipgloss.Style
}

func NewInputBox() InputState {
	input := textinput.NewModel()
	input.Placeholder = "Send a message..."
	// input.Focus()
	state := InputState{input, InputboxStyle}

	return state
}

func (s InputState) Update(msg tea.Msg) (InputState, tea.Cmd) {
	cmd := s.updateInputs(msg)
	return s, cmd
}

func (s *InputState) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.Input, cmd = s.Input.Update(msg)

	return cmd
}

func (s *InputState) SetSize(width int) {
	s.Style = s.Style.Width(width - s.Style.GetHorizontalBorderSize())
}

func (s InputState) View() string {
	return s.Style.Render(s.Input.View())
}
