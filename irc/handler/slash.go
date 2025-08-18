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
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/irc/commands"
)

func handleSlashPrivMsg(params []string, client *irc.Client) tea.Cmd {
	var batchedCmds []tea.Cmd

	if len(params) < 2 {
		client.SendCommand(commands.PRIVMSG, params...)
		return nil
	}

	target := strings.ToLower(params[0])

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
	}
	now := time.Now()
	timestamp := fmt.Sprintf("[%02d:%02d]", now.Hour(), now.Minute())
	msg := fmt.Sprintf("%s: %s", client.Nickname, strings.Join(params[1:], " "))

	c := client.RootChannel
	for {
		if c.Value.Name == target {
			client.ActiveChannel = c

			c.Value.AppendMsg(timestamp, msg, msgOpts)

			client.SendCommand(commands.PRIVMSG, params...)
			return cmds.SwitchChannels
		}

		c = c.Next
		if c == client.RootChannel {
			break
		}
	}

	// If we're messaging a user and their "channel" wasn't found in the previous loop, then create it and append it
	if target[0] != '#' && target[0] != '&' {
		newChannel := irc.Channel{
			Name:  target,
			Users: map[string]irc.User{target: {}},
		}

		newChannel.AppendMsg(timestamp, msg, msgOpts)
		client.ActiveChannel = client.AppendChannel(newChannel)

		batchedCmds = append(batchedCmds, cmds.UpdateTabBar)
	}

	client.SendCommand(commands.PRIVMSG, params...)
	batchedCmds = append(batchedCmds, cmds.SwitchChannels)

	return tea.Batch(batchedCmds...)
}

func handleSlashJoin(params []string, client *irc.Client) tea.Cmd {
	channel := ""
	if len(params) > 0 {
		channel = strings.ToLower(params[0])
	}

	c := client.RootChannel
	for {
		if c.Value.Name == channel {
			client.ActiveChannel = c
			return cmds.SwitchChannels
		}

		c = c.Next
		if c == client.RootChannel {
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
	case commands.PRIVMSG:
		return handleSlashPrivMsg(params, client)
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
