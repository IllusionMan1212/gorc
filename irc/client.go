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
	"github.com/illusionman1212/gorc/irc/commands"
	"github.com/illusionman1212/gorc/ui"
)

type Message struct {
	Timestamp  string
	Tags       MessageTags // starts with @ | Optional
	Source     string      // starts with : | Optional
	Command    string      // can either be a string or a numeric value | Required
	Parameters []string    // Optional (Dependant on command)
}

type MessageTags map[string]string

type User struct {
	// User prefix in channel
	Prefix string
}

type Channel struct {
	// Channel name
	Name string

	// Channel topic
	Topic string

	// Channel messages history
	History string

	// Users in this channel
	// The map key is the user's nickname
	// and the user struct holds data about that user
	// such as, prefixes in this channel and etc...
	Users map[string]User
}

type Client struct {
	// tcp connection
	TCPConn net.Conn

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

	// Index of the first visible tab in the tab bar
	FirstTabIndexInTabBar int

	// Index of the last visible tab in the tab bar
	LastTabIndexInTabBar int

	// The acknowledged capabilities
	EnabledCapabilities Capabilities
}

type Capabilities map[string]string

type MsgFmtOpts struct {
	WithTimestamp bool
	AsServerMsg   bool
	NotImpl       bool
	AsDate        bool
}

var notImpl = lipgloss.NewStyle().Foreground(ui.ErrorColor).Render("[NOT IMPL]")
var serverMsgStyle = lipgloss.NewStyle().Foreground(ui.ServerMsgColor)
var timestampStyle = serverMsgStyle
var dateStyle = lipgloss.NewStyle().Foreground(ui.DateColor)

const CRLF = "\r\n"

func (m *Message) SetTimestamp() {
	if serverTime, ok := m.Tags["time"]; ok {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", serverTime)
		if err != nil {
			// TODO: properly handle error
			log.Fatalln(err)
		}
		m.Timestamp = fmt.Sprintf("[%02d:%02d]", t.Local().Hour(), t.Local().Minute())
	} else {
		now := time.Now()
		m.Timestamp = fmt.Sprintf("[%02d:%02d]", now.Hour(), now.Minute())
	}
}

func (c *Channel) AppendMsg(timestamp string, fullMsg string, opts MsgFmtOpts) {
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

	c.History += prefixes + style.Render(fullMsg) + CRLF
}

func (c *Client) Initialize(host string, port string, tlsEnabled bool) {
	addr := fmt.Sprintf("%s:%s", host, port)

	if tlsEnabled {
		cfg := &tls.Config{ServerName: host}
		conn, err := tls.Dial("tcp", addr, cfg)
		if err != nil {
			// TODO: properly handle the error instead of Fatal-ing (failed to initiate a connection to server)
			log.Fatal(err)
		}

		c.TCPConn = conn
		c.Host = host
		c.Port = port
		c.ActiveChannel = host
		c.Channels = make([]Channel, 0)
		c.EnabledCapabilities = make(Capabilities, 0)
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

	c.TCPConn = conn
	c.Host = host
	c.Port = port
	c.ActiveChannel = host
	c.Channels = make([]Channel, 0)
	c.EnabledCapabilities = make(Capabilities, 0)
}

func (c *Client) Register(nick string, password string, channel string) {
	c.Channels = append(c.Channels, Channel{
		Name:  c.Host,
		Users: make(map[string]User),
	})

	c.SendCommand(commands.CAP, "LS", "302")
	if password != "" {
		c.SendCommand(commands.PASS, password)
	}
	c.SendCommand(commands.NICK, nick)
	// set user-wanted nickname
	c.Nickname = nick
	c.SendCommand(commands.USER, nick, "0", "*", nick)
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
	c.Channels[0].AppendMsg("", msg, msgOpts)

	// TODO: append new day to each channel whenever the day changes
}

func (c Client) SendCommand(cmd string, params ...string) {
	if c.TCPConn == nil {
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

	c.TCPConn.Write([]byte(cmd + paramsString + CRLF))
}
