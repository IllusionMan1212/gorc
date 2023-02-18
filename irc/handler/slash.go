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
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/irc/commands"
)

func HandleSlashCommand(msg string, client *irc.Client) tea.Cmd {
	substrs := strings.Fields(msg[1:])
	command := strings.ToUpper(substrs[0])
	var params []string
	if len(substrs) > 1 {
		params = substrs[1:]
	}

	switch command {
	case commands.JOIN:
		client.SendCommand(commands.JOIN, params...)
		return cmds.SwitchChannels
	case commands.QUIT:
		client.SendCommand(command, params...)
		return cmds.Quit(client)
	case "TEST":
		log.Println("running the spammer")
		return cmds.Tick(client)
	default:
		client.SendCommand(command, params...)
		return nil
	}
}
