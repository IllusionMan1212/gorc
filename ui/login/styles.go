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
	"github.com/charmbracelet/lipgloss"
)

const WelcomeMsg = `                                        
                                        
  /$$$$$$   /$$$$$$   /$$$$$$   /$$$$$$$
 /$$__  $$ /$$__  $$ /$$__  $$ /$$_____/
| $$  \ $$| $$  \ $$| $$  \__/| $$      
| $$  | $$| $$  | $$| $$      | $$      
|  $$$$$$$|  $$$$$$/| $$      |  $$$$$$$
 \____  $$ \______/ |__/       \_______/
 /$$  \ $$                              
|  $$$$$$/                              
 \______/                               

`

var (
	FocusedButton = lipgloss.NewStyle().
			Background(lipgloss.Color("105")).
			Foreground(lipgloss.Color("#000")).
			MarginTop(1).
			Padding(0, 2).
			Align(lipgloss.Center).
			Render("Connect")
	BlurredButton = lipgloss.NewStyle().
			Background(lipgloss.Color("#505050")).
			Foreground(lipgloss.Color("#EEEEEE")).
			MarginTop(1).
			Padding(0, 2).
			Align(lipgloss.Center).
			Render("Connect")

	FocusedCheckbox = lipgloss.NewStyle().
			Foreground(lipgloss.Color("105")).
			MarginTop(1).
			Render("[ ] Enable TLS")

	BlurredCheckbox = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EEEEEE")).
			MarginTop(1).
			Render("[ ] Enable TLS")

	FocusedCheckboxChecked = lipgloss.NewStyle().
				Foreground(lipgloss.Color("105")).
				MarginTop(1).
				Render("[x] Enable TLS")

	BlurredCheckboxChecked = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#EEEEEE")).
				MarginTop(1).
				Render("[x] Enable TLS")

	DialogStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)

	WelcomeMsgStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
)
