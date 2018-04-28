package logic

import (
	"time"

	"github.com/therocode/werewolf/logic/timeline"
)

// Base is a special timeline.Generator and event handler used for creating and managing basic game events.
type Base struct {
	data          Data
	communication Communication
}

// NewBase creates new Base instance
func NewBase(data Data, communication Communication) *Base {
	base := &Base{}
	base.data = data
	base.communication = communication
	return base
}

// Generate implements the timeline.Generator interface
func (*Base) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			Name:   "night_starts",
			Before: map[string]bool{},
			After:  map[string]bool{},
		},
		timeline.Event{
			Name:   "day_starts",
			Before: map[string]bool{"night_starts": true},
			After:  map[string]bool{},
		},
	}
}

// Handle implements the Role interface
func (instance *Base) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	if event.Name == "night_starts" {
		// Prior to nightfall, check if enough villagers or all werewolves are dead
		instance.communication.SendToChannel("Night falls.")
		instance.checkIfGameIsOver()
	} else if event.Name == "day_starts" {
		// Prior to nightfall, check if enough villagers or all werewolves are dead
		instance.communication.SendToChannel("Day breaks.")
		instance.checkIfGameIsOver()

		instance.communication.SendToChannel("5 minutes to go.")
		time.Sleep(4 * time.Minute)
		instance.communication.SendToChannel("1 minute to go.")
		time.Sleep(30 * time.Second)
		instance.communication.SendToChannel("30 seconds to go.")
		time.Sleep(30 * time.Second)
	}
}

func (instance *Base) checkIfGameIsOver() {
	// Are all werewolves dead?
	villagerCount := instance.data.CountRoles("villager")
	werewolfCount := instance.data.CountRoles("werewolf")

	if werewolfCount == 0 {
		instance.communication.SendToChannel("All werewolves are dead! Villagers win!")
		instance.data.EndGame()
	} else if villagerCount <= werewolfCount {
		instance.communication.SendToChannel("There are at least as many werewolves as villagers! Werewolves win!")
		instance.data.EndGame()
	}
}
