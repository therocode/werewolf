package logic

// Data interface for the game logic
type Data interface {
	Roles() []Role
	EndGame()
	CountComponent(component Component) int
	CountRoles(roleNames ...string) int
	Kill(player string)
	IsPlayer(name string) bool
	IsRole(player string, role string) bool
	PlayersWithRole(roleName string) []string
	Players() []string
	PlayerRole(player string) string
	Lock()
	Unlock()
	IncrementTurn()
	TurnCount() int
}
