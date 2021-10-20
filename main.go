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
	"bufio"
	"log"
)

func main() {
	// create new client with the provided host and port
	client := NewClient("irc.libera.chat", "6697")
	defer client.conn.Close()
	go client.Run()

	// main loop
	r := bufio.NewReaderSize(client.conn, 512)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Print(err)
		}

		client.receive <- msg
	}
}
