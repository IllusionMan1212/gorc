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
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/irc/commands"
	"github.com/illusionman1212/gorc/irc/parser"
)

func ReadLoop(client *irc.Client) {
	// 512 bytes as a base + 8192 additional bytes for tags
	r := bufio.NewReaderSize(client.TcpConn, 8192+512)

	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			client.TcpConn.Close()
			return
		}

		msg = strings.Replace(msg, "\r\n", "", 1)
		if ircMessage, valid := parser.ParseIRCMessage(msg); valid {
			HandleCommand(ircMessage, client)
		}
	}
}

func handlePing(msg parser.IRCMessage, client *irc.Client) {
	token := msg.Parameters[0]
	client.SendCommand(commands.PONG, token)
}

func handlePrivMsg(msg parser.IRCMessage, client *irc.Client) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	channel := msg.Parameters[0]
	msgContent := msg.Parameters[1]
	privMsg := fmt.Sprintf("%s: %s", nick, msgContent)

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
	}

	for i, c := range client.Channels {
		if c.Name == channel {
			client.Channels[i].AppendMsg(msg.Timestamp, privMsg, msgOpts)
		}
	}
}

func handleJoin(msg parser.IRCMessage, client *irc.Client) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	channel := msg.Parameters[0]

	joinMsg := fmt.Sprintf("%s has joined", nick)

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	if nick == client.Nickname {
		client.Channels = append(client.Channels, irc.Channel{
			Name:  channel,
			Users: make(map[string]irc.User),
		})
		client.ActiveChannelIndex = len(client.Channels) - 1
		client.ActiveChannel = channel
		client.Channels[client.ActiveChannelIndex].AppendMsg(msg.Timestamp, joinMsg, msgOpts)
		if _, exists := client.Channels[client.ActiveChannelIndex].Users[nick]; !exists {
			client.Channels[client.ActiveChannelIndex].Users[nick] = irc.User{}
		}
		client.Tea.Send(cmds.UpdateTabBar())
	} else {
		for i, c := range client.Channels {
			if c.Name == channel {
				client.Channels[i].AppendMsg(msg.Timestamp, joinMsg, msgOpts)
				if _, exists := client.Channels[i].Users[nick]; !exists {
					client.Channels[i].Users[nick] = irc.User{}
				}
			}
		}
	}

	if channel == client.ActiveChannel {
		client.Tea.Send(cmds.SwitchChannels())
	}
}

func handleQuit(msg parser.IRCMessage, client *irc.Client) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	reason := msg.Parameters[0]
	quitMsg := fmt.Sprintf("%s has quit (%s)", nick, reason)

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	for i := range client.Channels {
		// skip the server "channel"
		if i == 0 {
			continue
		}
		client.Channels[i].AppendMsg(msg.Timestamp, quitMsg, msgOpts)
		delete(client.Channels[i].Users, nick)
	}

	client.Tea.Send(cmds.SwitchChannels())
}

func handlePart(msg parser.IRCMessage, client *irc.Client) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	channel := msg.Parameters[0]
	reason := ""
	if len(msg.Parameters) > 1 {
		reason = msg.Parameters[1]
	}

	partMsg := fmt.Sprintf("%s has left %s (%s)", nick, channel, reason)

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	for i, c := range client.Channels {
		if c.Name == channel {
			if nick == client.Nickname {
				client.Channels = append(client.Channels[:i], client.Channels[i+1:]...)
				if client.ActiveChannelIndex >= i {
					client.ActiveChannelIndex--
					client.ActiveChannel = client.Channels[client.ActiveChannelIndex].Name
				}
				client.LastTabIndexInTabBar--
			} else {
				client.Channels[i].AppendMsg(msg.Timestamp, partMsg, msgOpts)
				delete(client.Channels[i].Users, nick)
			}
		}
	}

	client.Tea.Send(cmds.SwitchChannels())
}

func handleWELCOME(msg parser.IRCMessage, client *irc.Client) {
	nick := msg.Parameters[0]
	welcomeMsg := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	// set server-registered nickname
	client.Nickname = nick
	client.Channels[0].AppendMsg(msg.Timestamp, welcomeMsg, msgOpts)
}

func handleYOURHOST(msg parser.IRCMessage, client *irc.Client) {
	host := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, host, msgOpts)
}

func handleCREATED(msg parser.IRCMessage, client *irc.Client) {
	created := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, created, msgOpts)
}

