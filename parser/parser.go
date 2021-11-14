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

package parser

import (
	"strings"
)

type IRCMessage struct {
	Tags       []Tag    // starts with @ | Optional
	Source     string   // starts with : | Optional
	Command    string   // can either be a string or a numeric value | Required
	Parameters []string // Optional (Dependant on command)
}

type Tag struct {
	key   string
	value string
}

func parseTags(rawTags string) []Tag {
	// TODO:
	// TOOD: some tags don't have an "=" and they're basically considered "true"
	return nil
}

func ParseIRCMessage(line string) IRCMessage {
	ircMessage := IRCMessage{}

	if len(line) <= 0 {
		return IRCMessage{}
	}

	if line[0] == '@' {
		substrs := strings.SplitN(line, " ", 2)
		tags := parseTags(substrs[0][1:])
		ircMessage.Tags = tags

		line = substrs[1]
	}

	if line[0] == ':' {
		substrs := strings.SplitN(line, " ", 2)
		ircMessage.Source = substrs[0]

		line = substrs[1]
	}

	substrs := strings.SplitN(line, " ", 2)

	ircMessage.Command = substrs[0]

	// if parameters exist
	if len(substrs) > 1 {
		params := strings.SplitN(substrs[1], " :", 2)

		parameters := strings.Split(params[0], " ")
		// TODO: parameters can be empty. this is denoted by a ":"

		// if there's a trailing parameter
		if len(params) > 1 {
			parameters = append(parameters, params[1])
		}

		ircMessage.Parameters = parameters
	}

	return ircMessage
}
