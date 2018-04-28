package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/therocode/werewolf/ircgame"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/roles"
	"github.com/therocode/werewolf/testgame"
	ircevent "github.com/thoj/go-ircevent"
)

const lobbyChannel = "#wolfadmin"
const serverssl = "irc.boxbox.org:6697"

func runTestGame() {
	log.SetOutput(ioutil.Discard)

	communication := testgame.NewTestCommunication()

	game := testgame.NewTestGame(communication)

	lynchVote := components.NewVote("lynch")
	killVote := components.NewVote("kill")

	game.AddRole(roles.NewVillager(game, communication, lynchVote))
	game.AddRole(roles.NewWerewolf(game, communication, killVote, lynchVote))

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

	lobby := ircgame.NewIrcLobby(ircnick1, lobbyChannel, irccon)

	irccon.Debug = false                  //<--- set to true to get lots of IRC debug prints
	irccon.VerboseCallbackHandler = false //<--- set to true to get even more debug prints
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *ircevent.Event) { irccon.Join(lobbyChannel) })
	irccon.AddCallback("366", func(e *ircevent.Event) {})
	irccon.AddCallback("PRIVMSG", func(e *ircevent.Event) {
		lobby.HandleMessage(e.Arguments[0], e.Nick, e.Message())
	})
	err := irccon.Connect(serverssl)
	if err != nil {
		log.Printf("Err connecting: %s", err)
		return
	}

	irccon.Loop()
}

func main() {
	runIrcGame()
	//runTestGame()
}
