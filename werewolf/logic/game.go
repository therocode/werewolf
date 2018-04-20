package logic

import (
	"log"

	"github.com/therocode/werewolf/werewolf/logic/components"
	"github.com/therocode/werewolf/werewolf/timeline"
)

type Game struct {
	timeline []timeline.Event
	players  map[string]string
	base     *Base
	roles    map[string]Role
}

func NewGame(base *Base) *Game {
	instance := &Game{}
	instance.roles = map[string]Role{}
	instance.players = map[string]string{}
	instance.base = base
	return instance
}

func (instance *Game) getGeneratorSet() map[timeline.Generator]bool {
	result := map[timeline.Generator]bool{}
	result[instance.base] = true
	for _, role := range instance.roles {
		result[role] = true
	}
	return result
}

func (instance *Game) AddRole(role Role) {
	instance.roles[role.Name()] = role
}

func (instance *Game) AddPlayer(name string, roleName string) {
	instance.players[name] = roleName
}

func (instance *Game) GetRoles() []Role {
	roles := []Role{}
	for _, roleName := range instance.players {
		role := instance.roles[roleName]
		roles = append(roles, role)
	}
	return roles
}

func (instance *Game) IsPlayer(name string) bool {
	_, contains := instance.players[name]
	return contains
}

func (instance *Game) IsRole(name string, roleName string) bool {
	playerRoleName, _ := instance.players[name]
	return roleName == playerRoleName
}

func (instance *Game) CountComponent(component components.Component) int {
	count := 0
	for _, roleName := range instance.players {
		role := instance.roles[roleName]
		if role.HasComponent(component) {
			count++
		}
	}
	return count
}

func (instance *Game) Kill(player string) {
	delete(instance.players, player)
}

func (instance *Game) Run() {
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
