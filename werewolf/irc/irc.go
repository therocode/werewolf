package irc

import (
	"fmt"

	irc "github.com/thoj/go-ircevent"
)

type Irc struct {
	irc      *irc.Connection
	channel  string
	response map[string]chan string
}

func NewIrc(irc *irc.Connection, channel string) *Irc {
	this := &Irc{}
	this.irc = irc
	this.channel = channel
	this.response = map[string]chan string{}
	return this
}

func (this *Irc) SendToChannel(format string, params ...interface{}) {
	this.irc.Privmsg(this.channel, fmt.Sprintf(format, params...))
}

func (this *Irc) SendToPlayer(player string, format string, params ...interface{}) {
	this.irc.Privmsg(player, fmt.Sprintf(format, params...))
}

func (this *Irc) RequestName(requestFrom string, promptFormat string, params ...interface{}) string {
	this.irc.Privmsg(requestFrom, fmt.Sprintf(promptFormat, params...))
	this.response[requestFrom] = make(chan string, 1)
	response := <-this.response[requestFrom]
	delete(this.response, requestFrom)
	return response
}

func (this *Irc) Respond(sender string, message string) {
	if channel, contains := this.response[sender]; contains {
		channel <- message
	}
}
