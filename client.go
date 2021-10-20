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

package main

import (
	"crypto/tls"
	"fmt"
	"log"
)

type Client struct {
	// tcp connection
	conn *tls.Conn

	// channel to receive data that comes through from the connection
	receive chan string
}

func NewClient(host string, port string) *Client {
	// TODO: allow for insecure connections thru a checkbox (needed for irc servers that maybe don't have encryption on)
	// TODO: distinguish between secure and insecure connections
	cfg := &tls.Config{ServerName: host}
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := tls.Dial("tcp", addr, cfg)
	if err != nil {
		// TODO: properly handle the error instead of Fatal-ing (failed to initiate a connection to server)
		log.Fatal(err)
	}

	// TODO: read nickname and channel from user input (stdin ?? or TUI)
	conn.Write([]byte("CAP LS\r\n"))
	// conn.Write([]byte("PASS password"))
	conn.Write([]byte("NICK illusion434\r\n"))
	conn.Write([]byte("USER illusion434 0 * :illusion434\r\n"))
	// TODO: CAP REQ :whatever capability the client recognizes and supports
	conn.Write([]byte("CAP END\r\n"))
	conn.Write([]byte("JOIN ##programming\r\n"))

	return &Client{
		conn:    conn,
		receive: make(chan string),
	}
}

func (c Client) handleCommand(message IRCMessage) {
	fmt.Printf("Source: %s\n", message.Source)
	fmt.Printf("Command: %s\n", message.Command)
	fmt.Printf("Params: %s\n", message.Parameters)

	switch message.Command {
	case "PING":
		c.conn.Write([]byte("PONG\r\n"))
		break
	}
}

func (c Client) Run() {
	for {
		select {
		case line := <-c.receive:
			messages := parseIRCMessage(line)
			c.handleCommand(messages)
		}
	}
}
