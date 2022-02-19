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

package ui

import "github.com/charmbracelet/lipgloss"

var (
	MainStyle = lipgloss.NewStyle().
			Width(0).
			Height(0)

	DefaultStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor)

	AccentColor = lipgloss.AdaptiveColor{
		Dark:  "#1C6DD0",
		Light: "#1C6DD0",
	}

	PrimaryColor = lipgloss.AdaptiveColor{
		Dark:  "#EEE",
		Light: "#151515",
	}

	DisabledColor = lipgloss.AdaptiveColor{
		Dark:  "#444",
		Light: "#AAA",
	}

	DisabledColorFocus = lipgloss.AdaptiveColor{
		Dark:  "#666",
		Light: "#CCC",
	}

	ErrorColor = lipgloss.Color("#F00")
	DateColor  = lipgloss.AdaptiveColor{
		Dark:  "#0F0",
		Light: "#090",
	}
	ServerMsgColor = lipgloss.AdaptiveColor{
		Dark:  "#999",
		Light: "#777",
	}
)
