package testgame

import (
	"fmt"

	"github.com/therocode/werewolf/werewolf/logic"
	"github.com/therocode/werewolf/werewolf/logic/components"
)

type TestGame struct {
	run  bool
	game *logic.Game
}

func NewTestGame() *TestGame {
	instance := &TestGame{}
	instance.run = true
	instance.game = logic.NewGame(logic.NewBase(instance, instance))
	return instance
}

/*
 * Communication methods
 */

func (instance *TestGame) SendToChannel(format string, params ...interface{}) {
	if len(params) == 0 {
		fmt.Print(format + "\n")
	} else {
		fmt.Printf(format+"\n", params...)
	}
}

func (instance *TestGame) SendToPlayer(player string, format string, params ...interface{}) {
	fmt.Printf("PM for %s: ", player)
	if len(params) == 0 {
		fmt.Print(format + "\n")
	} else {
		fmt.Printf(format+"\n", params...)
	}
}

func (instance *TestGame) RequestName(requestFrom string, promptFormat string, params ...interface{}) string {
	fmt.Printf(promptFormat+"\n", params...)
	var text string
	fmt.Scanln(&text)

	return text
}

/*
 * Data methods
 */

func (instance *TestGame) EndGame() {
	instance.run = false
}

func (instance *TestGame) AddRole(role logic.Role) {
	instance.game.AddRole(role)
}

func (instance *TestGame) AddPlayer(name string, roleName string) {
	instance.game.AddPlayer(name, roleName)
}

func (instance *TestGame) GetRoles() []logic.Role {
	return instance.game.GetRoles()
}

func (instance *TestGame) IsPlayer(name string) bool {
	return instance.game.IsPlayer(name)
}

func (instance *TestGame) IsRole(name string, roleName string) bool {
	return instance.game.IsRole(name, roleName)
}

func (instance *TestGame) CountComponent(component components.Component) int {
	return instance.game.CountComponent(component)
}

func (instance *TestGame) Kill(player string) {
	instance.game.Kill(player)
}

func (instance *TestGame) RunGame() {
	for instance.run {
		instance.game.Run()
	}
}
