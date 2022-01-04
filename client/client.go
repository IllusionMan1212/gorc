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

package client

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	// tls connection
	TlsConn *tls.Conn

	// tcp connection
	TcpConn *net.TCPConn
}

const CRLF = "\r\n"

func NewClient(host string, port string, tlsEnabled bool) Client {
	addr := fmt.Sprintf("%s:%s", host, port)

	if tlsEnabled {
		cfg := &tls.Config{ServerName: host}
		conn, err := tls.Dial("tcp", addr, cfg)
		if err != nil {
			// TODO: properly handle the error instead of Fatal-ing (failed to initiate a connection to server)
			log.Fatal(err)
		}

		return Client{
			TlsConn: conn,
			TcpConn: nil,
		}
	}

	addrTCP, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, addrTCP)
	if err != nil {
		log.Fatal(err)
	}

	return Client{
		TlsConn: nil,
		TcpConn: conn,
	}
}

func (c Client) Register(nick string, password string, channel string) {
	c.SendCommand("CAP", "LS")
	if password != "" {
		c.SendCommand("PASS", password)
	}
	// TODO: check if nickname has spaces and remove them
	c.SendCommand("NICK", nick)
	c.SendCommand("USER", nick, "0", "*", nick)
	// TODO: CAP REQ :whatever capability the client recognizes and supports
	c.SendCommand("CAP", "REQ", ":message-tags")
	c.SendCommand("CAP", "END")
	c.SendCommand("JOIN", channel)
}

func (c Client) SendCommand(cmd string, params ...string) {
	if c.TlsConn == nil && c.TcpConn == nil {
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

		// if we have exactly 1 param
	} else if len(params) == 1 {
		// if this 1 param contains spaces, we prepend a colon
		if strings.Contains(params[0], " ") {
			paramsString = " :" + params[0]
		} else {
			paramsString = " " + params[0]
		}
	}

	if c.TlsConn != nil {
		c.TlsConn.Write([]byte(cmd + paramsString + CRLF))
	} else if c.TcpConn != nil {
		c.TcpConn.Write([]byte(cmd + paramsString + CRLF))
	} else {
		// TODO: properly handle the error instead of Fatal-ing
		log.Fatal("No valid connections to write to")
	}
}
