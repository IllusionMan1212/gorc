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

package commands

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
	PING = "PING"
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
var UserPrefixes = map[string]bool{
	"~": true,
	"&": true,
	"@": true,
	"%": true,
	"+": true,
}

/* --- Numeric Replies --- */
const (
	RPL_WELCOME         = "001" // RFC2812 - Implemented
	RPL_YOURHOST        = "002" // RFC2812 - Implemented
	RPL_CREATED         = "003" // RFC2812 - Implemented
	RPL_MYINFO          = "004" // RFC2812 - Implemented
	RPL_ISUPPORT        = "005" // ??????? - Not Implemented (TODO:)
	RPL_BOUNCE          = "010" // ??????? - Not Implemented (TODO:)
	RPL_REMOTEISUPPORT  = "105" // ??????? - Not Implemented (TODO:) // Exact same as RPL_ISUPPORT but for remote servers
	RPL_TRACELINK       = "200" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACECONNECTING = "201" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACEHANDSHAKE  = "202" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACEUNKNOWN    = "203" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACEOPERATOR   = "204" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACEUSER       = "205" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACESERVER     = "206" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACESERVICE    = "207" // RFC2812 - Not Implemented (TODO:)
	RPL_TRACENEWTYPE    = "208" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACECLASS      = "209" // RFC2812 - Not Implemented (TODO:)
	RPL_TRACERECONNECT  = "210" // RFC2812 - Not Implemented - Deprecated - Has Conflicts (TODO:)
	RPL_STATSLINKINFO   = "211" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSCOMMANDS   = "212" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSCLINE      = "213" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSNLINE      = "214" // RFC1459 - Not Implemented - Has Conflicts (TODO:)
	RPL_STATSILINE      = "215" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSKLINE      = "216" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSQLINE      = "217" // RFC1459 - Not Implemented - Has Conflicts (TODO:)
	RPL_STATSYLINE      = "218" // RFC1459 - Not Implemented (TODO:)
	RPL_ENDOFSTATS      = "219" // RFC1459 - Not Implemented (TODO:)
	RPL_UMODEIS         = "221" // RFC1459 - Not Implemented (TODO:)
	RPL_SERVICEINFO     = "231" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_ENDOFSERVICES   = "232" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_SERVICE         = "233" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_SERVLIST        = "234" // RFC2812 - Not Implemented (TODO:) replaces 231 and 233
	RPL_SERVLISTEND     = "235" // RFC2812 - Not Implemented (TODO:) replaces 232
	RPL_STATSVLINE      = "240" // RFC2812 - Not Implemented - Has Conflicts (TODO:)
	RPL_STATSLLINE      = "241" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSUPTIME     = "242" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSOLINE      = "243" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSHLINE      = "244" // RFC1459 - Not Implemented (TODO:)
	RPL_STATSPING       = "246" // RFC2812 - Not Implemented - Deprecated (TODO:)
	RPL_STATSBLINE      = "247" // RFC2812 - Not Implemented - Has Conflicts (TODO:)
	RPL_STATSDLINE      = "250" // RFC2812 - Not Implemented - Has Conflicts (TODO:)
	RPL_LUSERCLIENT     = "251" // RFC1459 - Implemented
	RPL_LUSEROP         = "252" // RFC1459 - Implemented
	RPL_LUSERUNKNOWN    = "253" // RFC1459 - Implemented
	RPL_LUSERCHANNELS   = "254" // RFC1459 - Implemented
	RPL_LUSERME         = "255" // RFC1459 - Implemented
	RPL_ADMINME         = "256" // RFC1459 - Not Implemented (TODO:)
	RPL_ADMINLOC1       = "257" // RFC1459 - Not Implemented (TODO:)
	RPL_ADMINLOC2       = "258" // RFC1459 - Not Implemented (TODO:)
	RPL_ADMINEMAIL      = "259" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACELOG        = "261" // RFC1459 - Not Implemented (TODO:)
	RPL_TRACEEND        = "262" // RFC2812 - Not Implemented - Has Conflicts (TODO:)
	RPL_TRYAGAIN        = "263" // RFC2812 - Not Implemented (TODO:)
	RPL_LOCALUSERS      = "265" // aircd, Hybrid, Bahamut - Implemented
	RPL_GLOBALUSERS     = "266" // aircd, Hybird, Bahamut - Implemented
	RPL_WHOISCERTFP     = "276" // oftc-hybrid - Not Implemented - Has Conflicts (TODO:)
	RPL_NONE            = "300" // RFC1459 - Not Implemented (TODO:)
	RPL_AWAY            = "301" // RFC1459 - Not Implemented (TODO:)
	RPL_USERHOST        = "302" // RFC1459 - Not Implemented (TODO:)
	RPL_ISON            = "303" // RFC1459 - Not Implemented (TODO:)
	RPL_TEXT            = "304" // irc2 - Not Implemented (TODO:)
	RPL_UNAWAY          = "305" // RFC1459 - Not Implemented (TODO:)
	RPL_NOWAWAY         = "306" // RFC1459 - Not Implemented (TODO:)
	RPL_WHOISUSER       = "311" // RFC1459 - Implemented
	RPL_WHOISSERVER     = "312" // RFC1459 - Implemented
	RPL_WHOISOPERATOR   = "313" // RFC1459 - Not Implemented (TODO:)
	RPL_WHOWASUSER      = "314" // RFC1459 - Not Implemented (TODO:)
	RPL_ENDOFWHO        = "315" // RFC1459 - Not Implemented (TODO:)
	RPL_WHOISIDLE       = "317" // RFC1459 - Implemented
	RPL_ENDOFWHOIS      = "318" // RFC1459 - Implemented
	RPL_WHOISCHANNELS   = "319" // RFC1459 - Implemented
	RPL_LISTSTART       = "321" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_LIST            = "322" // RFC1459 - Not Implemented (TODO:)
	RPL_LISTEND         = "323" // RFC1459 - Not Implemented (TODO:)
	RPL_CHANNELMODEIS   = "324" // RFC1459 - Not Implemented (TODO:)
	RPL_UNIQOPIS        = "325" // RFC2812 - Not Implemented - Has Conflicts (TODO:)
	RPL_CREATIONTIME    = "329" // Bahamut,InspIRCd - Not Implemented (TODO:)
	RPL_NOTOPIC         = "331" // RFC1459 - Not Implemented (TODO:)
	RPL_TOPIC           = "332" // RFC1459 - Not Implemented (TODO:)
	RPL_TOPICWHOTIME    = "333" // ircu,InspIRCd - Not Implemented (TODO:)
	RPL_INVITING        = "341" // RFC1459 - Not Implemented (TODO:)
	RPL_SUMMONING       = "342" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_INVITELIST      = "346" // RFC2812 - Not Implemented (TODO:)
	RPL_ENDOFINVITELIST = "347" // RFC2812 - Not Implemented (TODO:)
	RPL_EXCEPTLIST      = "348" // RFC2812 - Not Implemented (TODO:)
	RPL_ENDOFEXCEPTLIST = "349" // RFC2812 - Not Implemented (TODO:)
	RPL_WHOISGATEWAY    = "350" // InspIRCd - Not Implemented (TODO:)
	RPL_VERSION         = "351" // RFC1459 - Not Implemented (TODO:)
	RPL_WHOREPLY        = "352" // RFC1459 - Not Implemented (TODO:)
	RPL_NAMREPLY        = "353" // RFC1459 - Implemented
	RPL_KILLDONE        = "361" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_CLOSING         = "362" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_CLOSEEND        = "363" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_LINKS           = "364" // RFC1459 - Not Implemented (TODO:)
	RPL_ENDOFLINKS      = "365" // RFC1459 - Not Implemented (TODO:)
	RPL_ENDOFNAMES      = "366" // RFC1459 - Not Implemented (TODO:)
	RPL_BANLIST         = "367" // RFC1459 - Not Implemented (TODO:)
	RPL_ENDOFBANLIST    = "368" // RFC1459 - Not Implemented (TODO:)
	RPL_ENDOFWHOWAS     = "369" // RFC1459 - Not Implemented (TODO:)
	RPL_INFO            = "371" // RFC1459 - Not Implemented (TODO:)
	RPL_MOTD            = "372" // RFC1459 - Implemented
	RPL_INFOSTART       = "373" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_ENDOFINFO       = "374" // RFC1459 - Not Implemented (TODO:)
	RPL_MOTDSTART       = "375" // RFC1459 - Implemented
	RPL_ENDOFMOTD       = "376" // RFC1459 - Not Implemented (TODO:)
	RPL_WHOISHOST       = "378" // ??????? - Implemented - Has Conflicts
	RPL_WHOISMODES      = "379" // ??????? - Implemented - Has Conflicts
	RPL_YOUREOPER       = "381" // RFC1459 - Not Implemented (TODO:)
	RPL_REHASHING       = "382" // RFC1459 - Not Implemented (TODO:)
	RPL_YOURSERVICE     = "383" // RFC2812 - Not Implemented (TODO:)
	RPL_MYPORTIS        = "384" // RFC1459 - Not Implemented - Deprecated (TODO:)
	RPL_TIME            = "391" // RFC1459 - Not Implemented - Has Conflicts (TODO:)
	RPL_USERSSTART      = "392" // RFC1459 - Not Implemented (TODO:)
	RPL_USERS           = "393" // RFC1459 - Not Implemented (TODO:)
	RPL_ENDOFUSERS      = "394" // RFC1459 - Not Implemented (TODO:)
	RPL_NOUSERS         = "395" // RFC1459 - Not Implemented (TODO:)

	ERR_UNKNOWNERROR      = "400" // ??????? - Not Implemented (TODO:)
	ERR_NOSUCHNICK        = "401" // RFC1459 - Not Implemented (TODO:)
	ERR_NOSUCHSERVER      = "402" // RFC1459 - Implemented
	ERR_NOSUCHCHANNEL     = "403" // RFC1459 - Not Implemented (TODO:)
	ERR_CANNOTSENDTOCHAN  = "404" // RFC1459 - Not Implemented (TODO:)
	ERR_TOOMANYCHANNELS   = "405" // RFC1459 - Not Implemented (TODO:)
	ERR_WASNOSUCHNICK     = "406" // RFC1459 - Not Implemented (TODO:)
	ERR_NOSUCHSERVICE     = "408" // RFC2812 - Not Implemented - Has Conflicts (TODO:)
	ERR_NOORIGIN          = "409" // RFC1459 - Not Implemented (TODO:)
	ERR_NORECIPIENT       = "411" // RFC1459 - Not Implemented (TODO:)
	ERR_NOTEXTTOSEND      = "412" // RFC1459 - Not Implemented (TODO:)
	ERR_NOTOPLEVEL        = "413" // RFC1459 - Not Implemented (TODO:)
	ERR_WILDTOPLEVEL      = "414" // RFC1459 - Not Implemented (TODO:)
	ERR_BADMASK           = "415" // RFC2812 - Not Implemented (TODO:)
	ERR_UNKNOWNCOMMAND    = "421" // RFC1459 - Implemented
	ERR_NOMOTD            = "422" // RFC1459 - Not Implemented (TODO:)
	ERR_NOADMININFO       = "423" // RFC1459 - Not Implemented (TODO:)
	ERR_FILEERROR         = "424" // RFC1459 - Not Implemented (TODO:)
	ERR_NONICKNAMEGIVEN   = "431" // RFC1459 - Implemented
	ERR_ERRONEUSNICKNAME  = "432" // RFC1459 - Not Implemented (TODO:)
	ERR_NICKNAMEINUSE     = "433" // RFC1459 - Not Implemented (TODO:)
	ERR_NICKCOLLISION     = "436" // RFC1459 - Not Implemented (TODO:)
	ERR_UNAVAILRESOURCE   = "437" // RFC2812 - Not Implemented - Has Conflicts (TODO:)
	ERR_USERNOTINCHANNEL  = "441" // RFC1459 - Not Implemented (TODO:)
	ERR_NOTONCHANNEL      = "442" // RFC1459 - Not Implemented (TODO:)
	ERR_USERONCHANNEL     = "443" // RFC1459 - Not Implemented (TODO:)
	ERR_NOLOGIN           = "444" // RFC1459 - Not Implemented (TODO:)
	ERR_SUMMONDISABLED    = "445" // RFC1459 - Not Implemented (TODO:)
	ERR_USERSDISABLED     = "446" // RFC1459 - Not Implemented (TODO:)
	ERR_NOTREGISTERED     = "451" // RFC1459 - Not Implemented (TODO:)
	ERR_NEEDMOREPARAMS    = "461" // RFC1459 - Implemented
	ERR_ALREADYREGISTERED = "462" // RFC1459 - Implemented
	ERR_NOPERMFORHOST     = "463" // RFC1459 - Not Implemented (TODO:)
	ERR_PASSWDMISMATCH    = "464" // RFC1459 - Not Implemented (TODO:)
	ERR_YOUREBANNEDCREEP  = "465" // RFC1459 - Not Implemented (TODO:)
	ERR_YOUWILLBEBANNED   = "466" // RFC1459 - Not Implemented - Deprecated (TODO:)
	ERR_KEYSET            = "467" // RFC1459 - Not Implemented (TODO:)
	ERR_CHANNELISFULL     = "471" // RFC1459 - Not Implemented (TODO:)
	ERR_UNKNOWNMODE       = "472" // RFC1459 - Not Implemented (TODO:)
	ERR_INVITEONLYCHAN    = "473" // RFC1459 - Not Implemented (TODO:)
	ERR_BANNEDFROMCHAN    = "474" // RFC1459 - Not Implemented (TODO:)
	ERR_BADCHANNELKEY     = "475" // RFC1459 - Not Implemented (TODO:)
	ERR_BADCHANMASK       = "476" // RFC2812 - Not Implemented (TODO:)
	ERR_BANLISTFULL       = "478" // RFC2812 - Not Implemented (TODO:)
	ERR_NOPRIVILEGES      = "481" // RFC1459 - Not Implemented (TODO:)
	ERR_CHANOPRIVSNEEDED  = "482" // RFC1459 - Not Implemented (TODO:)
	ERR_CANTKILLSERVER    = "483" // RFC1459 - Not Implemented (TODO:)
	ERR_RESTRICTED        = "484" // RFC2812 - Not Implemented - Has Conflicts (TODO:)
	ERR_UNIQOPRIVSNEEDED  = "485" // RFC2812 - Not Implemented (TODO:)
	ERR_NOOPERHOST        = "491" // RFC1459 - Not Implemented (TODO:)
	ERR_NOSERVICEHOST     = "492" // RFC1459 - Not Implemented - Deprecated - Has Conflicts (TODO:)
	ERR_UMODEUNKNOWNFLAG  = "501" // RFC1459 - Not Implemented - Has Conflicts (TODO:)
	ERR_USERSDONTMATCH    = "502" // RFC1459 - Not Implemented (TODO:)
	RPL_STARTTLS          = "670" // IRCv3 - Not Implemented (TODO:)
	ERR_STARTTLS          = "691" // IRCv3 - Not Implemented (TODO:)
	ERR_NOPRIVS           = "723" // RatBox - Not Implemented (TODO:)
	RPL_WHOISKEYVALUE     = "760" // IRCv3 - Not Implemented (TODO:)
	RPL_KEYVALUE          = "761" // IRCv3 - Not Implemented (TODO:)
	RPL_METADATAEND       = "762" // IRCv3 - Not Implemented (TODO:)
	ERR_METADATALIMIT     = "764" // IRCv3 - Not Implemented (TODO:)
	ERR_TARGETINVALID     = "765" // IRCv3 - Not Implemented (TODO:)
	ERR_NOMATCHINGKEY     = "766" // IRCv3 - Not Implemented (TODO:)
	ERR_KEYINVALID        = "767" // IRCv3 - Not Implemented (TODO:)
	ERR_KEYNOTSET         = "768" // IRCv3 - Not Implemented (TODO:)
	ERR_KEYNOPERMISSION   = "769" // IRCv3 - Not Implemented (TODO:)
	RPL_LOGGEDIN          = "900" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
	RPL_LOGGEDOUT         = "901" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
	ERR_NICKLOCKED        = "902" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
	RPL_SASLSUCCESS       = "903" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
	ERR_SASLFAIL          = "904" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
	ERR_SASLTOOLONG       = "905" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
	ERR_SASLABORTED       = "906" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
	ERR_SASLALREADY       = "907" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
	RPL_SASLMECHS         = "908" // Charybdis/Atheme,IRCv3 - Not Implemented (TODO:)
)

