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

package app

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/ui"
	"github.com/illusionman1212/gorc/ui/login"
	"golang.org/x/term"
)

type Screen int

const (
	Login = iota
	MainScreen
)

// TODO: have a global State that holds a separate client State and separate tui State
type State struct {
	host string
	port string

	nick     string
	password string

	windowWidth   int
	windowHeight  int
	currentScreen Screen

	login login.State
}

func InitialState() State {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	s := State{
		login:         login.NewLogin(),
		windowWidth:   width,
		windowHeight:  height,
		currentScreen: Login,
	}

	return s
}

func (s State) Init() tea.Cmd {
	return nil
}

func (s State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// TODO: clean up and gracefully quit by leaving the IRC channel/server and then quit bubble.
			return s, tea.Quit
		}

	case tea.WindowSizeMsg:
		s.windowWidth = msg.Width
		s.windowHeight = msg.Height

		ui.MainStyle = ui.MainStyle.
			Width(msg.Width).
			Height(msg.Height)

		return s, nil
	}

	// switch between which screen is currently active and update its state
	switch s.currentScreen {
	case Login:
		loginState, loginCmd := s.login.Update(msg)
		s.login = loginState
		return s, loginCmd
	}

	return s, nil
}

func (s State) View() string {
	// switch between which screen is currently active and render it.
	switch s.currentScreen {
	case Login:
		return s.login.View()
	}

	return "Error: this screen shouldn't ever show"
}