func handleMYINFO(msg parser.IRCMessage, client *irc.Client) {
	info := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, info, msgOpts)
}

func handleLUSERCLIENT(msg parser.IRCMessage, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLUSEROP(msg parser.IRCMessage, client *irc.Client) {
	message := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLUSERUNKNOWN(msg parser.IRCMessage, client *irc.Client) {
	message := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLUSERCHANNELS(msg parser.IRCMessage, client *irc.Client) {
	message := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLUSERME(msg parser.IRCMessage, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLOCALUSERS(msg parser.IRCMessage, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleGLOBALUSERS(msg parser.IRCMessage, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleNAMREPLY(msg parser.IRCMessage, client *irc.Client) {
	// TODO: do i need these
	// client := msg.Parameters[0]
	// chanSymbol := msg.Parameters[1]
	channel := msg.Parameters[2]
	nicks := strings.Split(msg.Parameters[3], " ")

	for i, c := range client.Channels {
		if c.Name == channel {
			for _, nick := range nicks {
				prefix := ""
				_nick := nick

				if commands.UserPrefixes[string(nick[0])] {
					prefix = string(nick[0])
					_nick = nick[1:]
				}

				client.Channels[i].Users[_nick] = irc.User{
					Prefix: prefix,
				}
			}
		}
	}

	if channel == client.ActiveChannel {
		client.Tea.Send(cmds.SwitchChannels())
	}
}

func handleMOTDStart(msg parser.IRCMessage, client *irc.Client) {
	message := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleMOTD(msg parser.IRCMessage, client *irc.Client) {
	messageLine := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, messageLine, msgOpts)
}

func HandleCommand(msg parser.IRCMessage, client *irc.Client) {
	// TODO: handle different commands
	switch msg.Command {
	case commands.PING:
		handlePing(msg, client)
	case commands.PRIVMSG:
		handlePrivMsg(msg, client)
	case commands.JOIN:
		handleJoin(msg, client)
	case commands.QUIT:
		handleQuit(msg, client)
	case commands.PART:
		handlePart(msg, client)
	case commands.RPL_WELCOME:
		handleWELCOME(msg, client)
	case commands.RPL_YOURHOST:
		handleYOURHOST(msg, client)
	case commands.RPL_CREATED:
		handleCREATED(msg, client)
	case commands.RPL_MYINFO:
		handleMYINFO(msg, client)
	case commands.RPL_LUSERCLIENT:
		handleLUSERCLIENT(msg, client)
	case commands.RPL_LUSEROP:
		handleLUSEROP(msg, client)
	case commands.RPL_LUSERUNKNOWN:
		handleLUSERUNKNOWN(msg, client)
	case commands.RPL_LUSERCHANNELS:
		handleLUSERCHANNELS(msg, client)
	case commands.RPL_LUSERME:
		handleLUSERME(msg, client)
	case commands.RPL_LOCALUSERS:
		handleLOCALUSERS(msg, client)
	case commands.RPL_GLOBALUSERS:
		handleGLOBALUSERS(msg, client)
	case commands.RPL_NAMREPLY:
		handleNAMREPLY(msg, client)
	case commands.RPL_ENDOFNAMES:
		// TODO: toggle a flag on the channel to indicate
		// it received all the names correctly.

		// start a timeout and update said timeout on every RPL_NAMREPLY
		// and log an error if timeout ends without receiving this command.
	case commands.RPL_MOTDSTART:
		handleMOTDStart(msg, client)
	case commands.RPL_MOTD:
		handleMOTD(msg, client)
	case commands.RPL_ENDOFMOTD:
		// TODO: toggle a flag on the client/server to indicate
		// it received all the MOTD correctly

		// start a timeout and update said timeout on every RPL_MOTD
		// and log an error if timeout ends without receiving this command.
	default:
		fullMsg := fmt.Sprintf(
			"%s %s %s %s",
			msg.Tags,
			msg.Source,
			msg.Command,
			strings.Join(msg.Parameters, " "),
		)

		msgOpts := irc.MsgFmtOpts{
			WithTimestamp: true,
			NotImpl:       true,
		}

		client.Channels[0].AppendMsg(msg.Timestamp, fullMsg, msgOpts)
	}

	// send a receivedIRCmsg tea message so the ui can update
	// we also use this tea message to scroll the viewport down
	client.Tea.Send(cmds.ReceivedIRCMsg())
}
