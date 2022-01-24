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

package handler

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/irc/commands"
)

func handleSlashJoin(params []string, client *irc.Client) tea.Cmd {
	channel := ""
	if len(params) > 0 {
		channel = params[0]
	}

	for i, c := range client.Channels {
		if c.Name == channel {
			client.ActiveChannelIndex = i
			client.ActiveChannel = channel
			break
		}
	}
	client.SendCommand(commands.JOIN, params...)

	return cmds.SwitchChannels
}

func HandleSlashCommand(msg string, client *irc.Client) tea.Cmd {
	substrs := strings.Fields(msg[1:])
	command := strings.ToUpper(substrs[0])
	var params []string
	if len(substrs) > 1 {
		params = substrs[1:]
	}

	switch command {
	case commands.JOIN:
		return handleSlashJoin(params, client)
	case commands.QUIT:
		client.SendCommand(command, params...)
		return cmds.Quit(client)
	default:
		client.SendCommand(command, params...)
		return nil
	}
}
