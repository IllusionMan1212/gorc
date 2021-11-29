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

package main

/* --- IRC Ports --- */
const (
	INSECURE_PORT = 6667
	SECURE_PORT   = 6697
)

/* --- Prohibited characters in channel names --- */
const ( // do we need this as a client ????
	SPACE  = " "  // 0x20
	CTRL_G = "^G" // 0x07
	COMMA  = ","  // 0x2C
)

/* --- COMMANDS --- */
const (
	// Connection messages
	CAP          = "CAP"          // Negotiate capabilities of client and server
	AUTHENTICATE = "AUTHENTICATE" // Used for SASL authentication between client and server. needs to have successfully negotiated the 'sasl' capability
	PASS         = "PASS"         // Set a password to use when registering the client.
	NICK         = "NICK"         // Set a nickname to the client
	USER         = "USER"         // Register the client as a user
	OPER         = "OPER"         // Used to obtain operator privileges, needs a <name> and <password> parameters.
	QUIT         = "QUIT"         // Terminate client connection with an optional <reason>.

	// Channel operations
	JOIN   = "JOIN"   // Join a channel. can take multiple channels
	PART   = "PART"   // Leave a channel. can take multiple channels with an optional reason.
	TOPIC  = "TOPIC"  // Change the channel topic in a mode +t channel. may need operator privs.
	NAMES  = "NAMES"  // View nicknames of all clients in the provided channel. can take multiple channel params separated by comma ","
	LIST   = "LIST"   // Get list of channels with info about each channel.
	INVITE = "INVITE" // Invite a client to an invite-only channel (mode +i)
	KICK   = "KICK"   // Eject a client from the channel

	// Server Queries and Commands
	MOTD    = "MOTD"    // Message of the day command returns either the MOTD or ERR_NOMOTD
	VERSION = "VERSION" // Query version of software and the RPL_ISUPPORT params of the server.
	ADMIN   = "ADMIN"   // Find name of adminstrator of the server.
	CONNECT = "CONNECT" // Force a server to try to establish a new connection to another server. IRC Operators only (do we need this ??)
	TIME    = "TIME"    // Query local time of the server.
	STATS   = "STATS"   // Query statistics of a server.
	INFO    = "INFO"    // Get information about a server. e.g, software name/version, compile date of server, copyright. etc...
	MODE    = "MODE"    // Set or remove modes from a target, either a user(client) target, or a channel target.

	// Sending Messages
	PRIVMSG = "PRIVMSG" // Sends a "private" message to either a channel or another client
	NOTICE  = "NOTICE"  // Send notices between users and to channels.

	// Operator Messages
	KILL = "KILL" // CLose connection between a client and the server. IRC Operators only (do we need this ??)

	// Optional Messages
	AWAY     = "AWAY"     // Indicate that the client (user) is away/afk/etc..
	USERHOST = "USERHOST" // Get information about user with a given nickname. can take up to 5 nickanmes.

	// Left behind...
	PONG = "PONG" // Respond to a server PING
)

/* --- Channel Types --- */
// These are the channel prefixes for different types of channels
const (
	Regular = "#" // 0x23
	Local   = "&" // 0x26
)

/* --- User Modes --- */
const (
	Invisible  = "+i"
	Oper       = "+o"
	Local_Oper = "+O"
	Registered = "+r"
	Wallops    = "+w"
)

/* --- Channel Modes --- */
const (
	Ban                  = "+b"
	Exception            = "+e"
	Client_Limit         = "+l"
	Invite_Only          = "+i"
	Invite_Exception     = "+I"
	Key                  = "+k"
	Moderated            = "+m"
	Secret               = "+s"
	Protected_Topic      = "+t"
	No_External_Messages = "+n"
)

/* --- Channel membership prefixes --- */
const (
	Founder   = "+q"
	Protected = "+a"
	Operator  = "@"
	HalfOp    = "%"
	Voice     = "+v"
)

