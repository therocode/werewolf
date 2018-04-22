package ircgame

import (
	"testing"

	"github.com/therocode/werewolf/werewolf/testgame"
)

func TestIrcGame_assignRoles(t *testing.T) {
	communication := testgame.NewTestCommunication()

	game := NewIrcGame(communication)

	game.AddPlayer("a")
	game.AddPlayer("b")
	game.AddPlayer("c")

	game.assignRoles()
}
