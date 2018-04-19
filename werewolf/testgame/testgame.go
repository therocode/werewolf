package testgame

import (
	"fmt"
	"log"

	"github.com/therocode/werewolf/werewolf/logic"
	"github.com/therocode/werewolf/werewolf/logic/components"

	"github.com/therocode/werewolf/werewolf/timeline"
)

type TestGame struct {
	run      bool
	timeline []timeline.Event
	players  map[string]string
	base     *logic.Base
	roles    map[string]logic.Role
}

func NewTestGame() *TestGame {
	instance := &TestGame{}
	instance.run = true
	instance.base = logic.NewBase(instance, instance)
	instance.roles = map[string]logic.Role{}
	instance.players = map[string]string{}
	return instance
}

func (instance *TestGame) AddRole(role logic.Role) {
	instance.roles[role.Name()] = role
}

func (instance *TestGame) AddPlayer(name string, roleName string) {
	instance.players[name] = roleName
}

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

func (instance *TestGame) getGeneratorSet() map[timeline.Generator]bool {
	result := map[timeline.Generator]bool{}
	result[instance.base] = true
	for _, role := range instance.roles {
		result[role] = true
	}
	return result
}

func (instance *TestGame) GetRoles() []logic.Role {
	roles := []logic.Role{}
	for _, roleName := range instance.players {
		role := instance.roles[roleName]
		roles = append(roles, role)
	}
	return roles
}

func (instance *TestGame) EndGame() {
	instance.run = false
}

func (instance *TestGame) RequestName(requestFrom string, promptFormat string, params ...interface{}) string {
	fmt.Printf(promptFormat, params...)
	var text string
	fmt.Scanln(&text)

	return text
}

func (instance *TestGame) IsPlayer(name string) bool {
	_, contains := instance.players[name]
	return contains
}

func (instance *TestGame) IsRole(name string, roleName string) bool {
	playerRoleName, _ := instance.players[name]
	return roleName == playerRoleName
}

func (instance *TestGame) CountRole(roleName string) int {
	roleCount := 0
	for _, playerRoleName := range instance.players {
		if roleName == playerRoleName {
			roleCount++
		}
	}
	return roleCount
}

func (instance *TestGame) CountComponent(component components.Component) int {
	count := 0
	for _, roleName := range instance.players {
		role := instance.roles[roleName]
		if role.HasComponent(component) {
			count++
		}
	}
	return count
}

func (instance *TestGame) Kill(player string) {
	delete(instance.players, player)
	fmt.Printf("%s was killed!\n", player)
}

func (instance *TestGame) RunGame() {
	for instance.run {
		// If there are no events in the timeline, generate more
		if len(instance.timeline) == 0 {
			log.Printf("Timeline is empty, generating events.")
			instance.timeline = timeline.Generate(instance.getGeneratorSet())
			log.Printf("Generated timeline: %s", instance.timeline)
			if len(instance.timeline) == 0 {
				panic("Couldn't generate a timeline!")
			}
		}

		// Pop the first event in the timeline
		var event timeline.Event
		event, instance.timeline = instance.timeline[0], instance.timeline[1:]
		log.Printf("Popped event: %s", event.Name)

		// Create a goroutine executing the game role handler
		instance.base.Handle("", event, nil)

		// Create a goroutine executing the role handler for the event for each player
		hasTerminated := map[string]chan bool{}
		for name, roleName := range instance.players {
			role := instance.roles[roleName]
			hasTerminated[name] = make(chan bool, 1)
			go role.Handle(name, event, hasTerminated[name])
		}

		// Block until all role handlers have finished
		for name, channel := range hasTerminated {
			log.Printf("Waiting for %s's action to terminate...", name)
			<-channel
			log.Printf("%s's action terminated.", name)
		}
	}
}
