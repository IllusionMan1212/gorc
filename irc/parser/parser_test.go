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

import "testing"

// tests provided by https://github.com/ircdocs/parser-tests/blob/master/tests/msg-split.yaml

func TestParser(t *testing.T) {
	t.Run("Test simple command", func(t *testing.T) {
		testMessage := "foo bar baz asdf"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing message with only command and params")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing message with only command and params (1st param)")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing message with only command and params (2nd param)")
		}

		if ircMessage.Parameters[2] != "asdf" {
			t.Fatal("Parsing message with only command and params (3rd param)")
		}
	})

	t.Run("Test source", func(t *testing.T) {
		testMessage := ":coolguy foo bar baz asdf"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "coolguy" {
			t.Fatal("Parsing regular source")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing command with source")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param of regular message")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param of regular message")
		}

		if ircMessage.Parameters[2] != "asdf" {
			t.Fatal("Parsing 3rd param of regular message")
		}
	})

	t.Run("Test trailing param", func(t *testing.T) {
		testMessage := "foo bar baz :asdf quux"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param in message with trailing param")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param in message with trailing param")
		}

		if ircMessage.Parameters[2] != "asdf quux" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test empty trailing param", func(t *testing.T) {
		testMessage := "foo bar baz :"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param in message with trailing param")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param in message with trailing param")
		}

		if ircMessage.Parameters[2] != "" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test multiple colons at trailing param", func(t *testing.T) {
		testMessage := "foo bar baz ::asdf"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param in message with trailing param")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param in message with trailing param")
		}

		if ircMessage.Parameters[2] != ":asdf" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test source and trailing param", func(t *testing.T) {
		testMessage := ":coolguy foo bar baz :asdf quux"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "coolguy" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "asdf quux" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test source and trailing param with spaces", func(t *testing.T) {
		testMessage := ":coolguy foo bar baz :  asdf quux "
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "coolguy" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "  asdf quux " {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test source and trailing param 2", func(t *testing.T) {
		testMessage := ":coolguy PRIVMSG bar :lol :) "
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "coolguy" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "PRIVMSG" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "lol :) " {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test source and empty trailing param", func(t *testing.T) {
		testMessage := ":coolguy foo bar baz :"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "coolguy" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test source and trailing param with spaces 2", func(t *testing.T) {
		testMessage := ":coolguy foo bar baz :  "
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "coolguy" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "  " {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test tags", func(t *testing.T) {
		testMessage := "@id=234AB;rose :dan!d@localhost PRIVMSG #chan :Hey what's up!"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["id"] != "234AB" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["rose"] != "" {
			t.Fatal("Parsing value-less tag")
		}

		if ircMessage.Source != "dan!d@localhost" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "PRIVMSG" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "#chan" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "Hey what's up!" {
			t.Fatal("Parsing 2nd param")
		}
	})

	t.Run("Test tags 2", func(t *testing.T) {
		testMessage := "@a=b;c=32;k;rt=q17 foo"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["a"] != "b" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["c"] != "32" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["k"] != "" {
			t.Fatal("Parsing value-less tag")
		}

		if ircMessage.Tags["rt"] != "q17" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing command")
		}
	})

	// TODO: fails
	t.Run("Test escaped tags", func(t *testing.T) {
		testMessage := "@a=b\\\\and\\nk;c=72\\s45;d=gh\\:764 foo"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["a"] != "b\\and\nk" {
			t.Fatal("Parsing escaped tag")
		}

		if ircMessage.Tags["c"] != "72 45" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["d"] != "gh;764" {
			t.Fatal("Parsing escaped tag")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing command")
		}
	})

	t.Run("Test tags with source", func(t *testing.T) {
		testMessage := "@c;h=;a=b :quux ab cd"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["c"] != "" {
			t.Log(ircMessage.Tags)
			t.Fatal("Parsing value-less tag")
		}

		if ircMessage.Tags["h"] != "" {
			t.Fatal("Parsing value-less tag")
		}

		if ircMessage.Tags["a"] != "b" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Source != "quux" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "ab" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "cd" {
			t.Fatal("Parsing param")
		}
	})

	// TODO: FAIL
	t.Run("Test single param with preceeding :", func(t *testing.T) {
		testMessage := ":src JOIN :#chan"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Parameters[0] != "#chan" {
			t.Fatal("Parsing a single param with a preceeding \":\"")
		}
	})

	t.Run("Test with no space as last param", func(t *testing.T) {
		testMessage := ":src AWAY"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "src" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "AWAY" {
			t.Fatal("Parsing command")
		}

		if len(ircMessage.Parameters) > 0 {
			t.Fatal("Parsing params")
		}
	})

	t.Run("Test with space as last param", func(t *testing.T) {
		testMessage := ":src AWAY "
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "src" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "AWAY" {
			t.Fatal("Parsing command")
		}

		if len(ircMessage.Parameters) > 0 {
			t.Logf("'%v'", ircMessage.Parameters[0])
			t.Fatal("Parsing params")
		}
	})

	t.Run("Test tab as space", func(t *testing.T) {
		testMessage := ":cool\tguy foo bar baz"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "cool\tguy" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "foo" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "bar" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "baz" {
			t.Fatal("Parsing 2nd param")
		}
	})

	t.Run("Test weird control codes in source", func(t *testing.T) {
		testMessage := ":coolguy!ag@net\x035w\x03ork.admin PRIVMSG foo :bar baz"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "coolguy!ag@net\x035w\x03ork.admin" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "PRIVMSG" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "foo" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "bar baz" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test more weird control codes", func(t *testing.T) {
		testMessage := ":coolguy!~ag@n\x02et\x0305w\x0fork.admin PRIVMSG foo :bar baz"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "coolguy!~ag@n\x02et\x0305w\x0fork.admin" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "PRIVMSG" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "foo" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "bar baz" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test full message", func(t *testing.T) {
		testMessage := "@tag1=value1;tag2;vendor1/tag3=value2;vendor2/tag4= :irc.example.com COMMAND param1 param2 :param3 param3"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["tag1"] != "value1" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["tag2"] != "" {
			t.Fatal("Parsing value-less tag")
		}

		if ircMessage.Tags["vendor1/tag3"] != "value2" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["vendor2/tag4"] != "" {
			t.Fatal("Parsing value-less tag")
		}

		if ircMessage.Source != "irc.example.com" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "param1" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "param2" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "param3 param3" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test message without tags", func(t *testing.T) {
		testMessage := ":irc.example.com COMMAND param1 param2 :param3 param3"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "irc.example.com" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "param1" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "param2" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "param3 param3" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test message without source", func(t *testing.T) {
		testMessage := "@tag1=value1;tag2;vendor1/tag3=value2;vendor2/tag4= COMMAND param1 param2 :param3 param3"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["tag1"] != "value1" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["tag2"] != "" {
			t.Fatal("Parsing value-less tag")
		}

		if ircMessage.Tags["vendor1/tag3"] != "value2" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["vendor2/tag4"] != "" {
			t.Fatal("Parsing value-less tag")
		}

		if ircMessage.Source != "" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "param1" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "param2" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "param3 param3" {
			t.Fatal("Parsing trailing param")
		}
	})

	// TODO: FAIL
	t.Run("Test yaml encoding with slashes", func(t *testing.T) {
		testMessage := "@foo=\\\\\\\\\\:\\\\s\\s\\r\\n COMMAND"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["foo"] != "\\\\;\\s \r\n" {
			t.Fatal("Parsing escaped tag")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}
	})

	t.Run("Test broken messages 1", func(t *testing.T) {
		testMessage := ":gravel.mozilla.org 432  #momo :Erroneous Nickname: Illegal characters"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "gravel.mozilla.org" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "432" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "#momo" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "Erroneous Nickname: Illegal characters" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test broken messages 2", func(t *testing.T) {
		testMessage := ":gravel.mozilla.org MODE #tckk +n "
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "gravel.mozilla.org" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "MODE" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "#tckk" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "+n" {
			t.Fatal("Parsing 2nd param")
		}
	})

	t.Run("Test broken messages 3", func(t *testing.T) {
		testMessage := ":services.esper.net MODE #foo-bar +o foobar  "
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "services.esper.net" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "MODE" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "#foo-bar" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "+o" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "foobar" {
			t.Fatal("Parsing trailing param")
		}
	})

	// TODO: FAIL
	t.Run("Test tags. they should be parsed char-at-a-time", func(t *testing.T) {
		testMessage := "@tag1=value\\\\ntest COMMAND"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["tag1"] != "value\\ntest" {
			t.Fatal("Parsing escaped tag")
		}

		if ircMessage.Source != "" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}
	})

	// TODO: FAIl
	t.Run("Test tag escape for char that doesn't need it", func(t *testing.T) {
		testMessage := "@tag1=value\\1 COMMAND"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["tag1"] != "value1" {
			t.Fatal("Parsing escaped tag")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}
	})

	// TODO: FAIL
	t.Run("Test slash at end of tag", func(t *testing.T) {
		testMessage := "@tag1=value1\\ COMMAND"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["tag1"] != "value1" {
			t.Fatal("Parsing escaped tag")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}
	})

	t.Run("Test duplicate tags", func(t *testing.T) {
		testMessage := "@tag1=1;tag2=3;tag3=4;tag1=5 COMMAND"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["tag1"] != "5" {
			t.Fatal("Parsing duplicate tag")
		}

		if ircMessage.Tags["tag2"] != "3" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["tag3"] != "4" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}
	})

	t.Run("Test vendored tags can have the same name as unvendored tag", func(t *testing.T) {
		testMessage := "@tag1=1;tag2=3;tag3=4;tag1=5;vendor/tag2=8 COMMAND"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Tags["tag1"] != "5" {
			t.Fatal("Parsing duplicate tag")
		}

		if ircMessage.Tags["tag2"] != "3" {
			t.Fatal("Parsing unvendored tag")
		}

		if ircMessage.Tags["tag3"] != "4" {
			t.Fatal("Parsing regular tag")
		}

		if ircMessage.Tags["vendor/tag2"] != "8" {
			t.Fatal("Parsing vendored tag")
		}

		if ircMessage.Command != "COMMAND" {
			t.Fatal("Parsing command")
		}
	})

	t.Run("Test: some parsers handle /MODE in a special way", func(t *testing.T) {
		testMessage := ":SomeOp MODE #channel :+i"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "SomeOp" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "MODE" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "#channel" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "+i" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test: some parsers handle /MODE in a special way 2", func(t *testing.T) {
		testMessage := ":SomeOp MODE #channel +oo SomeUser :AnotherUser"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "SomeOp" {
			t.Fatal("Parsing source")
		}

		if ircMessage.Command != "MODE" {
			t.Fatal("Parsing command")
		}

		if ircMessage.Parameters[0] != "#channel" {
			t.Fatal("Parsing 1st param")
		}

		if ircMessage.Parameters[1] != "+oo" {
			t.Fatal("Parsing 2nd param")
		}

		if ircMessage.Parameters[2] != "SomeUser" {
			t.Fatal("Parsing 3rd param")
		}

		if ircMessage.Parameters[3] != "AnotherUser" {
			t.Fatal("Parsing trailing param")
		}
	})

	t.Run("Test multiple spaces", func(t *testing.T) {
		testMessage := ":src    JOIN    #chan     :thing lol haha"
		ircMessage, value := ParseIRCMessage(testMessage)

		if !value {
			t.Fatal("Invalid irc message")
		}

		if ircMessage.Source != "src" {
			t.Fatal("Parsing source with multiple spaces in between")
		}

		if ircMessage.Command != "JOIN" {
			t.Fatal("Parsing command with multiple spaces in between")
		}

		if ircMessage.Parameters[0] != "#chan" {
			t.Fatal("Parsing params with multiple spaces in between")
		}

		if ircMessage.Parameters[1] != "thing lol haha" {
			t.Fatal("Parsing trailing param with multiple spaces in between")
		}
	})
}
