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

package app

import (
	"bufio"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/client"
	"github.com/illusionman1212/gorc/ui"
	"github.com/illusionman1212/gorc/ui/login"
	"github.com/illusionman1212/gorc/ui/mainscreen"
)

type Screen int

const (
	Login = iota
	MainScreen
)

type UI struct {
	CurrentScreen Screen
	Login         login.State
	MainScreen    mainscreen.State
}

type State struct {
	UI     UI
	Client *client.Client
}

func initialUiState(client *client.Client) UI {
	return UI{
		Login:      login.NewLogin(),
		MainScreen: mainscreen.NewMainScreen(client),
	}
}

func InitialState() State {
	client := &client.Client{}

	s := State{
		Client: client,
		UI:     initialUiState(client),
	}

	return s
}

func (s State) Init() tea.Cmd {
	return textinput.Blink
}

func (s State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, s.Quit
		}

	case tea.WindowSizeMsg:
		ui.MainStyle = ui.MainStyle.
			Width(msg.Width).
			Height(msg.Height)

		s.UI.MainScreen.SetSize(msg.Width, msg.Height)

		return s, nil
	case login.ConnectingMsg:
		host := s.UI.Login.Inputs[0].Value()
		port := s.UI.Login.Inputs[1].Value()
		channel := s.UI.Login.Inputs[2].Value()
		nickname := s.UI.Login.Inputs[3].Value()
		password := s.UI.Login.Inputs[4].Value()
		tlsEnabled := s.UI.Login.TLS

		s.UI.CurrentScreen = MainScreen

		// create new client with the provided host and port
		(*s.Client) = client.NewClient(host, port, tlsEnabled)
		s.Client.Register(nickname, password, channel)

		// 512 bytes as a base + 8192 additional bytes for tags
		r := bufio.NewReaderSize(s.Client.TcpConn, 8192+512)
		s.UI.MainScreen.ConnReader = r

		go s.UI.MainScreen.ReadLoop()

		return s, textinput.Blink
	}

	// switch between which screen is currently active and update its state
	switch s.UI.CurrentScreen {
	case Login:
		state, cmd := s.UI.Login.Update(msg)
		s.UI.Login = state
		return s, cmd
	case MainScreen:
		state, cmd := s.UI.MainScreen.Update(msg)
		s.UI.MainScreen = state
		return s, cmd
	}

	return s, nil
}

func (s State) View() string {
	// switch between which screen is currently active and render it.
	switch s.UI.CurrentScreen {
	case Login:
		return s.UI.Login.View()
	case MainScreen:
		return s.UI.MainScreen.View()
	}

	return "Error: this screen shouldn't ever show"
}
