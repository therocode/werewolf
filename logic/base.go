package logic

import (
	"time"

	"github.com/therocode/werewolf/logic/roles"
	"github.com/therocode/werewolf/logic/timeline"
)

// Base is a special timeline.Generator and event handler used for creating and managing basic game events.
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

// Generate implements the timeline.Generator interface
func (*Base) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			Name:   NightStarts,
			Before: map[string]bool{},
			After:  map[string]bool{},
		},
		timeline.Event{
			Name:   DayStarts,
			Before: map[string]bool{NightStarts: true},
			After:  map[string]bool{},
		},
	}
}

func (instance *Base) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	if event.Name == NightStarts {
		instance.data.IncrementTurn()
		instance.communication.SendToChannel("Night falls.")
		instance.checkIfGameIsOver()
	} else if event.Name == DayStarts {
		instance.communication.SendToChannel("Day breaks.")
		instance.unmuteChannel()
		instance.checkIfGameIsOver()

		//instance.communication.SendToChannel("5 minutes to go.")
		//time.Sleep(4 * time.Minute)
		instance.communication.SendToChannel("1 minute to go.")
		time.Sleep(30 * time.Second)
		instance.communication.SendToChannel("30 seconds to go.")
		time.Sleep(30 * time.Second)
		instance.muteChannel()
	}
}

func (instance *Base) checkIfGameIsOver() {
	villagerCount := instance.data.CountRoles(roles.Villager, roles.Seer)
	werewolfCount := instance.data.CountRoles(roles.Werewolf)

	if werewolfCount == 0 {
		instance.communication.SendToChannel("All werewolves are dead! Villagers win!")
		instance.data.EndGame()
	} else if villagerCount <= werewolfCount {
		instance.communication.SendToChannel("There are at least as many werewolves as villagers! Werewolves win!")
		instance.data.EndGame()
	}
}

func (instance *Base) muteChannel() {
	for _, player := range instance.data.Players() {
		instance.communication.MutePlayer(player)
	}
}

func (instance *Base) unmuteChannel() {
	for _, player := range instance.data.Players() {
		instance.communication.UnmutePlayer(player)
	}
}
