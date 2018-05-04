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
	"github.com/therocode/werewolf/logic/roles/villager"
	"github.com/therocode/werewolf/logic/roles/werewolf"
	"github.com/therocode/werewolf/testgame"
	ircevent "github.com/thoj/go-ircevent"
)

const lobbyChannel = "#wolfadmin"
const serverssl = "irc.boxbox.org:6697"

//nolint
func runTestGame() {
	log.SetOutput(ioutil.Discard)

	communication := testgame.NewTestCommunication()

	game := testgame.NewTestGame(communication)

	lynch := components.NewLynch(game, communication)
	killVote := components.NewVote("kill")

	game.AddRole(villager.NewVillager(game, communication, lynch))
	game.AddRole(werewolf.NewWerewolf(game, communication, killVote, lynch))

	game.AddPlayer("ulf", roles.Werewolf)
	game.AddPlayer("wulf", roles.Werewolf)
	game.AddPlayer("dolph", roles.Werewolf)
	game.AddPlayer("stig", roles.Villager)
	game.AddPlayer("nils", roles.Villager)
	game.AddPlayer("hans", roles.Villager)
	game.AddPlayer("göte", roles.Villager)
	game.AddPlayer("lennart", roles.Villager)

	game.RunGame()
}

func runIrcGame() {
	rand.Seed(time.Now().UTC().UnixNano())

	ircnick := "ulfmann"
	irccon := ircevent.IRC(ircnick, "Ulf Mannerstrom")

	lobby := ircgame.NewIrcLobby(ircnick, lobbyChannel, irccon)

	irccon.Debug = false                  //<--- set to true to get lots of IRC debug prints
	irccon.VerboseCallbackHandler = false //<--- set to true to get even more debug prints
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true} //nolint
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
