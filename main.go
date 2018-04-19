package main

import (
	"io/ioutil"
	"log"

	"github.com/therocode/werewolf/werewolf/logic/components"
	"github.com/therocode/werewolf/werewolf/logic/roles"
	"github.com/therocode/werewolf/werewolf/testgame"
)

const channel = "#wolfadmin"
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

func main() {
	runTestGame()
	return
	/*
		rand.Seed(time.Now().UTC().UnixNano())

		ircnick1 := "ulfmann"
		irccon := irc.IRC(ircnick1, "Ulf Mannerstrom")

		var config werewolf.Config //later load from file or something
		var werewolfInstance *werewolf.Game

		irccon.Debug = false                  //<--- set to true to get lots of IRC debug prints
		irccon.VerboseCallbackHandler = false //<--- set to true to get even more debug prints
		irccon.UseTLS = true
		irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
		irccon.AddCallback("366", func(e *irc.Event) {})
		irccon.AddCallback("PRIVMSG", func(e *irc.Event) {

			if cmd, err := werewolf.ParseCommand(e.Arguments[0], e.Nick, e.Message()); err == nil {
				if cmd.Command == "newgame" {
					if werewolfInstance == nil {
						werewolfInstance = werewolf.NewWerewolfGame(irccon, config, "#wolfgame", cmd.Nick) //parse #wolfgame from message or randomize
					} else {
						irccon.Privmsgf(channel, "Cannot start new game with game already in progress")
					}
				} else {
					if werewolfInstance != nil {
						werewolfInstance.HandleCommand(cmd)
					} else {
						irccon.Privmsg(channel, "Start a new game with !newgame first")
					}
				}
			}
		})
		err := irccon.Connect(serverssl)
		if err != nil {
			log.Printf("Err connecting: %s", err)
			return
		}

		irccon.Loop()
	*/
}
