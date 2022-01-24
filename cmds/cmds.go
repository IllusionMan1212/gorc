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

package cmds

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/irc"
)

type ConnectMsg struct{}

func Connect() tea.Msg {
	return ConnectMsg{}
}

func Quit(client *irc.Client) tea.Cmd {
	return func() tea.Msg {
		if client.TcpConn != nil {
			client.SendCommand("QUIT")
		}

		return tea.Quit()
	}
}

type SendPrivMsgMsg struct {
	Msg string
}

func SendPrivMsg(msg string) tea.Cmd {
	return func() tea.Msg {
		return SendPrivMsgMsg{Msg: msg}
	}
}

type ReceivedIRCMsgMsg struct{}

func ReceivedIRCMsg() tea.Msg {
	return ReceivedIRCMsgMsg{}
}

type SwitchChannelsMsg struct{}

func SwitchChannels() tea.Msg {
	return SwitchChannelsMsg{}
}

type UpdateTabBarMsg struct{}

func UpdateTabBar() tea.Msg {
	return UpdateTabBarMsg{}
}
