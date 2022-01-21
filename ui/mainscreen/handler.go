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
	"fmt"
	"strings"

	"github.com/illusionman1212/gorc/client"
	"github.com/illusionman1212/gorc/commands"
	"github.com/illusionman1212/gorc/parser"
)

func (s *State) handlePing(msg parser.IRCMessage) {
	token := msg.Parameters[0]
	s.Client.SendCommand(commands.PONG, token)
}

func (s *State) handlePrivMsg(msg parser.IRCMessage) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	channel := msg.Parameters[0]
	msgContent := msg.Parameters[1]
	privMsg := fmt.Sprintf("%s: %s", nick, msgContent)

	for i, c := range s.Client.Channels {
		if c.Name == channel {
			s.Client.Channels[i].History += privMsg + client.CRLF
		}
	}
}

func (s *State) handleJoin(msg parser.IRCMessage) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	channel := msg.Parameters[0]

	joinMsg := fmt.Sprintf("== %s has joined", nick)
	for i, c := range s.Client.Channels {
		if c.Name == channel {
			s.Client.Channels[i].History += joinMsg + client.CRLF
			if _, exists := s.Client.Channels[i].Users[nick]; !exists {
				s.Client.Channels[i].Users[nick] = client.User{}
			}
		}
	}

	if channel == s.Client.ActiveChannel {
		s.SidePanel.UpdateNicks()
	}
}

func (s *State) handleQuit(msg parser.IRCMessage) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	reason := msg.Parameters[0]
	quitMsg := fmt.Sprintf("== %s has quit (%s)", nick, reason)

	for i := range s.Client.Channels {
		// skip the server "channel"
		if i == 0 {
			continue
		}
		s.Client.Channels[i].History += quitMsg + client.CRLF
		delete(s.Client.Channels[i].Users, nick)
	}

	s.SidePanel.UpdateNicks()
}

func (s *State) handlePart(msg parser.IRCMessage) {
	nick := strings.SplitN(msg.Source, "!", 2)[0]
	channel := msg.Parameters[0]
	reason := ""
	if len(msg.Parameters) > 1 {
		reason = msg.Parameters[1]
	}

	partMsg := fmt.Sprintf("== %s has left %s (%s)", nick, channel, reason)

	for i, c := range s.Client.Channels {
		if c.Name == channel {
			s.Client.Channels[i].History += partMsg + client.CRLF
			delete(s.Client.Channels[i].Users, nick)
		}
	}

	s.SidePanel.UpdateNicks()
}

func (s *State) handleWELCOME(msg parser.IRCMessage) {
	welcomeMsg := msg.Parameters[1]

	s.Client.Channels[0].History += welcomeMsg + client.CRLF
}

func (s *State) handleYOURHOST(msg parser.IRCMessage) {
	host := msg.Parameters[1]

	s.Client.Channels[0].History += host + client.CRLF
}

func (s *State) handleCREATED(msg parser.IRCMessage) {
	created := msg.Parameters[1]

	s.Client.Channels[0].History += created + client.CRLF
}

func (s *State) handleMYINFO(msg parser.IRCMessage) {
	info := strings.Join(msg.Parameters[1:], " ")

	s.Client.Channels[0].History += info + client.CRLF
}

func (s *State) handleLUSERCLIENT(msg parser.IRCMessage) {
	message := msg.Parameters[1]

	s.Client.Channels[0].History += message + client.CRLF
}

func (s *State) handleLUSEROP(msg parser.IRCMessage) {
	message := strings.Join(msg.Parameters[1:], " ")

	s.Client.Channels[0].History += message + client.CRLF
}

func (s *State) handleLUSERUNKNOWN(msg parser.IRCMessage) {
	message := strings.Join(msg.Parameters[1:], " ")

	s.Client.Channels[0].History += message + client.CRLF
}

func (s *State) handleLUSERCHANNELS(msg parser.IRCMessage) {
	message := strings.Join(msg.Parameters[1:], " ")

	s.Client.Channels[0].History += message + client.CRLF
}

func (s *State) handleLUSERME(msg parser.IRCMessage) {
	message := msg.Parameters[1]

	s.Client.Channels[0].History += message + client.CRLF
}

