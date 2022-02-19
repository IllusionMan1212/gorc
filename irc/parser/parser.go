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

package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type IRCMessage struct {
	Timestamp  string
	Tags       Tags     // starts with @ | Optional
	Source     string   // starts with : | Optional
	Command    string   // can either be a string or a numeric value | Required
	Parameters []string // Optional (Dependant on command)
}

type Tags map[string]string

func (m *IRCMessage) setTimestamp() {
	// TODO: get timestamp from time tag
	now := time.Now()
	m.Timestamp = fmt.Sprintf("[%02d:%02d]", now.Hour(), now.Minute())
}

func parseTags(rawTags string) Tags {
	tags := make(Tags, 0)
	rawTagsSlice := strings.Split(rawTags, ";")

	for _, tag := range rawTagsSlice {
		str := strings.Split(tag, "=")

		key := str[0]
		value := ""

		// we need to make sure a "value" exists before accessing str[1]
		if len(str) > 1 {
			value = str[1]
		}

		tags[key] = value
	}

	return tags
}

func ParseIRCMessage(line string) (IRCMessage, bool) {
	ircMessage := IRCMessage{}
	multipleSpacesRegex := regexp.MustCompile("[^\\S\\t]+")
	trailingParamRegex := regexp.MustCompile("[^\\S\\t]+:")
	whitespaceRegex := regexp.MustCompile("^[^\\S\\t]+$")

	if len(line) <= 0 {
		// TODO: do something here ???
		// log.Println("empty message")
		return IRCMessage{}, false
	}

	if line[0] == '@' {
		substrs := multipleSpacesRegex.Split(line, 2)
		tags := parseTags(substrs[0][1:])
		ircMessage.Tags = tags

		line = substrs[1]
	}

	if line[0] == ':' {
		substrs := multipleSpacesRegex.Split(line, 2)
		ircMessage.Source = substrs[0][1:]

		line = substrs[1]
	}

	substrs := multipleSpacesRegex.Split(line, 2)

	// TODO: return error(??) if the irc message is malformed (i.e doesn't contain command)
	ircMessage.Command = substrs[0]

	// if parameters exist and it's not an empty string and it's not whitespace without ":"
	if len(substrs) > 1 && len(substrs[1]) != 0 && !whitespaceRegex.MatchString(substrs[1]) {
		params := trailingParamRegex.Split(substrs[1], 2)

		parameters := multipleSpacesRegex.Split(params[0], -1)
		finalParams := make([]string, 0)

		// if we have 1 param and it starts with a colon
		if len(params) == 1 && params[0][0] == ':' {
			finalParams = append(finalParams, params[0][1:])
		}

		// if there's regular param(s) with trailing param
		// this parses empty params as well
		if len(params) > 1 {
			finalParams = append(parameters, params[1])
		} else {
			// if we only have regular param(s)
			finalParams = append(finalParams, parameters...)
		}

		ircMessage.Parameters = finalParams
	}

	ircMessage.setTimestamp()

	return ircMessage, true
}
