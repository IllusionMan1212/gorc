// gorc project
// Copyright (C) 2021 IllusionMan1212 and contributors
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

package client

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	"github.com/illusionman1212/gorc/parser"
)

type Client struct {
	// tcp connection
	Conn *tls.Conn

	// channel to receive data that comes through from the connection
	Receive chan string
}

const CRLF = "\r\n"

func NewClient(host string, port string, tlsEnabled bool) *Client {
	// TODO: distinguish between secure and insecure connections
	addr := fmt.Sprintf("%s:%s", host, port)

	cfg := &tls.Config{ServerName: host}
	conn, err := tls.Dial("tcp", addr, cfg)
	if err != nil {
		// TODO: properly handle the error instead of Fatal-ing (failed to initiate a connection to server)
		log.Fatal(err)
	}

	return &Client{
		Conn:    conn,
		Receive: make(chan string),
	}
}

func (c Client) Register(nick string, password string, channel string) {
	if c.Conn == nil {
		// TODO: properly handle the error instead of Fatal-ing
		log.Fatal("Attempted to write data to nil connection")
	}

	c.sendCommand("CAP", "LS")
	if password != "" {
		c.sendCommand("PASS", password)
	}
	// TODO: check if nickname has spaces and remove them
	c.sendCommand("NICK", nick)
	c.sendCommand("USER", nick, "0", "*", nick)
	// TODO: CAP REQ :whatever capability the client recognizes and supports
	c.sendCommand("CAP", "END")
	c.sendCommand("JOIN", channel)
}

func (c Client) handleCommand(message parser.IRCMessage) {
	fmt.Printf("Tags: %s\n", message.Tags)
	fmt.Printf("Source: %s\n", message.Source)
	fmt.Printf("Command: %s\n", message.Command)
	fmt.Printf("Params: %s\n", message.Parameters)

	switch message.Command {
	case "PING":
		c.sendCommand("PONG")
		break
	}
	// TODO: handle all other commands
}

func (c Client) sendCommand(cmd string, params ...string) {
	paramsString := ""
	if len(params) > 0 {
		paramsString = " " + strings.Join(params, " ")
	}

	// if we have more than 1 param then replace the last param's space with a " :"
	// NOTE: not sure if we should only check if we have more than 1 param
	// cuz one param could have spaces and would require a colon
	if len(params) > 1 {
		i := strings.LastIndex(paramsString, " ")
		paramsString = paramsString[:i] + strings.Replace(paramsString[i:], " ", " :", 1)
	}

	c.Conn.Write([]byte(cmd + paramsString + CRLF))
}

func (c Client) Run() {
	for {
		select {
		case line := <-c.Receive:
			message := parser.ParseIRCMessage(line)
			c.handleCommand(message)
		}
	}
}
