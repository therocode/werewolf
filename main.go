package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/therocode/werewolf/werewolf"
	"github.com/therocode/werewolf/werewolf/irc"
	"github.com/therocode/werewolf/werewolf/ircgame"
	"github.com/therocode/werewolf/werewolf/logic/components"
	"github.com/therocode/werewolf/werewolf/logic/roles"
	"github.com/therocode/werewolf/werewolf/testgame"
	ircevent "github.com/thoj/go-ircevent"
)

const adminChannel = "#wolfadmin"
const serverssl = "irc.boxbox.org:6697"

func runTestGame() {
	log.SetOutput(ioutil.Discard)

	game := testgame.NewTestGame()

	lynchVote := components.NewVote("lynch")
	killVote := components.NewVote("kill")

	game.AddRole(roles.NewVillager(game, game, lynchVote))
	game.AddRole(roles.NewWerewolf(game, game, killVote, lynchVote))

	game.AddPlayer("ulf", "werewolf")
	game.AddPlayer("wulf", "werewolf")
	game.AddPlayer("dolph", "werewolf")
	game.AddPlayer("stig", "villager")
	game.AddPlayer("nils", "villager")
	game.AddPlayer("hans", "villager")
	game.AddPlayer("göte", "villager")
	game.AddPlayer("lennart", "villager")

	game.RunGame()
}

func runIrcGame() {
	rand.Seed(time.Now().UTC().UnixNano())

	ircnick1 := "ulfmann"
	irccon := ircevent.IRC(ircnick1, "Ulf Mannerstrom")

	//var config werewolf.Config //later load from file or something
	//var werewolfInstance *werewolf.Game

	var communication *irc.Irc
	var game *ircgame.IrcGame

	irccon.Debug = false                  //<--- set to true to get lots of IRC debug prints
	irccon.VerboseCallbackHandler = false //<--- set to true to get even more debug prints
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *ircevent.Event) { irccon.Join(adminChannel) })
	irccon.AddCallback("366", func(e *ircevent.Event) {})
	irccon.AddCallback("PRIVMSG", func(e *ircevent.Event) {
		if cmd, err := werewolf.ParseCommand(e.Arguments[0], e.Nick, e.Message()); err == nil {
			if cmd.Command == "newgame" {
				if game == nil {
					communication = irc.NewIrc(irccon, "#"+cmd.Args[0])
					game = newIrcGame(communication)
				} else {
					irccon.Privmsgf(adminChannel, "Cannot start new game with game already in progress")
				}
			} else if game == nil {
				irccon.Privmsgf(adminChannel, "Start a new game with !newgame [channel] first!")
			} else {
				switch cmd.Command {
				case "join":
					game.AddPlayer(e.Nick, cmd.Args[0])
				case "start":
					go game.Run()
				}
			}
		} else {
			if communication != nil {
				communication.Respond(e.Nick, e.Message())
			}
		}
	})
	err := irccon.Connect(serverssl)
	if err != nil {
		log.Printf("Err connecting: %s", err)
		return
	}

	irccon.Loop()
}

func newIrcGame(communication *irc.Irc) *ircgame.IrcGame {
	game := ircgame.NewIrcGame(communication)

	lynchVote := components.NewVote("lynch")
	killVote := components.NewVote("kill")

	game.AddRole(roles.NewVillager(communication, game, lynchVote))
	game.AddRole(roles.NewWerewolf(communication, game, killVote, lynchVote))

	return game
}

func main() {
	runIrcGame()
}