/* --- Numeric Replies --- */
const (
	RPL_WELCOME         = "001"
	RPL_YOURHOST        = "002"
	RPL_CREATED         = "003"
	RPL_MYINFO          = "004"
	RPL_ISUPPORT        = "005"
	RPL_BOUNCE          = "010"
	RPL_UMODEIS         = "221"
	RPL_LUSERCLIENT     = "251"
	RPL_LUSEROP         = "252"
	RPL_LUSERUNKNOWN    = "253"
	RPL_LUSERCHANNELS   = "254"
	RPL_LUSERME         = "255"
	RPL_ADMINME         = "256"
	RPL_ADMINLOC1       = "257"
	RPL_ADMINLOC2       = "258"
	RPL_ADMINEMAIL      = "259"
	RPL_TRYAGAIN        = "263"
	RPL_LOCALUSERS      = "265"
	RPL_GLOBALUSERS     = "266"
	RPL_WHOISCERTFP     = "276"
	RPL_NONE            = "300"
	RPL_AWAY            = "301"
	RPL_USERHOST        = "302"
	RPL_ISON            = "303"
	RPL_UNAWAY          = "305"
	RPL_NOWAWAY         = "306"
	RPL_WHOISUSER       = "311"
	RPL_WHOISSERVER     = "312"
	RPL_WHOISOPERATOR   = "313"
	RPL_WHOWASUSER      = "314"
	RPL_WHOISIDLE       = "317"
	RPL_ENDOFWHOIS      = "318"
	RPL_WHOISCHANNELS   = "319"
	RPL_LISTSTART       = "321"
	RPL_LIST            = "322"
	RPL_LISTEND         = "323"
	RPL_CHANNELMODEIS   = "324"
	RPL_CREATIONTIME    = "329"
	RPL_NOTOPIC         = "331"
	RPL_TOPIC           = "332"
	RPL_TOPICWHOTIME    = "333"
	RPL_INVITING        = "341"
	RPL_INVITELIST      = "346"
	RPL_ENDOFINVITELIST = "347"
	RPL_EXCEPTLIST      = "348"
	RPL_ENDOFEXCEPTLIST = "349"
	RPL_VERSION         = "351"
	RPL_NAMREPLY        = "353"
	RPL_ENDOFNAMES      = "366"
	RPL_BANLIST         = "367"
	RPL_ENDOFBANLIST    = "368"
	RPL_ENDOFWHOWAS     = "369"
	RPL_MOTDSTART       = "375"
	RPL_MOTD            = "372"
	RPL_ENDOFMOTD       = "376"
	RPL_YOUREOPER       = "381"
	RPL_REHASHING       = "382"

	ERR_UNKNOWNERROR      = "400"
	ERR_NOSUCHNICK        = "401"
	ERR_NOSUCHSERVER      = "402"
	ERR_NOSUCHCHANNEL     = "403"
	ERR_CANNOTSENDTOCHAN  = "404"
	ERR_TOOMANYCHANNELS   = "405"
	ERR_UNKNOWNCOMMAND    = "421"
	ERR_NOMOTD            = "422"
	ERR_ERRONUSNICKNAME   = "432"
	ERR_NICKNAMEINUSE     = "433"
	ERR_USERNOTINCHANNEL  = "441"
	ERR_NOTONCHANNEL      = "442"
	ERR_USERONCHANNEL     = "443"
	ERR_NOTREGISTERED     = "451"
	ERR_NEEDMOREPARAMS    = "461"
	ERR_ALREADYREGISTERED = "462"
	ERR_PASSWDMISMATCH    = "464"
	ERR_YOUREBANNEDCREEP  = "465"
	ERR_CHANNELISFULL     = "471"
	ERR_UNKNOWNMODE       = "472"
	ERR_INVITEONLYCHAN    = "473"
	ERR_BANNEDFROMCHAN    = "474"
	ERR_BADCHANNELKEY     = "475"
	ERR_BADCHANMASK       = "476"
	ERR_NOPRIVILEGES      = "481"
	ERR_CHANOPRIVSNEEDED  = "482"
	ERR_CANTKILLSERVER    = "483"
	ERR_NOOPERHOST        = "491"
	ERR_UMODEUNKNOWNFLAG  = "501"
	ERR_USERSDONTMATCH    = "502"
	RPL_STARTTLS          = "670"
	ERR_STARTTLS          = "691"
	ERR_NOPRIVS           = "723"
	RPL_LOGGEDIN          = "900"
	RPL_LOGGEDOUT         = "901"
	ERR_NICKLOCKED        = "902"
	RPL_SASLSUCCESS       = "903"
	ERR_SASLFAIL          = "904"
	ERR_SASLTOOLONG       = "905"
	ERR_SASLABORTED       = "906"
	ERR_SASLALREADY       = "907"
	RPL_SASLMECHS         = "908"
)

/* --- RPL_ISUPPORT Parameters --- */
// TODO:

// Messages MUST end in CRLF and SHOULD be encoded and decoded using UTF-8 (with fallbacks such as Latin-1).
// Names are casemapped, read the casemapping from the RPL_ISUPPORT the server sends when registration is completed.
// Silently ignore empty messages and only parse messages when you encounter a CRLF.

// Client SHOULD NOT include a source when sending a message, but if included it MUST be the nickname of the client.
// When receiving a numeric reply, client MUST be able to handle any number of parameters on a numeric reply.
