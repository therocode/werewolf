package testgame

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/therocode/werewolf/werewolf/roles"
)

func TestRunGame(t *testing.T) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	game := newTestGame()

	game.addRole(roles.NewWerewolf())
	game.addRole(roles.NewVillager())

	game.addPlayer("ulf", "werewolf")
	game.addPlayer("stig", "villager")
	game.addPlayer("lennart", "villager")

	game.runGame()
}
