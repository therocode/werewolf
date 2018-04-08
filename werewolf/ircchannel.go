package werewolf

import irc "github.com/thoj/go-ircevent"

type IRCChannel struct {
	irc  *irc.Connection
	name string
}

func newIRCChannel(irc *irc.Connection, name string) IRCChannel {
	irc.Join(name)
	return IRCChannel{irc, name}
}
