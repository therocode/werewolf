package roles

type Data interface {
	GetRoles() []Role
	EndGame()
	CountRole(role string) int
	Kill(player string)
	IsPlayer(name string) bool
	IsRole(player string, role string) bool
}
