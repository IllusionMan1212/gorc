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

package login

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/ui"
)

type State struct {
	FocusIndex int
	Inputs     []textinput.Model
	TLS        bool
}

func NewLogin() State {
	ls := State{
		Inputs: make([]textinput.Model, 5),
	}

	var t textinput.Model
	for i := range ls.Inputs {
		t = textinput.New()
		t.CursorStyle = ui.CursorStyle
		t.Prompt = "| "

		switch i {
		case 0:
			t.Placeholder = "Host"
			t.Focus()
			t.CharLimit = 40
			t.TextStyle = ui.FocusedStyle
		case 1:
			t.Placeholder = "Port"
			t.CharLimit = 5
		case 2:
			t.Placeholder = "Channel"
			t.CharLimit = 32
		case 3:
			t.Placeholder = "Nickname"
			t.CharLimit = 32
		case 4:
			t.Placeholder = "Password"
			t.CharLimit = 32
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		ls.Inputs[i] = t
	}

	return ls
}

func (s State) Update(msg tea.Msg) (State, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case " ", "enter":
			// run the newClient cmd when pressing "enter" while focused on the connect button
			if key == "enter" && s.FocusIndex == len(s.Inputs)+1 {
				return s, connect
			}

			// toggle tls state when pressing "enter" or "space" while focused on the checkbox
			if (key == "enter" || key == " ") && s.FocusIndex == len(s.Inputs) {
				s.TLS = !s.TLS
				return s, nil
			}
		case "tab", "shift+tab", "up", "down":
			if key == "up" || key == "shift+tab" {
				s.FocusIndex--
			} else {
				s.FocusIndex++
			}

			if s.FocusIndex > len(s.Inputs)+1 {
				s.FocusIndex = 0
			} else if s.FocusIndex < 0 {
				s.FocusIndex = len(s.Inputs) + 1
			}

			cmds := make([]tea.Cmd, len(s.Inputs))
			for i := 0; i <= len(s.Inputs)-1; i++ {
				if i == s.FocusIndex {
					cmds[i] = s.Inputs[i].Focus()
					s.Inputs[i].PromptStyle = ui.FocusedStyle
					s.Inputs[i].TextStyle = ui.FocusedStyle
					continue
				}

				s.Inputs[i].Blur()
				s.Inputs[i].PromptStyle = ui.NoStyle
				s.Inputs[i].TextStyle = ui.NoStyle
			}

			return s, tea.Batch(cmds...)
		}
	}

	cmd := s.updateInputs(msg)

	return s, cmd
}

func (s *State) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(s.Inputs))

	for i := range s.Inputs {
		s.Inputs[i], cmds[i] = s.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (s State) View() string {
	var sb strings.Builder

	for i := range s.Inputs {
		sb.WriteString(s.Inputs[i].View())
		if i < len(s.Inputs)-1 {
			sb.WriteRune('\n')
		}
	}

	checkbox := &BlurredCheckbox
	// if the checkbox is focused
	if s.FocusIndex == len(s.Inputs) {
		checkbox = &FocusedCheckbox
		// if tls is enabled
		if s.TLS {
			checkbox = &FocusedCheckboxChecked
		}
	}

	// if the checkbox is not focused and tls is enabled
	if s.FocusIndex != len(s.Inputs) && s.TLS {
		checkbox = &BlurredCheckboxChecked
	}

	fmt.Fprintf(&sb, "\n%s\n", *checkbox)

	button := &BlurredButton
	if s.FocusIndex == len(s.Inputs)+1 {
		button = &FocusedButton
	}
	fmt.Fprintf(&sb, "\n%s\n", *button)

	screen := lipgloss.JoinVertical(0, WelcomeMsgStyle.Render(WelcomeMsg), ui.DialogStyle.Render(sb.String()))

	final := lipgloss.Place(ui.MainStyle.GetWidth(), ui.MainStyle.GetHeight(), lipgloss.Center, lipgloss.Top, screen)

	return ui.MainStyle.Render(final)
}
