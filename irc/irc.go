package irc

import (
	"fmt"
	"log"
	"sync"
	"time"

	irc "github.com/thoj/go-ircevent"
)

const timeoutDelayInSeconds = 30

type Irc struct {
	irccon           *irc.Connection
	channel          string
	response         map[string]chan string
	responseMapMutex sync.Mutex
}

func NewIrc(irc *irc.Connection, channel string) *Irc {
	this := &Irc{}
	this.irccon = irc
	this.channel = channel
	this.response = map[string]chan string{}
	this.responseMapMutex = sync.Mutex{}
	return this
}

func (irc *Irc) SendToChannel(format string, params ...interface{}) {
	irc.irccon.Privmsg(irc.channel, fmt.Sprintf(format, params...))
}

func (irc *Irc) SendToPlayer(player string, format string, params ...interface{}) {
	irc.irccon.Privmsg(player, fmt.Sprintf(format, params...))
}

func (irc *Irc) MutePlayer(player string) {
	irc.irccon.Mode(irc.channel, fmt.Sprintf("-v %s", player))
}

func (irc *Irc) UnmutePlayer(player string) {
	irc.irccon.Mode(irc.channel, fmt.Sprintf("+v %s", player))
}

func (irc *Irc) MuteChannel() {
	irc.irccon.Mode(irc.channel, "+m")
}

func (irc *Irc) UnmuteChannel() {
	irc.irccon.Mode(irc.channel, "-m")
}

func (irc *Irc) Request(requestFrom string, promptFormat string, params ...interface{}) (string, bool) {
	irc.irccon.Privmsg(requestFrom, fmt.Sprintf(promptFormat, params...))

	channel := make(chan string, 1)
	irc.responseMapMutex.Lock()
	irc.response[requestFrom] = channel
	irc.responseMapMutex.Unlock()

	// Create timeout goroutine
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(timeoutDelayInSeconds * time.Second)
		timeout <- true
	}()

	log.Printf("Requested [%s] from %s", promptFormat, requestFrom)

	// Wait for response, or time out

	select {
	case response := <-channel:
		log.Printf("Response received")
		irc.responseMapMutex.Lock()
		delete(irc.response, requestFrom)
		irc.responseMapMutex.Unlock()
		return response, false
	case <-timeout:
		log.Printf("Timeout occurred")
		irc.responseMapMutex.Lock()
		delete(irc.response, requestFrom)
		irc.responseMapMutex.Unlock()
		return "", true
	}
}

func (irc *Irc) Respond(sender string, message string) {
	irc.responseMapMutex.Lock()
	channel, contains := irc.response[sender]
	irc.responseMapMutex.Unlock()

	if contains {
		channel <- message
	}
}

func (irc *Irc) Leave() {
	irc.irccon.Part(irc.channel)
}
