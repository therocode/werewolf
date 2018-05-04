package logic

import (
	"log"

	"github.com/therocode/werewolf/logic/timeline"
)

// Game contains basic game functionality
type Game struct {
	timeline []timeline.Event
	players  map[string]string
	base     *Base
	roles    map[string]Role
}

func NewGame(data Data, communication Communication) *Game {
	instance := &Game{}
	instance.roles = map[string]Role{}
	instance.players = map[string]string{}
	instance.base = NewBase(data, communication)
	return instance
}

func (instance *Game) generatorSet() map[timeline.Generator]bool {
	result := map[timeline.Generator]bool{}
	result[instance.base] = true
	for _, role := range instance.roles {
		result[role] = true
	}
	return result
}

func (instance *Game) AddRole(role Role) {
	log.Printf("Added role %s", role.Name())
	instance.roles[role.Name()] = role
}

func (instance *Game) AddPlayer(name string, roleName string) {
	log.Printf("Added player %s with role %s", name, roleName)
	instance.players[name] = roleName
}

// Roles return the role of each player (duplicates may occur)
func (instance *Game) Roles() []Role {
	roles := []Role{}
	for _, roleName := range instance.players {
		role := instance.roles[roleName]
		roles = append(roles, role)
	}
	return roles
}

// ContainsRole returns true if the game is configured with the specified role
func (instance *Game) ContainsRole(roleName string) bool {
	_, contains := instance.roles[roleName]
	return contains
}

func (instance *Game) IsPlayer(name string) bool {
	_, contains := instance.players[name]
	return contains
}

func (instance *Game) IsRole(name string, roleName string) bool {
	playerRoleName := instance.players[name]
	return roleName == playerRoleName
}

// CountComponent returns the number of players with a role with the specified component
func (instance *Game) CountComponent(component Component) int {
	count := 0
	for _, roleName := range instance.players {
		role := instance.roles[roleName]
		if role.HasComponent(component) {
			count++
		}
	}
	return count
}

// CountRoles returns the sum of the role count for each supplied role name
func (instance *Game) CountRoles(roleNames ...string) int {
	count := 0
	for _, playerRoleName := range instance.players {
		for _, roleName := range roleNames {
			if roleName == playerRoleName {
				count++
			}
		}
	}
	return count
}

func (instance *Game) Kill(player string) {
	delete(instance.players, player)
}

func (instance *Game) PlayersWithRole(roleName string) []string {
	result := []string{}
	for player, playerRoleName := range instance.players {
		if roleName == playerRoleName {
			result = append(result, player)
		}
	}
	return result
}

// Players returns a list of all living players
func (instance *Game) Players() []string {
	result := []string{}
	for player := range instance.players {
		result = append(result, player)
	}
	return result
}

func (instance *Game) PlayerRole(player string) string {
	return instance.players[player]
}

func (instance *Game) Run() {
	// If there are no events in the timeline, generate more
	if len(instance.timeline) == 0 {
		log.Printf("Timeline is empty, generating events.")
		instance.timeline = timeline.Generate(instance.generatorSet())
		log.Printf("Generated timeline: %s", instance.timeline) //nolint
		if len(instance.timeline) == 0 {
			panic("Couldn't generate a timeline!")
		}
	}

	// Pop the first event in the timeline
	var event timeline.Event
	event, instance.timeline = instance.timeline[0], instance.timeline[1:]
	log.Printf("Popped event: %s", event.Name)

	// Execute the Base handler single-threaded
	instance.base.Handle("", event, nil)

	// Create a goroutine executing the role handler for the event for each player
	hasTerminated := map[string]chan bool{}
	for name, roleName := range instance.players {
		role := instance.roles[roleName]
		hasTerminated[name] = make(chan bool, 1)
		go role.Handle(name, event, hasTerminated[name])
	}

	// Block until all role handlers have finished
	for _, channel := range hasTerminated {
		<-channel
	}
}
