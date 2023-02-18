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

package irc

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/illusionman1212/gorc/db"
	"github.com/illusionman1212/gorc/irc/commands"
	"github.com/illusionman1212/gorc/ui"
	"gorm.io/gorm"
)

type User struct {
	// User prefix in channel
	Prefix string
}

type Channel struct {
	// Channel name
	Name string

	// Channel topic
	Topic string

	// Scrollback buffer for the channel message history
	ScrollBackBuf []string

	// Whethere the viewport is at the bottom or not
	AtBottom bool

	// Total amount of lines that have been received from the tcp conn
	// and written to the history db
	TotalLines uint64

	// Buffer of messages that need to be inserted into the db
	ToInsert []db.History

	// Position/Index of first line in the scrollback buffer as mapped
	// to TotalLines.
	// Last line is ScrollBackBufFirstLine + SCROLLBACK_MAX - 1
	ScrollBackBufFirstLine uint64

	// Unique sha1 hash for the channel made from server hostname and channel name
	Hash string

	// Users in this channel
	// The map key is the user's nickname
	// and the user struct holds data about that user
	// such as, prefixes in this channel and etc...
	Users map[string]User
}

type Client struct {
	// tcp connection
	TcpConn net.Conn

	// Host that client is connected to
	Host string

	// Port that client is connected to
	Port string

	// Nickname currently in use by the user
	Nickname string

	// Joined channels
	Channels []Channel

	// Active channel name
	ActiveChannel string

	// Active channel index
	ActiveChannelIndex int

	// Reference to the bubbletea program
	Tea *tea.Program

	// TODO: remove this
	Ticks int

	// DB reference
	DB *gorm.DB

	// Index of the first visible tab in the tab bar
	FirstTabIndexInTabBar int

	// Index of the last visible tab in the tab bar
	LastTabIndexInTabBar int
}

type MsgFmtOpts struct {
	WithTimestamp bool
	AsServerMsg   bool
	NotImpl       bool
	AsDate        bool
}

var notImpl = lipgloss.NewStyle().Foreground(ui.ErrorColor).Render("[NOT IMPL]")
var serverMsgStyle = lipgloss.NewStyle().Foreground(ui.ServerMsgColor)
var timestampStyle = serverMsgStyle.Copy()
var dateStyle = lipgloss.NewStyle().Foreground(ui.DateColor)

const (
	CRLF           = "\r\n"
	SCROLLBACK_MAX = 400
)

func (c *Channel) AppendMsg(host string, database *gorm.DB, timestamp string, fullMsg string, opts MsgFmtOpts) {
	prefixes := ""
	style := ui.DefaultStyle

	if opts.WithTimestamp {
		prefixes += timestampStyle.Render(timestamp) + " "
	}

	if opts.NotImpl {
		prefixes += notImpl + " "
	}

	if opts.AsServerMsg {
		prefixes += serverMsgStyle.Render("==") + " "
		style = serverMsgStyle
	}

	if opts.AsDate {
		style = dateStyle
	}

	finalMsg := prefixes + style.Render(fullMsg)

	c.TotalLines++

	history := db.History{Hash: c.Hash, Line: c.TotalLines, Message: finalMsg}
	c.ToInsert = append(c.ToInsert, history)

	if c.TotalLines%10 == 0 {
		database.Create(c.ToInsert)
		c.ToInsert = nil
	}

	if len(c.ScrollBackBuf) < SCROLLBACK_MAX {
		c.ScrollBackBuf = append(c.ScrollBackBuf, prefixes+style.Render(fullMsg))
		// We haven't reached the bottom yet so the first line is still 1
		c.ScrollBackBufFirstLine = 1
	} else if c.AtBottom && len(c.ScrollBackBuf) >= SCROLLBACK_MAX {
		// TODO: when pressing 'g' in viewport, we need to set this to 1
		// TODO: when pressing 'G' in viewport, we need to set this to TotalLines - 1000 + 1
		// TODO: when going up in viewport, we need to decrement the amount of lines by how many lines we need to prepend to the scrollbackbuf
		// TODO: when going down in viewport, we need to increment the amount of lines by how many lines we need to append to the scrollbackbuf
		// If we're at bottom and we're appending a line, then we increment the first line
		c.ScrollBackBufFirstLine++
		c.ScrollBackBuf = append(c.ScrollBackBuf, prefixes+style.Render(fullMsg))
		c.ScrollBackBuf = c.ScrollBackBuf[1:]
	}
}

func (s *Client) Initialize(host string, port string, tlsEnabled bool) {
	addr := fmt.Sprintf("%s:%s", host, port)

	if tlsEnabled {
		cfg := &tls.Config{ServerName: host}
		conn, err := tls.Dial("tcp", addr, cfg)
		if err != nil {
			// TODO: properly handle the error instead of Fatal-ing (failed to initiate a connection to server)
			log.Fatal(err)
		}

		s.TcpConn = conn
		s.Host = host
		s.Port = port
		s.ActiveChannel = host
		s.Channels = make([]Channel, 0)
		s.Ticks = 1
		return
	}

	addrTCP, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, addrTCP)
	if err != nil {
		log.Fatal(err)
	}

	s.TcpConn = conn
	s.Host = host
	s.Port = port
	s.ActiveChannel = host
	s.Channels = make([]Channel, 0)
	s.Ticks = 1
}

func (c *Client) Register(nick string, password string, channel string) {
	c.Channels = append(c.Channels, Channel{
		Name:  c.Host,
		Users: make(map[string]User),
	})

	c.SendCommand(commands.CAP, "LS")
	if password != "" {
		c.SendCommand(commands.PASS, password)
	}
	// TODO: check if nickname has spaces and remove them
	c.SendCommand(commands.NICK, nick)
	// set user-wanted nickname
	c.Nickname = nick
	c.SendCommand(commands.USER, nick, "0", "*", nick)
	// TODO: CAP REQ :whatever capability the client recognizes and supports
	c.SendCommand(commands.CAP, "END")
	// joining a channel when registering is optional
	if channel != "" {
		c.SendCommand(commands.JOIN, channel)
	}
}

func (c *Client) SetDay() {
	msgOpts := MsgFmtOpts{
		AsDate: true,
	}

	now := time.Now()
	msg := fmt.Sprintf("————— %s %d —————", now.Month().String(), now.Day())
	c.Channels[0].AppendMsg(c.Host, c.DB, "", msg, msgOpts)

	// TODO: append new day to each channel whenever the day changes
}

func (c Client) SendCommand(cmd string, params ...string) {
	if c.TcpConn == nil {
		// TODO: properly handle the error instead of Fatal-ing
		log.Fatal("Attempted to write data to nil connection")
	}

	paramsString := ""

	// if we have more than 1 param
	if len(params) > 1 {
		lastParam := params[len(params)-1]
		// if the last param is a trailing param
		// we prepend a colon to it
		if strings.Contains(lastParam, " ") {
			lastParam = " :" + lastParam
		} else {
			lastParam = " " + lastParam
		}

		paramsString = " " + strings.Join(params[:len(params)-1], " ")
		paramsString += lastParam

		// if we have exactly 1 param and it's not empty
	} else if len(params) == 1 && params[0] != "" {
		// if this 1 param contains spaces, we prepend a colon
		if strings.Contains(params[0], " ") {
			paramsString = " :" + params[0]
		} else {
			paramsString = " " + params[0]
		}
	}

	c.TcpConn.Write([]byte(cmd + paramsString + CRLF))
}
