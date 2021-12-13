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
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/client"
	"github.com/illusionman1212/gorc/ui"
	"github.com/illusionman1212/gorc/ui/login"
	"github.com/illusionman1212/gorc/ui/mainscreen"
	"golang.org/x/term"
)

type Screen int

const (
	Login = iota
	MainScreen
)

type UI struct {
	windowWidth   int
	windowHeight  int
	currentScreen Screen
	login         login.State
	mainScreen    mainscreen.State
}

type State struct {
	ui     UI
	client *client.Client
}

func initialUiState() UI {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// TODO: properly handle this error
		log.Fatal(err)
	}

	return UI{
		login:        login.NewLogin(),
		mainScreen:   mainscreen.NewMainScreen(),
		windowWidth:  width,
		windowHeight: height,
	}
}

func InitialState() State {
	s := State{
		client: &client.Client{},
		ui:     initialUiState(),
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
			// TODO: clean up and gracefully quit by leaving the IRC channel/server and then quit bubble.
			if s.client.TlsConn != nil {
				s.client.TlsConn.Close()
			}
			if s.client.TcpConn != nil {
				s.client.TcpConn.Close()
			}
			return s, tea.Quit
		}

	case tea.WindowSizeMsg:
		s.ui.windowWidth = msg.Width
		s.ui.windowHeight = msg.Height

		ui.MainStyle = ui.MainStyle.
			Width(msg.Width).
			Height(msg.Height)

		s.ui.mainScreen.Viewport.Width = msg.Width
		s.ui.mainScreen.Viewport.Height = msg.Height - mainscreen.InputBoxHeight

		return s, nil
	case login.ConnectingMsg:
		host := s.ui.login.Inputs[0].Value()
		port := s.ui.login.Inputs[1].Value()
		channel := s.ui.login.Inputs[2].Value()
		nickname := s.ui.login.Inputs[3].Value()
		password := s.ui.login.Inputs[4].Value()
		tlsEnabled := s.ui.login.TLS

		s.ui.currentScreen = MainScreen

		// create new client with the provided host and port
		client := client.NewClient(host, port, tlsEnabled)
		client.Register(nickname, password, channel)
		s.client = client
		s.ui.mainScreen.Client = client

		if tlsEnabled {
			r := bufio.NewReaderSize(client.TlsConn, 512)
			s.ui.mainScreen.Reader = r
		} else {
			r := bufio.NewReaderSize(client.TcpConn, 512)
			s.ui.mainScreen.Reader = r
		}

		return s, mainscreen.InitialRead
	}

	// switch between which screen is currently active and update its state
	switch s.ui.currentScreen {
	case Login:
		state, cmd := s.ui.login.Update(msg)
		s.ui.login = state
		return s, cmd
	case MainScreen:
		state, cmd := s.ui.mainScreen.Update(msg)
		s.ui.mainScreen = state
		return s, cmd
	}

	return s, nil
}

func (s State) View() string {
	// switch between which screen is currently active and render it.
	switch s.ui.currentScreen {
	case Login:
		return s.ui.login.View()
	case MainScreen:
		return s.ui.mainScreen.View()
	}

	return "Error: this screen shouldn't ever show"
}
