package logic

// Data interface for the game logic
type Data interface {
	GetRoles() []Role
	EndGame()
	CountComponent(component Component) int
	CountRoles(roleName string) int
	Kill(player string)
	IsPlayer(name string) bool
	IsRole(player string, role string) bool
	GetPlayersWithRole(roleName string) []string
	GetPlayers() []string
	Lock()
	Unlock()
	IncrementTurn()
	GetTurnCount() int
}
