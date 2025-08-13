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
	"strconv"
	"strings"
	"time"

	"slices"

	"github.com/illusionman1212/gorc/cmds"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/irc/commands"
	"github.com/illusionman1212/gorc/irc/parser"
)

func ReadLoop(client *irc.Client) {
	// 512 bytes as a base + 8192 additional bytes for tags
	r := bufio.NewReaderSize(client.TCPConn, 8192+512)

	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			client.TCPConn.Close()
			return
		}

		msg = strings.Replace(msg, "\r\n", "", 1)
		if ircMessage, valid := parser.ParseIRCMessage(msg); valid {
			HandleCommand(ircMessage, client)
		}
	}
}

func handlePing(msg irc.Message, client *irc.Client) {
	token := msg.Parameters[0]
	client.SendCommand(commands.PONG, token)
}

func handlePrivMsg(msg irc.Message, client *irc.Client) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	channels := strings.Split(msg.Parameters[0], ",")
	msgContent := msg.Parameters[1]
	privMsg := fmt.Sprintf("%s: %s", nick, msgContent)

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
	}

	for _, channel := range channels {
		for i, c := range client.Channels {
			if c.Name == channel {
				client.Channels[i].AppendMsg(msg.Timestamp, privMsg, msgOpts)
			}
		}
	}
}

func handleNotice(msg irc.Message, client *irc.Client) {
	source := strings.SplitN(msg.Source, "!", 2)[0]
	targets := strings.Split(msg.Parameters[0], ",")
	msgContent := msg.Parameters[1]
	notice := fmt.Sprintf("%s: %s", source, msgContent)

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
	}

	for _, target := range targets {
		if target == "*" || target == client.Nickname {
			client.Channels[0].AppendMsg(msg.Timestamp, notice, msgOpts)
			continue
		}

		for i, c := range client.Channels {
			if c.Name == targets[0] {
				client.Channels[i].AppendMsg(msg.Timestamp, notice, msgOpts)
			}
		}
	}

}

func handleJoin(msg irc.Message, client *irc.Client) {
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

func handleQuit(msg irc.Message, client *irc.Client) {
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

func handlePart(msg irc.Message, client *irc.Client) {
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
				client.Channels = slices.Delete(client.Channels, i, i+1)
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

func handleTopic(msg irc.Message, client *irc.Client) {
	channel := msg.Parameters[0]
	topic := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		AsServerMsg: true,
	}

	for i, c := range client.Channels {
		if c.Name == channel {
			client.Channels[i].Topic = topic
			client.Channels[i].AppendMsg(msg.Timestamp, fmt.Sprintf("Topic changed: %v", topic), msgOpts)
		}
	}
}

func handleCAP(msg irc.Message, client *irc.Client) {
	switch msg.Parameters[1] {
	case "LIST":
		fallthrough
	case "NEW":
		fallthrough
	case "LS":
		caps := strings.Split(msg.Parameters[2], " ")
		recognizedCaps := make([]irc.Capabilities, 0)
		capabilities := make(irc.Capabilities, 0)
		length := 0

		for _, capability := range caps {
			// NOTE: not sure if I should store the cap value the server advertises or use my own
			parts := strings.Split(capability, "=")
			key := parts[0]
			value := ""
			if len(parts) > 1 {
				value = parts[1]
			}

			if commands.Capabilities[key] {
				// When the final parameter approaches 510 bytes,
				// we send multiple REQ commands
				// (https://ircv3.net/specs/extensions/capability-negotiation.html#the-cap-req-subcommand)
				if length+len(capability) >= 500 {
					recognizedCaps = append(recognizedCaps, capabilities)
					capabilities = make(irc.Capabilities, 0)
					length = 0
				}

				capabilities[key] = value
				length += len(capability)
			} else {
				log.Println("Unknown capability:", capability)
			}
		}
		recognizedCaps = append(recognizedCaps, capabilities)

		for _, capsList := range recognizedCaps {
			list := make([]string, 0)
			for key, value := range capsList {
				if len(value) > 0 {
					list = append(list, key+"="+value)
				} else {
					list = append(list, key)
				}
			}
			client.SendCommand(commands.CAP, "REQ", strings.Join(list, " "))
		}
		client.SendCommand(commands.CAP, "END")
	case "ACK":
		caps := strings.Split(msg.Parameters[2], " ")
		for _, capability := range caps {
			if capability[0] == '-' {
				delete(client.EnabledCapabilities, capability[1:])
			} else {
				parts := strings.Split(capability, "=")
				key := parts[0]
				value := ""
				if len(parts) > 1 {
					value = parts[1]
				}
				client.EnabledCapabilities[key] = value
			}
		}
	case "NAK":
		caps := strings.Split(msg.Parameters[2], " ")
		msgOpts := irc.MsgFmtOpts{
			WithTimestamp: true,
			AsServerMsg:   true,
		}
		client.Channels[0].AppendMsg(msg.Timestamp, "Unrecognized capabilities: "+strings.Join(caps, " "), msgOpts)
	case "DEL":
		caps := strings.Split(msg.Parameters[2], " ")
		for _, capability := range caps {
			delete(client.EnabledCapabilities, capability)
		}
	}
}

func handleWELCOME(msg irc.Message, client *irc.Client) {
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

func handleYOURHOST(msg irc.Message, client *irc.Client) {
	host := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, host, msgOpts)
}

func handleCREATED(msg irc.Message, client *irc.Client) {
	created := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, created, msgOpts)
}

func handleMYINFO(msg irc.Message, client *irc.Client) {
	info := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, info, msgOpts)
}

