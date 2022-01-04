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
	"github.com/illusionman1212/gorc/client"
)

type InputState struct {
	Input  textinput.Model
	Style  lipgloss.Style
	Client *client.Client
}

func NewInputBox(client *client.Client) InputState {
	input := textinput.NewModel()
	input.Placeholder = "Send a message..."
	input.Focus()
	state := InputState{
		Input:  input,
		Style:  InputboxStyle.Copy(),
		Client: client,
	}

	return state
}

func (s InputState) Update(msg tea.Msg) (InputState, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "enter":
			value := s.Input.Value()
			s.Input.Reset()

			return s, s.SendingPrivMsg(value)
		}
	}

	var cmd tea.Cmd
	s.Input, cmd = s.Input.Update(msg)
	return s, cmd
}

func (s *InputState) SetSize(width int) {
	s.Style = s.Style.Width(width - s.Style.GetHorizontalBorderSize())
}

func (s InputState) View() string {
	return s.Style.Render(s.Input.View())
}
