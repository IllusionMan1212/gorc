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

package login

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/ui"
)

type State struct {
	FocusIndex                int
	Inputs                    []textinput.Model
	TLS                       bool
	CanConnect                bool
	ConnectButtonBlurredStyle string
	ConnectButtonFocusedStyle string
	DialogStyle               lipgloss.Style
	WelcomeMsgStyle           lipgloss.Style
	Style                     lipgloss.Style
}

func NewLogin() State {
	state := State{
		Inputs:                    make([]textinput.Model, 5),
		ConnectButtonBlurredStyle: BlurredDisabledButton,
		ConnectButtonFocusedStyle: FocusedDisabledButton,
		DialogStyle:               DialogStyle,
		WelcomeMsgStyle:           WelcomeMsgStyle,
		Style:                     ui.MainStyle,
	}

	var t textinput.Model
	for i := range state.Inputs {
		t = textinput.New()
		t.Cursor.Style = CursorStyle
		t.Prompt = ""

		switch i {
		case 0:
			t.Placeholder = "Host"
			t.Focus()
			t.CharLimit = 40
			t.Width = len(t.Placeholder)
			t.TextStyle = FocusedStyle
			t.Validate = NoSpacesValidation
		case 1:
			t.Placeholder = "Port"
			t.CharLimit = 5
			t.Width = len(t.Placeholder)
			t.Validate = ValidatePort
		case 2:
			t.Placeholder = "Channel"
			t.CharLimit = 32
			t.Width = len(t.Placeholder)
		case 3:
			t.Placeholder = "Nickname"
			t.CharLimit = 32
			t.Width = len(t.Placeholder)
			t.Validate = NoSpacesValidation
		case 4:
			t.Placeholder = "Password"
			t.CharLimit = 64
			t.Width = len(t.Placeholder)
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		state.Inputs[i] = t
	}

	return state
}

func (s State) Update(msg tea.Msg) (State, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case " ", "enter":
			// run the Connect cmd when pressing "enter" while focused on the connect button
			// if the CanConnect flag is set
			if key == "enter" && s.FocusIndex == len(s.Inputs)+1 && s.CanConnect {
				return s, cmds.Connect
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
					s.Inputs[i].PromptStyle = FocusedStyle
					s.Inputs[i].TextStyle = FocusedStyle
					continue
				}

				s.Inputs[i].Blur()
				s.Inputs[i].PromptStyle = ui.DefaultStyle
				s.Inputs[i].TextStyle = ui.DefaultStyle
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
		s.Inputs[i].Width = max(len(s.Inputs[i].Placeholder), len(s.Inputs[i].Value()))
	}

	if s.Inputs[0].Value() != "" && s.Inputs[1].Value() != "" && s.Inputs[3].Value() != "" {
		s.ConnectButtonBlurredStyle = BlurredButton
		s.ConnectButtonFocusedStyle = FocusedButton
		s.CanConnect = true
	} else {
		s.ConnectButtonBlurredStyle = BlurredDisabledButton
		s.ConnectButtonFocusedStyle = FocusedDisabledButton
		s.CanConnect = false
	}

	return tea.Batch(cmds...)
}

func (s *State) SetSize(width, height int) {
	s.DialogStyle = s.DialogStyle.Width(width - s.DialogStyle.GetHorizontalFrameSize())
	s.DialogStyle = s.DialogStyle.Height(height*5/10 - s.DialogStyle.GetVerticalFrameSize())

	s.WelcomeMsgStyle = s.WelcomeMsgStyle.Width(width)
	s.WelcomeMsgStyle = s.WelcomeMsgStyle.Height(height * 5 / 10)
}

func (s State) View() string {
	var sb strings.Builder
	for i, input := range s.Inputs {
		sb.WriteString(input.View())
		if i < len(s.Inputs) {
			sb.WriteRune('\n')
		}
	}

	checkbox := BlurredCheckbox
	// if the checkbox is focused
	if s.FocusIndex == len(s.Inputs) {
		checkbox = FocusedCheckbox
		// if tls is enabled
		if s.TLS {
			checkbox = FocusedCheckboxChecked
		}
	}

	// if the checkbox is not focused and tls is enabled
	if s.FocusIndex != len(s.Inputs) && s.TLS {
		checkbox = BlurredCheckboxChecked
	}

	button := s.ConnectButtonBlurredStyle
	if s.FocusIndex == len(s.Inputs)+1 {
		button = s.ConnectButtonFocusedStyle
	}

	sb.WriteString(lipgloss.JoinVertical(lipgloss.Center, checkbox, button))

	screen := lipgloss.JoinVertical(lipgloss.Center, s.WelcomeMsgStyle.Render(WelcomeMsg), s.DialogStyle.Render(sb.String()))

	return ui.MainStyle.Render(screen)
}
