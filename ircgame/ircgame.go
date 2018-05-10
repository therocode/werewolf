package ircgame

import (
	"log"
	"math/rand"
	"sync"

	"github.com/therocode/werewolf/logic"
	"github.com/therocode/werewolf/logic/roles"
)

const (
	gameStateLobby    = "GameStateLobby"
	gameStateStarted  = "GameStateStarted"
	gameStateGameOver = "GameStateGameOver"
)

// IrcGame is an IRC-based implementation of the game
type IrcGame struct {
	state         string
	communication logic.Communication
	game          *logic.Game
	players       []string
	dataMutex     sync.Mutex
	turnCount     int
}

func NewIrcGame(communication logic.Communication) *IrcGame {
	instance := &IrcGame{}
	instance.state = gameStateLobby
	instance.communication = communication
	instance.game = logic.NewGame(instance, communication)
	instance.players = []string{}
	return instance
}

func (instance *IrcGame) Run() {
	// Recover from a general panic by ending the game and printing the error message
	defer func() {
		if r := recover(); r != nil {
			instance.communication.SendToChannel("Game was terminated, please start a new one: %s", r)
		}

		instance.EndGame()
	}()

	instance.assignRoles()
	instance.communication.MuteChannel()

	instance.state = gameStateStarted
	for instance.state == gameStateStarted {
		instance.game.Run()
	}
}

func (instance *IrcGame) AddRole(role logic.Role) {
	instance.game.AddRole(role)
}

func (instance *IrcGame) AddPlayer(name string) {
	log.Printf("%s joined the game", name)
	instance.communication.SendToChannel("%s joined the game", name)
	instance.players = append(instance.players, name)
}

// assignRoles randomly assigns roles to the current set of players
func (instance *IrcGame) assignRoles() {
	playerCount := len(instance.players)

	// Verify that there are at least 4 players
	if playerCount < 4 {
		panic("Cannot assign roles with less than 4 players!")
	}

	remainingRoles := []string{}

	// If there are 4-5 players, there should be only one werewolf. If there's 6 or more, there should be two.
	if playerCount == 4 || playerCount == 5 {
		remainingRoles = append(remainingRoles, roles.Werewolf)
	} else if playerCount >= 6 {
		remainingRoles = append(remainingRoles, roles.Werewolf, roles.Werewolf)
	}

	// There should be one seer
	if instance.game.ContainsRole(roles.Seer) {
		remainingRoles = append(remainingRoles, roles.Seer)
	}

	// The remaining players are villagers
	remainingRoleCount := len(remainingRoles)
	for i := 0; i < playerCount-remainingRoleCount; i++ {
		remainingRoles = append(remainingRoles, roles.Villager)
	}

	for _, player := range instance.players {
		// Remove a random role from the remaining roles list
		i := rand.Intn(len(remainingRoles))
		role := remainingRoles[i]
		remainingRoles = append(remainingRoles[:i], remainingRoles[i+1:]...)

		// Add the player with that role
		instance.game.AddPlayer(player, role)
	}
}

func (instance *IrcGame) IsFinished() bool {
	return instance.state == gameStateGameOver
}

func (instance *IrcGame) IsRunning() bool {
	return instance.state == gameStateStarted
}

/*
 * Data methods
 */

func (instance *IrcGame) EndGame() {
	instance.state = gameStateGameOver
	instance.communication.UnmuteChannel()
	instance.communication.Leave()
}

func (instance *IrcGame) Roles() []logic.Role {
	return instance.game.Roles()
}

func (instance *IrcGame) IsPlayer(name string) bool {
	return instance.game.IsPlayer(name)
}

func (instance *IrcGame) IsRole(name string, roleName string) bool {
	return instance.game.IsRole(name, roleName)
}

func (instance *IrcGame) CountComponent(component logic.Component) int {
	return instance.game.CountComponent(component)
}

func (instance *IrcGame) CountRoles(roleNames ...string) int {
	return instance.game.CountRoles(roleNames...)
}

func (instance *IrcGame) Kill(player string) {
	instance.game.Kill(player)
}

func (instance *IrcGame) PlayersWithRole(roleName string) []string {
	return instance.game.PlayersWithRole(roleName)
}

func (instance *IrcGame) Players() []string {
	return instance.game.Players()
}

func (instance *IrcGame) PlayerRole(player string) string {
	return instance.game.PlayerRole(player)
}

func (instance *IrcGame) Lock() {
	instance.dataMutex.Lock()
}

func (instance *IrcGame) Unlock() {
	instance.dataMutex.Unlock()
}

func (instance *IrcGame) IncrementTurn() {
	instance.turnCount++
}

func (instance *IrcGame) TurnCount() int {
	return instance.turnCount
}
