package testgame

import "github.com/therocode/werewolf/logic"

// TestGame is a local command-line test implementation of the game
type TestGame struct {
	run  bool
	game *logic.Game
}

// NewTestGame creates a new TestGame instance
func NewTestGame(communication logic.Communication) *TestGame {
	instance := &TestGame{}
	instance.run = true
	instance.game = logic.NewGame(instance, communication)
	return instance
}

// RunGame runs the game
func (instance *TestGame) RunGame() {
	for instance.run {
		instance.game.Run()
	}
}

// AddRole adds a role to the game
func (instance *TestGame) AddRole(role logic.Role) {
	instance.game.AddRole(role)
}

// AddPlayer adds a player to the game
func (instance *TestGame) AddPlayer(name string, roleName string) {
	instance.game.AddPlayer(name, roleName)
}

/*
 * Data methods
 */

// EndGame implements the Data interface
func (instance *TestGame) EndGame() {
	instance.run = false
}

// GetRoles implements the Data interface
func (instance *TestGame) GetRoles() []logic.Role {
	return instance.game.GetRoles()
}

// IsPlayer implements the Data interface
func (instance *TestGame) IsPlayer(name string) bool {
	return instance.game.IsPlayer(name)
}

// IsRole implements the Data interface
func (instance *TestGame) IsRole(name string, roleName string) bool {
	return instance.game.IsRole(name, roleName)
}

// CountComponent implements the Data interface
func (instance *TestGame) CountComponent(component logic.Component) int {
	return instance.game.CountComponent(component)
}

// CountRoles implements the Data interface
func (instance *TestGame) CountRoles(roleName string) int {
	return instance.game.CountRoles(roleName)
}

// Kill implements the Data interface
func (instance *TestGame) Kill(player string) {
	instance.game.Kill(player)
}

// GetPlayersWithRole implements the Data interface
func (instance *TestGame) GetPlayersWithRole(roleName string) []string {
	return instance.game.GetPlayersWithRole(roleName)
}
