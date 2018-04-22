package irc

import (
	"fmt"

	irc "github.com/thoj/go-ircevent"
)

// Irc is an IRC implementation of the Communication interface
type Irc struct {
	irccon   *irc.Connection
	channel  string
	response map[string]chan string
}

// NewIrc creates a new Irc instance
func NewIrc(irc *irc.Connection, channel string) *Irc {
	this := &Irc{}
	this.irccon = irc
	this.channel = channel
	this.response = map[string]chan string{}
	return this
}

// SendToChannel implements the Communication interface
func (irc *Irc) SendToChannel(format string, params ...interface{}) {
	irc.irccon.Privmsg(irc.channel, fmt.Sprintf(format, params...))
}

// SendToPlayer implements the Communication interface
func (irc *Irc) SendToPlayer(player string, format string, params ...interface{}) {
	irc.irccon.Privmsg(player, fmt.Sprintf(format, params...))
}

// Request implements the Communication interface
func (irc *Irc) Request(requestFrom string, promptFormat string, params ...interface{}) string {
	irc.irccon.Privmsg(requestFrom, fmt.Sprintf(promptFormat, params...))
	irc.response[requestFrom] = make(chan string, 1)
	response := <-irc.response[requestFrom]
	delete(irc.response, requestFrom)
	return response
}

// Respond is used to supply player input as responses to game queries initiated by Request
func (irc *Irc) Respond(sender string, message string) {
	if channel, contains := irc.response[sender]; contains {
		channel <- message
	}
}
