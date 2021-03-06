package testgame

import (
	"github.com/therocode/werewolf/logic"
)

// TestGame is a local command-line test implementation of the game
type TestGame struct {
	run       bool
	game      *logic.Game
	turnCount int
}

func NewTestGame(communication logic.Communication) *TestGame {
	instance := &TestGame{}
	instance.run = true
	instance.game = logic.NewGame(instance, communication)
	return instance
}

func (instance *TestGame) RunGame() {
	for instance.run {
		instance.game.Run()
	}
}

func (instance *TestGame) AddRole(role logic.Role) {
	instance.game.AddRole(role)
}

func (instance *TestGame) AddPlayer(name string, roleName string) {
	instance.game.AddPlayer(name, roleName)
}

/*
 * Data interface
 */

func (instance *TestGame) EndGame() {
	instance.run = false
}

func (instance *TestGame) Roles() []logic.Role {
	return instance.game.Roles()
}

func (instance *TestGame) IsPlayer(name string) bool {
	return instance.game.IsPlayer(name)
}

func (instance *TestGame) IsRole(name string, roleName string) bool {
	return instance.game.IsRole(name, roleName)
}

func (instance *TestGame) CountComponent(component logic.Component) int {
	return instance.game.CountComponent(component)
}

func (instance *TestGame) CountRoles(roleNames ...string) int {
	return instance.game.CountRoles(roleNames...)
}

func (instance *TestGame) Kill(player string) {
	instance.game.Kill(player)
}

func (instance *TestGame) PlayersWithRole(roleName string) []string {
	return instance.game.PlayersWithRole(roleName)
}

func (instance *TestGame) Players() []string {
	return instance.game.Players()
}

func (instance *TestGame) PlayerRole(player string) string {
	return instance.game.PlayerRole(player)
}

func (instance *TestGame) Lock() {
}

func (instance *TestGame) Unlock() {
}

func (instance *TestGame) IncrementTurn() {
	instance.turnCount++
}

func (instance *TestGame) TurnCount() int {
	return instance.turnCount
}
