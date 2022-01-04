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

import tea "github.com/charmbracelet/bubbletea"

func (s *State) Quit() tea.Msg {
	if s.Client.TlsConn != nil {
		s.Client.SendCommand("QUIT")
		s.Client.TlsConn.Close()
	}
	if s.Client.TcpConn != nil {
		s.Client.SendCommand("QUIT")
		s.Client.TcpConn.Close()
	}

	return tea.Quit()
}