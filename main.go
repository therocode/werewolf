package main

import (
	"crypto/tls"
	"log"
	"math/rand"
	"time"

	"github.com/therocode/werewolf/werewolf"
	"github.com/thoj/go-ircevent"
)

const channel = "#wolfadmin"
const serverssl = "irc.boxbox.org:6697"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	ircnick1 := "ulfmann"
	irccon := irc.IRC(ircnick1, "Ulf Mannerstrom")

	var config werewolf.Config //later load from file or something
	var werewolfInstance *werewolf.Werewolf = werewolf.NewWerewolf(irccon, config)

	irccon.Debug = false                  //<--- set to true to get lots of IRC debug prints
	irccon.VerboseCallbackHandler = false //<--- set to true to get even more debug prints
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		message := e.Message()
		if message[0] == '!' { //only process messages starting with '!'
			nick := e.Nick
			channel := e.Arguments[0] //arg 0 for privmsg is the channel name
			werewolfInstance.HandleMessage(channel, nick, message)
		}
	})
	err := irccon.Connect(serverssl)
	if err != nil {
		log.Printf("Err connecting: %s", err)
		return
	}

	irccon.Loop()
}
