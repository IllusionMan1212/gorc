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

package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/ui/app"
)

func main() {
	if err := tea.NewProgram(app.InitialState(), tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
