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

package mainscreen

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/ui"
)

var (
	InputboxStyle = ui.DefaultStyle.
			Border(lipgloss.NormalBorder(), true).
			Padding(0, 1).
			Width(ui.MainStyle.GetWidth()).
			BorderForeground(ui.AccentColor)
	MessagesStyle = ui.DefaultStyle.
			Border(lipgloss.NormalBorder(), false, true, true, true).
			Width(ui.MainStyle.GetWidth() * 8 / 10).
			BorderForeground(ui.PrimaryColor)
	SidePanelStyle = ui.DefaultStyle.
			Border(lipgloss.NormalBorder(), true).
			Width(ui.MainStyle.GetWidth() * 2 / 10).
			BorderForeground(ui.PrimaryColor)

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	leftArrowBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "┐",
		BottomLeft:  "│",
		BottomRight: "└",
	}

	rightArrowBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "│",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		Foreground(ui.PrimaryColor).
		Padding(0, 1)

	activeTab = tab.
			Border(activeTabBorder, true).
			BorderForeground(ui.AccentColor).
			Foreground(ui.AccentColor).
			Italic(true).
			Bold(true)

	leftArrowDim = tab.
			Border(leftArrowBorder, true).
			BorderForeground(ui.PrimaryColor).
			Foreground(ui.DisabledColor)
	rightArrowDim = tab.
			Border(rightArrowBorder, true).
			BorderForeground(ui.PrimaryColor).
			Foreground(ui.DisabledColor)

	leftArrowLit = tab.
			Border(leftArrowBorder, true).
			Foreground(ui.PrimaryColor)
	rightArrowLit = tab.
			Border(rightArrowBorder, true).
			Foreground(ui.PrimaryColor)

	tabLine = lipgloss.NewStyle().
		Foreground(ui.PrimaryColor)
)
