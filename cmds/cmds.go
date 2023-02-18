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
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/illusionman1212/gorc/irc"
	"github.com/illusionman1212/gorc/irc/commands"
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

type TickMsg struct{}

func Tick(client *irc.Client) tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		client.SendCommand(commands.PRIVMSG, "#hell", fmt.Sprintf("%d", client.Ticks))
		return TickMsg{}
	})
}

type SendPrivMsgMsg struct {
	Msg       string
	Timestamp string
}

func SendPrivMsg(msg string) tea.Cmd {
	return func() tea.Msg {
		now := time.Now()
		timestamp := fmt.Sprintf("[%02d:%02d]", now.Hour(), now.Minute())

		return SendPrivMsgMsg{
			Msg:       msg,
			Timestamp: timestamp,
		}
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

type ReadHistoryFromDBMsg struct {
	LinesRead uint64
	Direction string
}

func ReadHistoryFromDB(client *irc.Client, linesToRead int, direction string) tea.Cmd {
	hash := client.Channels[client.ActiveChannelIndex].Hash
	scrollBackBufFirstLine := client.Channels[client.ActiveChannelIndex].ScrollBackBufFirstLine

	messages := make([]string, 0)

	switch direction {
	case "top":
		client.DB.
			Table("histories").
			Select("message").
			Where("hash = ?", hash).
			Limit(linesToRead).
			Scan(&messages).
			Order("line ASC")

		// set the new messages
		copy(client.Channels[client.ActiveChannelIndex].ScrollBackBuf, messages)
		// adjust the scrollback buf first line
		client.Channels[client.ActiveChannelIndex].ScrollBackBufFirstLine = 1

		return func() tea.Msg {
			return ReadHistoryFromDBMsg{
				LinesRead: 0,
				Direction: "top",
			}
		}
	case "up":
		limit := linesToRead
		linesToGet := 0

		if linesToRead > int(scrollBackBufFirstLine) {
			limit = int(scrollBackBufFirstLine) - 1
			linesToGet = int(scrollBackBufFirstLine) - limit
		} else {
			linesToGet = int(scrollBackBufFirstLine) - linesToRead
		}

		client.DB.
			Table("histories").
			Select("message").
			Where("hash = ? AND line >= ?", hash, linesToGet).
			Limit(limit).
			Scan(&messages).
			Order("line ASC")

		if len(messages) != 0 {
			// prepend to the scroll back buf
			client.Channels[client.ActiveChannelIndex].ScrollBackBuf = append(messages, client.Channels[client.ActiveChannelIndex].ScrollBackBuf...)
			client.Channels[client.ActiveChannelIndex].ScrollBackBuf = client.Channels[client.ActiveChannelIndex].ScrollBackBuf[:irc.SCROLLBACK_MAX]
			// adjust the scrollback buf first line
			client.Channels[client.ActiveChannelIndex].ScrollBackBufFirstLine -= uint64(limit)

			return func() tea.Msg {
				return ReadHistoryFromDBMsg{
					LinesRead: uint64(limit),
					Direction: "up",
				}
			}
		}
	case "down":
		client.DB.
			Table("histories").
			Select("message").
			Where("hash = ? AND line >= ?", hash, client.Channels[client.ActiveChannelIndex].ScrollBackBufFirstLine+irc.SCROLLBACK_MAX).
			Limit(linesToRead).
			Scan(&messages).
			Order("line ASC")
		linesGot := len(messages)
		if len(messages) != 0 {
			extraMessagesLen := 0
			// append to the scroll back buf
			client.Channels[client.ActiveChannelIndex].ScrollBackBuf = append(client.Channels[client.ActiveChannelIndex].ScrollBackBuf, messages...)
			// if we reached the end of the messages in the db. try appending the messages in the ToInsert buffer
			if linesGot < linesToRead {
				for _, history := range client.Channels[client.ActiveChannelIndex].ToInsert {
					client.Channels[client.ActiveChannelIndex].ScrollBackBuf = append(client.Channels[client.ActiveChannelIndex].ScrollBackBuf, history.Message)
				}
				extraMessagesLen = len(client.Channels[client.ActiveChannelIndex].ToInsert)
			}
			client.Channels[client.ActiveChannelIndex].ScrollBackBuf = client.Channels[client.ActiveChannelIndex].ScrollBackBuf[len(client.Channels[client.ActiveChannelIndex].ScrollBackBuf)-irc.SCROLLBACK_MAX:]
			// adjust the scrollback buf first line
			client.Channels[client.ActiveChannelIndex].ScrollBackBufFirstLine += uint64(linesGot) + uint64(extraMessagesLen)

			return func() tea.Msg {
				return ReadHistoryFromDBMsg{
					LinesRead: uint64(linesGot) + uint64(extraMessagesLen),
					Direction: "down",
				}
			}
		}
	case "bottom":
		client.DB.Raw("SELECT t.message FROM (SELECT line, message FROM histories WHERE hash = ? ORDER BY line DESC LIMIT ?) t ORDER BY t.line ASC;", hash, linesToRead).
			Scan(&messages)

		toInsert := make([]string, 0)

		for _, history := range client.Channels[client.ActiveChannelIndex].ToInsert {
			toInsert = append(toInsert, history.Message)
		}

		// set the new messages
		copy(client.Channels[client.ActiveChannelIndex].ScrollBackBuf, messages)
		client.Channels[client.ActiveChannelIndex].ScrollBackBuf = append(client.Channels[client.ActiveChannelIndex].ScrollBackBuf, toInsert...)
		client.Channels[client.ActiveChannelIndex].ScrollBackBuf = client.Channels[client.ActiveChannelIndex].ScrollBackBuf[len(toInsert):]
		// adjust the scrollback buf first line
		client.Channels[client.ActiveChannelIndex].ScrollBackBufFirstLine = client.Channels[client.ActiveChannelIndex].TotalLines - irc.SCROLLBACK_MAX + uint64(len(toInsert))

		return func() tea.Msg {
			return ReadHistoryFromDBMsg{
				LinesRead: 0,
				Direction: "bottom",
			}
		}
	}

	return func() tea.Msg {
		return ReadHistoryFromDBMsg{
			LinesRead: 0,
			Direction: "unknown",
		}
	}
}
