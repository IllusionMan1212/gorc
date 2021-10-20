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
	"fmt"
	"strings"
)

func parseIRCMessage(line string) IRCMessage {
	ircMessage := &IRCMessage{}

	switch line[0] {
	case '@':
		substrs := strings.SplitN(line, " ", 4)

		tagsSlice := make(map[string]string)

		tags := strings.Split(substrs[0][1:], ";")

		for _, tag := range tags {
			splitTag := strings.Split(tag, "=")

			// TOOD: some tags don't have an "=" and they're basically considered "true"
			tagsSlice[splitTag[0]] = splitTag[1]
		}

		ircMessage := IRCMessage{
			Tags:    tagsSlice,
			Source:  substrs[1][1:],
			Command: substrs[2],
		}

		ircMessage.Tags = tagsSlice
		ircMessage.Source = substrs[1][1:]
		ircMessage.Command = substrs[2]

		fmt.Printf("%s\n", ircMessage.Tags)
		fmt.Printf("%s\n", ircMessage.Source)
		fmt.Printf("%s\n", ircMessage.Command)
		fmt.Printf("%s\n", ircMessage.Parameters)
		// TODO:

		break
	case ':':
		substrs := strings.SplitN(line, " ", 3)

		parameters := make([]string, 0)

		// if parameters exist
		if len(substrs) > 2 {
			params := strings.SplitN(substrs[2], " :", 2)

			parameters = strings.Split(params[0], " ")
			// TODO: parameters can be empty. this is denoted by a ":"

			// if there's a final space-included param
			if len(params) > 1 {
				parameters = append(parameters, params[1])
			}
		}

		ircMessage.Source = substrs[0][1:]
		ircMessage.Command = substrs[1]
		ircMessage.Parameters = parameters

		break
	default:
		substrs := strings.SplitN(line, " ", 2)

		parameters := make([]string, 0)

		// if parameters exist
		if len(substrs) > 1 {
			params := strings.SplitN(substrs[1], " :", 2)

			parameters = strings.Split(params[0], " ")
			// TODO: parameters can be empty. this is denoted by a ":"

			// if there's a trailing parameter
			if len(params) > 1 {
				parameters = append(parameters, params[1])
			}
		}

		ircMessage.Source = ""
		ircMessage.Command = substrs[0]
		ircMessage.Parameters = parameters

		break
	}

	return *ircMessage
}