/* --- RPL_ISUPPORT Parameters --- */
var Features = map[string]bool{
	"ACCEPT":      true,
	"AWAYLEN":     true,
	"BOT":         true,
	"CALLERID":    true,
	"CASEMAPPING": true,
	"CHANLIMIT":   true,
	"CHANMODES":   true,
	"CHANNELLEN":  true,
	"CHANTYPES":   true, // Channel Types. Default is #. Available are #&
	"CHARSET":     true, // Deprecated but might still be used
	"CLIENTVER":   true, // Deprecated but might still be used
	"CNOTICE":     true,
	"CPRIVMSG":    true,
	"DEAF":        true,
	"ELIST":       true,
	"ESILENCE":    true,
	"ETRACE":      true,
	"EXCEPTS":     true,
	"EXTBAN":      true,
	"FNC":         true, // Deprecated but might still be used
	"INVEX":       true,
	"KEYLEN":      true,
	"KICKLEN":     true,
	"KNOCK":       true,
	"LINELEN":     true, // Proposed
	"MAP":         true, // Deprecated but might still be used
	"MAXBANS":     true, // Deprecated but might still be used
	"MAXCHANNELS": true, // Deprecated but might still be used
	"MAXLIST":     true,
	"MAXNICKLEN":  true,
	"MAXPARA":     true, // Deprecated but might still be used
	"MAXTARGETS":  true,
	"METADATA":    true,
	"MODES":       true,
	"MONITOR":     true,
	"NAMESX":      true, // Deprecated but might still be used
	"NETWORK":     true,
	"NICKLEN":     true,
	"OVERRIDE":    true,
	"PREFIX":      true,
	"SAFELIST":    true,
	"SECURELIST":  true,
	"SILENCE":     true,
	"SSL":         true, // Deprecated but might still be used
	"STARTTLS":    true, // Deprecated but might still be used
	"STATUSMSG":   true,
	"STD":         true, // Deprecated but might still be used
	"TARGMAX":     true,
	"TOPICLEN":    true,
	"UHNAMES":     true, // Deprecated but might still be used
	"USERIP":      true,
	"USERLEN":     true, // Proposed
	"VBANLIST":    true, // Deprecared but might still be used
	"VLIST":       true,
	"WALLCHOPS":   true, // Deprecated but might still be used
	"WALLVOICES":  true, // Deprecated but might still be used
	"WATCH":       true,
	"WHOX":        true,
}

var Capabilities = map[string]bool{
	// Twitch-specific capabilities
	"twitch.tv/membership": false,
	"twitch.tv/commands":   false,
	"twitch.tv/tags":       false,

	// Solanum-specific capabilities
	"solanum.chat/identify-msg": false,
	"solanum.chat/oper":         false,
	"solanum.chat/realhost":     false,

	// IRCv3 capabilities
	"account-notify":       false,
	"account-registration": false, // Draft
	"account-tag":          false,
	"away-notify":          false,
	"batch":                false,
	"cap-notify":           true,
	"channel-rename":       false, // Draft
	"chathistory":          false, // Draft
	"chghost":              false,
	"echo-message":         false,
	"event-playback":       false, // Draft
	"extended-join":        false,
	"extended-monitor":     false,
	"invite-notify":        false,
	"labeled-response":     false,
	"message-tags":         false,
	"metadata":             false,
	"monitor":              false,
	"multi-prefix":         false,
	"multiline":            false, // Draft
	"read-marker":          false, // Draft
	"sasl":                 false,
	"server-time":          true,
	"setname":              false,
	"tls":                  false, // Deprecated
	"userhost-in-names":    false,
}
