// gorc project
// Copyright (C) 2021 IllusionMan1212 and contributors
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

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	border = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	mainStyle = lipgloss.NewStyle().
			Width(172).
			Height(40)

	dialogStyle = lipgloss.NewStyle().
			Border(border, true).
			Align(lipgloss.Center).
			Width(50).
			Padding(3, 6).
			BorderForeground(lipgloss.Color("105"))

	welcomeMsgStyle = lipgloss.NewStyle().
			Padding(0, 5).
			MarginBottom(5)

	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("105"))
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("105"))
	noStyle      = lipgloss.NewStyle()

	focusedButton = lipgloss.NewStyle().
			Background(lipgloss.Color("105")).
			Foreground(lipgloss.Color("255")).
			MarginTop(1).
			Padding(0, 2).
			Align(lipgloss.Center).
			Render("Connect")
	blurredButton = lipgloss.NewStyle().
			Background(lipgloss.Color("#809070")).
			Foreground(lipgloss.Color("#EEEEEE")).
			MarginTop(1).
			Padding(0, 2).
			Align(lipgloss.Center).
			Render("Connect")

	focusedCheckbox = lipgloss.NewStyle().
			Foreground(lipgloss.Color("105")).
			MarginTop(1).
			Render("[ ] Enable TLS")

	blurredCheckbox = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EEEEEE")).
			MarginTop(1).
			Render("[ ] Enable TLS")

	focusedCheckedCheckbox = lipgloss.NewStyle().
				Foreground(lipgloss.Color("105")).
				MarginTop(1).
				Render("[x] Enable TLS")

	blurredCheckedCheckbox = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#EEEEEE")).
				MarginTop(1).
				Render("[x] Enable TLS")
)

type loginState struct {
	focusIndex int
	inputs     []textinput.Model

	tls bool
}

// TODO: have a global state that holds a separate client state and separate tui state
type state struct {
	host string
	port string

	nick     string
	password string

	windowWidth  int
	windowHeight int

	login loginState
}

func initialLoginState() loginState {
	ls := loginState{
		inputs: make([]textinput.Model, 5),
	}

	var t textinput.Model
	for i := range ls.inputs {
		t = textinput.NewModel()
		t.CursorStyle = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "Host"
			t.Focus()
			t.CharLimit = 40
			t.TextStyle = focusedStyle
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

		ls.inputs[i] = t
	}

	return ls
}

func initialState() state {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	s := state{
		login:        initialLoginState(),
		windowWidth:  width,
		windowHeight: height,
	}

	return s
}

func (s state) Init() tea.Cmd {
	return textinput.Blink
}

func (s state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down", "space":
			str := msg.String()

			// run the newClient cmd when pressing "enter" while focused on the connect button
			if str == "enter" && s.login.focusIndex == len(s.login.inputs)+1 {
				return s, tea.Quit // TODO: change this to create a newClient and connect using a tea.Cmd
			}

			// toggle tls state when pressing "enter" while focused on the checkbox
			// TODO: figure out why space isn't working here. all other keys works just fine
			if (str == "enter" || str == "space") && s.login.focusIndex == len(s.login.inputs) {
				s.login.tls = !s.login.tls
				return s, nil
			}

			if str == "up" || str == "shift+tab" {
				s.login.focusIndex--
			} else {
				s.login.focusIndex++
			}

			if s.login.focusIndex > len(s.login.inputs)+1 {
				s.login.focusIndex = 0
			} else if s.login.focusIndex < 0 {
				s.login.focusIndex = len(s.login.inputs) + 1
			}

			cmds := make([]tea.Cmd, len(s.login.inputs))
			for i := 0; i <= len(s.login.inputs)-1; i++ {
				if i == s.login.focusIndex {
					cmds[i] = s.login.inputs[i].Focus()
					s.login.inputs[i].PromptStyle = focusedStyle
					s.login.inputs[i].TextStyle = focusedStyle
					continue
				}

				s.login.inputs[i].Blur()
				s.login.inputs[i].PromptStyle = noStyle
				s.login.inputs[i].TextStyle = noStyle
			}

			return s, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		// NOTE: this kinda works. still need to fix it tho
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			log.Print(err)
			return s, tea.Quit
		}

		s.windowWidth = width
		s.windowHeight = height

		mainStyle = mainStyle.
			Width(width).
			Height(height)

		return s, nil
	}

	cmd := s.updateInputs(msg)

	return s, cmd
}

func (s *state) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(s.login.inputs))

	for i := range s.login.inputs {
		s.login.inputs[i], cmds[i] = s.login.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (s state) View() string {
	var sb strings.Builder

	for i := range s.login.inputs {
		sb.WriteString(s.login.inputs[i].View())
		if i < len(s.login.inputs)-1 {
			sb.WriteRune('\n')
		}
	}

	checkbox := &blurredCheckbox
	// if the checkbox is focused
	if s.login.focusIndex == len(s.login.inputs) {
		checkbox = &focusedCheckbox
		// if tls is enabled
		if s.login.tls {
			checkbox = &focusedCheckedCheckbox
		}
	}

	// if the checkbox is not focused and tls is enabled
	if s.login.focusIndex != len(s.login.inputs) && s.login.tls {
		checkbox = &blurredCheckedCheckbox
	}

	fmt.Fprintf(&sb, "\n%s\n", *checkbox)

	button := &blurredButton
	if s.login.focusIndex == len(s.login.inputs)+1 {
		button = &focusedButton
	}
	fmt.Fprintf(&sb, "\n%s\n", *button)

	ui := lipgloss.JoinVertical(0, welcomeMsgStyle.Render(welcomeMsg), dialogStyle.Render(sb.String()))

	final := lipgloss.Place(s.windowWidth, s.windowHeight, lipgloss.Center, lipgloss.Top, ui)

	return mainStyle.Render(final)
}
