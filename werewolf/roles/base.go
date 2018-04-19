package roles

import (
	"github.com/therocode/werewolf/werewolf/timeline"
)

// Base is a special role used for managing the basic game events.
type Base struct {
	data          Data
	communication Communication
}

func NewBase(data Data, communication Communication) *Base {
	base := &Base{}
	base.data = data
	base.communication = communication
	return base
}

func (*Base) Name() string {
	return "game"
}

func (*Base) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			"night_starts",
			map[string]bool{},
			map[string]bool{},
		},
		timeline.Event{
			"day_starts",
			map[string]bool{"night_starts": true},
			map[string]bool{},
		},
	}
}

func (instance *Base) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	if event.Name == "night_starts" {
		// Prior to nightfall, check if enough villagers or all werewolves are dead
		instance.communication.SendToChannel("Night falls.")
		instance.checkIfGameIsOver()
	} else if event.Name == "day_starts" {
		// Prior to nightfall, check if enough villagers or all werewolves are dead
		instance.communication.SendToChannel("Day breaks.")
		instance.checkIfGameIsOver()
	}
}

func (instance *Base) checkIfGameIsOver() {
	// Are all werewolves dead?
	villagerCount := 0
	werewolfCount := 0
	for _, role := range instance.data.GetRoles() {
		switch role.Name() {
		case "villager":
			villagerCount += 1
		case "werewolf":
			werewolfCount += 1
		}
	}

	if werewolfCount == 0 {
		instance.communication.SendToChannel("All werewolves are dead! Villagers win!")
		instance.data.EndGame()
	} else if villagerCount <= werewolfCount {
		instance.communication.SendToChannel("There are at least as many werewolves as villagers! Werewolves win!")
		instance.data.EndGame()
	}
}
