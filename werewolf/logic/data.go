package logic

import (
	"github.com/therocode/werewolf/werewolf/logic/components"
)

type Data interface {
	GetRoles() []Role
	EndGame()
	CountComponent(component components.Component) int
	CountRoles(roleName string) int
	Kill(player string)
	IsPlayer(name string) bool
	IsRole(player string, role string) bool
	GetPlayersWithRole(roleName string) []string
}
