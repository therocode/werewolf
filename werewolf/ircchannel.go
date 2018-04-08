package werewolf

import (
	"fmt"

	irc "github.com/thoj/go-ircevent"
)

type IRCChannel struct {
	irc  *irc.Connection
	name string
}

func newIRCChannel(irc *irc.Connection, name string) IRCChannel {
	irc.Join(name)
	return IRCChannel{irc, name}
}

func (channel *IRCChannel) message(msg string, args ...interface{}) {
	channel.irc.Privmsg(channel.name, fmt.Sprintf(msg, args...))
}