func (s *State) handleLOCALUSERS(msg parser.IRCMessage) {
	message := msg.Parameters[1]

	s.Client.Channels[0].History += message + client.CRLF
}

func (s *State) handleGLOBALUSERS(msg parser.IRCMessage) {
	message := msg.Parameters[1]

	s.Client.Channels[0].History += message + client.CRLF
}

func (s *State) handleNAMREPLY(msg parser.IRCMessage) {
	// TODO: do i need these
	// client := msg.Parameters[0]
	// chanSymbol := msg.Parameters[1]
	channel := msg.Parameters[2]
	nicks := strings.Split(msg.Parameters[3], " ")

	for i, c := range s.Client.Channels {
		if c.Name == channel {
			for _, nick := range nicks {
				prefix := ""
				_nick := nick

				if commands.UserPrefixes[string(nick[0])] {
					prefix = string(nick[0])
					_nick = nick[1:]
				}

				s.Client.Channels[i].Users[_nick] = client.User{
					Prefix: prefix,
				}
			}
		}
	}

	if channel == s.Client.ActiveChannel {
		s.SidePanel.UpdateNicks()
	}
}

func (s *State) handleMOTDStart(msg parser.IRCMessage) {
	params := strings.Join(msg.Parameters[1:], " ")

	s.Client.Channels[0].History += params + client.CRLF
}

func (s *State) handleMOTD(msg parser.IRCMessage) {
	messageLine := msg.Parameters[1]

	s.Client.Channels[0].History += messageLine + client.CRLF
}

func (s *State) HandleCommand(msg parser.IRCMessage) {
	// TODO: handle different commands
	switch msg.Command {
	case commands.PING:
		s.handlePing(msg)
	case commands.PRIVMSG:
		s.handlePrivMsg(msg)
	case commands.JOIN:
		s.handleJoin(msg)
	case commands.QUIT:
		s.handleQuit(msg)
	case commands.PART:
		s.handlePart(msg)
	case commands.RPL_WELCOME:
		s.handleWELCOME(msg)
	case commands.RPL_YOURHOST:
		s.handleYOURHOST(msg)
	case commands.RPL_CREATED:
		s.handleCREATED(msg)
	case commands.RPL_MYINFO:
		s.handleMYINFO(msg)
	case commands.RPL_LUSERCLIENT:
		s.handleLUSERCLIENT(msg)
	case commands.RPL_LUSEROP:
		s.handleLUSEROP(msg)
	case commands.RPL_LUSERUNKNOWN:
		s.handleLUSERUNKNOWN(msg)
	case commands.RPL_LUSERCHANNELS:
		s.handleLUSERCHANNELS(msg)
	case commands.RPL_LUSERME:
		s.handleLUSERME(msg)
	case commands.RPL_LOCALUSERS:
		s.handleLOCALUSERS(msg)
	case commands.RPL_GLOBALUSERS:
		s.handleGLOBALUSERS(msg)
	case commands.RPL_NAMREPLY:
		s.handleNAMREPLY(msg)
	case commands.RPL_ENDOFNAMES:
		// TODO: toggle a flag on the channel to indicate
		// it received all the names correctly.

		// start a timeout and update said timeout on every RPL_NAMREPLY
		// and log an error if timeout ends without receiving this command.
	case commands.RPL_MOTDSTART:
		s.handleMOTDStart(msg)
	case commands.RPL_MOTD:
		s.handleMOTD(msg)
	case commands.RPL_ENDOFMOTD:
		// TODO: toggle a flag on the client/server to indicate
		// it received all the MOTD correctly

		// start a timeout and update said timeout on every RPL_MOTD
		// and log an error if timeout ends without receiving this command.
	default:
		fullMsg := fmt.Sprintf("%s %s %s %s", msg.Tags, msg.Source, msg.Command, strings.Join(msg.Parameters, " "))

		s.Client.Channels[0].History += fullMsg + client.CRLF
	}

	// send a receivedIRCmsg tea message so the ui can update
	// we also use this tea message to scroll the viewport down
	s.Client.Tea.Send(ReceivedIRCMsg())
}
