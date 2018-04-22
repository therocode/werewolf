package ircgame

import (
	"log"
	"math/rand"

	"github.com/therocode/werewolf/werewolf/logic"
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
}

// NewIrcGame creates a new instance of IrcGame
func NewIrcGame(communication logic.Communication) *IrcGame {
	instance := &IrcGame{}
	instance.state = gameStateLobby
	instance.communication = communication
	instance.game = logic.NewGame(instance, communication)
	instance.players = []string{}
	return instance
}

// Run the game
func (instance *IrcGame) Run() {
	// Recover from a general panic by ending the game and printing the error message
	defer func() {
		if r := recover(); r != nil {
			instance.EndGame()
			instance.communication.SendToChannel("Game was terminated, please start a new one: %s", r)
		}
	}()

	instance.assignRoles()

	instance.state = gameStateStarted
	for instance.state == gameStateStarted {
		instance.game.Run()
	}

	log.Printf("Game ended.")
}

// AddRole adds a role to the game
func (instance *IrcGame) AddRole(role logic.Role) {
	instance.game.AddRole(role)
}

// AddPlayer adds a player to the game
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
		remainingRoles = append(remainingRoles, "werewolf")
	} else if playerCount >= 6 {
		remainingRoles = append(remainingRoles, "werewolf", "werewolf")
	}

	// The remaining players are villagers
	remainingRoleCount := len(remainingRoles)
	for i := 0; i < playerCount-remainingRoleCount; i++ {
		remainingRoles = append(remainingRoles, "villager")
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

// IsFinished returns true if the game is over
func (instance *IrcGame) IsFinished() bool {
	return instance.state == gameStateGameOver
}

/*
 * Data methods
 */

// EndGame implements the Data interface
func (instance *IrcGame) EndGame() {
	instance.state = gameStateGameOver
}

// GetRoles implements the Data interface
func (instance *IrcGame) GetRoles() []logic.Role {
	return instance.game.GetRoles()
}

// IsPlayer implements the Data interface
func (instance *IrcGame) IsPlayer(name string) bool {
	return instance.game.IsPlayer(name)
}

// IsRole implements the Data interface
func (instance *IrcGame) IsRole(name string, roleName string) bool {
	return instance.game.IsRole(name, roleName)
}

// CountComponent implements the Data interface
func (instance *IrcGame) CountComponent(component logic.Component) int {
	return instance.game.CountComponent(component)
}

// CountRoles implements the Data interface
func (instance *IrcGame) CountRoles(roleName string) int {
	return instance.game.CountRoles(roleName)
}

// Kill implements the Data interface
func (instance *IrcGame) Kill(player string) {
	instance.game.Kill(player)
}

// GetPlayersWithRole implements the Data interface
func (instance *IrcGame) GetPlayersWithRole(roleName string) []string {
	return instance.game.GetPlayersWithRole(roleName)
}
