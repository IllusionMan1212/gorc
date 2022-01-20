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

package mainscreen

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/ui"
)

var (
	InputboxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			Padding(0, 1).
			Width(ui.MainStyle.GetWidth()).
			BorderForeground(lipgloss.Color("105"))
	MessagesStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, true, true).
			Width(ui.MainStyle.GetWidth() * 8 / 10).
			BorderForeground(lipgloss.Color("#EEE"))
	SidePanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true).
			Width(ui.MainStyle.GetWidth() * 2 / 10).
			BorderForeground(lipgloss.Color("#EEE"))

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
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "105"}
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
		Padding(0, 1)

	activeTab = tab.Copy().
			Border(activeTabBorder, true).
			BorderForeground(highlight).
			Foreground(lipgloss.Color("105")).
			Italic(true).
			Bold(true)

	leftArrowDim = tab.Copy().
			Border(leftArrowBorder, true).
			Foreground(lipgloss.Color("#444"))
	rightArrowDim = tab.Copy().
			Border(rightArrowBorder, true).
			Foreground(lipgloss.Color("#444"))

	leftArrowLit = tab.Copy().
			Border(leftArrowBorder, true).
			Foreground(lipgloss.Color("#FFF"))
	rightArrowLit = tab.Copy().
			Border(rightArrowBorder, true).
			Foreground(lipgloss.Color("#FFF"))

	tabLine = lipgloss.NewStyle()
)
