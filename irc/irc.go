package irc

import (
	"fmt"
	"log"
	"sync"
	"time"

	irc "github.com/thoj/go-ircevent"
)

const timeoutDelayInSeconds = 30

// Irc is an IRC implementation of the Communication interface
type Irc struct {
	irccon   *irc.Connection
	channel  string
	response map[string]chan string
	mutex    sync.Mutex
}

// NewIrc creates a new Irc instance
func NewIrc(irc *irc.Connection, channel string) *Irc {
	this := &Irc{}
	this.irccon = irc
	this.channel = channel
	this.response = map[string]chan string{}
	this.mutex = sync.Mutex{}
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
func (irc *Irc) Request(requestFrom string, promptFormat string, params ...interface{}) (string, bool) {
	irc.irccon.Privmsg(requestFrom, fmt.Sprintf(promptFormat, params...))

	channel := make(chan string, 1)
	irc.mutex.Lock()
	irc.response[requestFrom] = channel
	irc.mutex.Unlock()

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
		delete(irc.response, requestFrom)
		return response, false
	case <-timeout:
		log.Printf("Timeout occurred")
		delete(irc.response, requestFrom)
		return "", true
	}
}

// Respond is used to supply player input as responses to game queries initiated by Request
func (irc *Irc) Respond(sender string, message string) {
	irc.mutex.Lock()
	channel, contains := irc.response[sender]
	irc.mutex.Unlock()

	if contains {
		channel <- message
	}
}
