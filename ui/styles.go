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

package ui

import "github.com/charmbracelet/lipgloss"

var (
	BorderRound = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	MainStyle = lipgloss.NewStyle().
			Width(172).
			Height(40)

	DialogStyle = lipgloss.NewStyle().
			Border(BorderRound, true).
			Align(lipgloss.Center).
			Width(50).
			Padding(3, 6).
			BorderForeground(lipgloss.Color("105"))

	CursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("105"))
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("105"))
	NoStyle      = lipgloss.NewStyle()
)
