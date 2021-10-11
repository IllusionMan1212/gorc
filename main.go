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
	"io"
	"log"
)

func main() {
	// TOOD: read the server (and port?) from user input (stdin ?? or TUI)
	// TODO: allow for insecure connections thru a checkbox (needed for irc servers that maybe don't have encryption on)
	cfg := &tls.Config{ServerName: "irc.libera.chat"}
	conn, err := tls.Dial("tcp", "irc.libera.chat:6697", cfg)
	if err != nil {
		// TODO: properly handle the error instead of Fatal-ing (failed to initiate a connection to server)
		log.Fatal(err)
	}
	defer conn.Close()

	// TODO: read nickname and channel from user input (stdin ?? or TUI)
	conn.Write([]byte("CAP LS\r\n"))
	// TODO: PASS if it requires a password
	conn.Write([]byte("NICK illusion\r\n"))
	conn.Write([]byte("USER illusion 0 * :illusion\r\n"))
	// TODO: CAP REQ :whatever capability the client recognizes and supports
	conn.Write([]byte("CAP END\r\n"))
	// conn.Write([]byte("JOIN #libera\r\n"))

	for {
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Print("End of file received\n")
				return
			}
		}

		// TODO: parse received data according to the protocol

		// TOOD: reply to PINGs with PONGs to keep the connection alive

		fmt.Printf("Data: %s\n", buf)
	}
}
