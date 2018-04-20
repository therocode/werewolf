package ircgame

import (
	"github.com/therocode/werewolf/werewolf/irc"
	"github.com/therocode/werewolf/werewolf/logic"
	"github.com/therocode/werewolf/werewolf/logic/components"
)

const (
	GameStateLobby    = "GameStateLobby"
	GameStateStarted  = "GameStateStarted"
	GameStateGameOver = "GameStateGameOver"
)

type IrcGame struct {
	state string
	game  *logic.Game
}

func NewIrcGame(irc *irc.Irc) *IrcGame {
	instance := &IrcGame{}
	instance.state = GameStateLobby
	instance.game = logic.NewGame(logic.NewBase(instance, irc))
	return instance
}

func (instance *IrcGame) Run() {
	instance.state = GameStateStarted
	for instance.state == GameStateStarted {
		instance.game.Run()
	}
}

/*
 * Data methods
 */

func (instance *IrcGame) EndGame() {
	instance.state = GameStateGameOver
}

func (instance *IrcGame) AddRole(role logic.Role) {
	instance.game.AddRole(role)
}

func (instance *IrcGame) AddPlayer(name string, roleName string) {
	instance.game.AddPlayer(name, roleName)
}

func (instance *IrcGame) GetRoles() []logic.Role {
	return instance.game.GetRoles()
}

func (instance *IrcGame) IsPlayer(name string) bool {
	return instance.game.IsPlayer(name)
}

func (instance *IrcGame) IsRole(name string, roleName string) bool {
	return instance.game.IsRole(name, roleName)
}

func (instance *IrcGame) CountComponent(component components.Component) int {
	return instance.game.CountComponent(component)
}

func (instance *IrcGame) Kill(player string) {
	instance.game.Kill(player)
}
