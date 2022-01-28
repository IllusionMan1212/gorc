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
	"github.com/illusionman1212/gorc/ui"
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
	BlurredDisabledButton = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ui.DisabledColor).
				Foreground(ui.DisabledColor).
				MarginTop(1).
				Padding(0, 2).
				Align(lipgloss.Center).
				Render("Connect")
	FocusedDisabledButton = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ui.DisabledColorFocus).
				Foreground(ui.DisabledColorFocus).
				MarginTop(1).
				Padding(0, 2).
				Align(lipgloss.Center).
				Render("Connect")
	FocusedButton = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ui.AccentColor).
			Foreground(ui.AccentColor).
			MarginTop(1).
			Padding(0, 2).
			Align(lipgloss.Center).
			Bold(true).
			Render("Connect")
	BlurredButton = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ui.PrimaryColor).
			Foreground(ui.PrimaryColor).
			MarginTop(1).
			Padding(0, 2).
			Align(lipgloss.Center).
			Render("Connect")

	FocusedCheckbox = lipgloss.NewStyle().
			Foreground(ui.AccentColor).
			MarginTop(1).
			Render("TLS Enabled [✕]")

	BlurredCheckbox = lipgloss.NewStyle().
			Foreground(ui.PrimaryColor).
			MarginTop(1).
			Render("TLS Enabled [✕]")

	FocusedCheckboxChecked = lipgloss.NewStyle().
				Foreground(ui.AccentColor).
				MarginTop(1).
				Render("TLS Enabled [🗸]")

	BlurredCheckboxChecked = lipgloss.NewStyle().
				Foreground(ui.PrimaryColor).
				MarginTop(1).
				Render("TLS Enabled [🗸]")

	DialogStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)

	WelcomeMsgStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(ui.PrimaryColor)

	CursorStyle = lipgloss.NewStyle().
			Foreground(ui.AccentColor)
	FocusedStyle = lipgloss.NewStyle().
			Foreground(ui.AccentColor)
	NoStyle = lipgloss.NewStyle()
)