func handleLUSERCLIENT(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLUSEROP(msg irc.Message, client *irc.Client) {
	message := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLUSERUNKNOWN(msg irc.Message, client *irc.Client) {
	message := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLUSERCHANNELS(msg irc.Message, client *irc.Client) {
	message := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLUSERME(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleLOCALUSERS(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleGLOBALUSERS(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleWHOISUSER(msg irc.Message, client *irc.Client) {
	nick := msg.Parameters[1]
	user := msg.Parameters[2]
	host := msg.Parameters[3]
	realName := msg.Parameters[5]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	message := fmt.Sprintf(
		"Nick: %s | User: %s | Host: %s | Real name: %s",
		nick,
		user,
		host,
		realName,
	)

	client.Channels[0].AppendMsg(msg.Timestamp, "WHOIS Information", msgOpts)
	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleWHOISSERVER(msg irc.Message, client *irc.Client) {
	server := msg.Parameters[2]
	serverInfo := msg.Parameters[3]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	message := fmt.Sprintf(
		"server: %s [%s]",
		server,
		serverInfo,
	)

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleWHOISIDLE(msg irc.Message, client *irc.Client) {
	idleSeconds := msg.Parameters[2]
	connectedTimestamp, err := strconv.ParseInt(msg.Parameters[3], 10, 64)

	if err != nil {
		// TODO: return this err and handle it in the parent
		log.Println(err)
	}

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	since := time.Unix(connectedTimestamp, 0)

	message := fmt.Sprintf(
		"idle for: %s seconds, connected since: %s",
		idleSeconds,
		since.Format(time.ANSIC),
	)
	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleENDOFWHOIS(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleWHOISCHANNELS(msg irc.Message, client *irc.Client) {
	chans := msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	message := fmt.Sprintf("channels: %s", chans)

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleNOTOPIC(msg irc.Message, client *irc.Client) {
	channel := msg.Parameters[1]
	msgStr := msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		AsServerMsg: true,
	}

	for i, c := range client.Channels {
		if c.Name == channel {
			client.Channels[i].AppendMsg(msg.Timestamp, msgStr, msgOpts)
		}
	}
}

func handleTOPIC(msg irc.Message, client *irc.Client) {
	channel := msg.Parameters[1]
	topic := msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		AsServerMsg: true,
	}

	for i, c := range client.Channels {
		if c.Name == channel {
			client.Channels[i].Topic = topic
			client.Channels[i].AppendMsg(msg.Timestamp, fmt.Sprintf("TOPIC: %v", topic), msgOpts)
		}
	}
}

func handleNAMREPLY(msg irc.Message, client *irc.Client) {
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

func handleMOTDStart(msg irc.Message, client *irc.Client) {
	message := strings.Join(msg.Parameters[1:], " ")

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleMOTD(msg irc.Message, client *irc.Client) {
	messageLine := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, messageLine, msgOpts)
}

func handleWHOISHOST(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1] + " " + msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleWHOISMODES(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1] + " " + msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleNOSUCHSERVER(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1] + " " + msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleUNKNOWNCOMMAND(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1] + " " + msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleNONICKNAMEGIVEN(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleNEEDMOREPARAMS(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1] + " " + msg.Parameters[2]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func handleALREADYREGISTERED(msg irc.Message, client *irc.Client) {
	message := msg.Parameters[1]

	msgOpts := irc.MsgFmtOpts{
		WithTimestamp: true,
		AsServerMsg:   true,
	}

	client.Channels[0].AppendMsg(msg.Timestamp, message, msgOpts)
}

func HandleCommand(msg irc.Message, client *irc.Client) {
	// TODO: handle different commands
	switch msg.Command {
	case commands.PING:
		handlePing(msg, client)
	case commands.PRIVMSG:
		handlePrivMsg(msg, client)
	case commands.NOTICE:
		handleNotice(msg, client)
	case commands.JOIN:
		handleJoin(msg, client)
	case commands.QUIT:
		handleQuit(msg, client)
	case commands.PART:
		handlePart(msg, client)
	case commands.TOPIC:
		handleTopic(msg, client)
	case commands.CAP:
		handleCAP(msg, client)
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
	case commands.RPL_WHOISUSER:
		handleWHOISUSER(msg, client)
	case commands.RPL_WHOISSERVER:
		handleWHOISSERVER(msg, client)
	case commands.RPL_WHOISIDLE:
		handleWHOISIDLE(msg, client)
	case commands.RPL_ENDOFWHOIS:
		handleENDOFWHOIS(msg, client)
	case commands.RPL_WHOISCHANNELS:
		handleWHOISCHANNELS(msg, client)
	case commands.RPL_NOTOPIC:
		handleNOTOPIC(msg, client)
	case commands.RPL_TOPIC:
		handleTOPIC(msg, client)
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
	case commands.RPL_WHOISHOST:
		handleWHOISHOST(msg, client)
	case commands.RPL_WHOISMODES:
		handleWHOISMODES(msg, client)
	case commands.ERR_NOSUCHSERVER:
		handleNOSUCHSERVER(msg, client)
	case commands.ERR_UNKNOWNCOMMAND:
		handleUNKNOWNCOMMAND(msg, client)
	case commands.ERR_NONICKNAMEGIVEN:
		handleNONICKNAMEGIVEN(msg, client)
	case commands.ERR_NEEDMOREPARAMS:
		handleNEEDMOREPARAMS(msg, client)
	case commands.ERR_ALREADYREGISTERED:
		handleALREADYREGISTERED(msg, client)
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
