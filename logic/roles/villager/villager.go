package villager

import (
	"github.com/therocode/werewolf/logic"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/roles"
	"github.com/therocode/werewolf/logic/timeline"
)

// Villager is a role with no special abilities
type Villager struct {
	data          logic.Data
	communication logic.Communication
	lynch         *components.Lynch
}

func NewVillager(data logic.Data, communication logic.Communication, lynch *components.Lynch) *Villager {
	instance := &Villager{}
	instance.communication = communication
	instance.data = data
	instance.lynch = lynch
	return instance
}

/*
 * Role interface
 */

func (*Villager) Name() string {
	return roles.Villager
}

func (villager *Villager) HasComponent(component logic.Component) bool {
	return villager.lynch.Name() == component.Name()
}

func (*Villager) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			Name:   logic.Lynch,
			Before: map[string]bool{logic.DayStarts: true},
			After:  map[string]bool{},
		},
	}
}

func (villager *Villager) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	villager.data.Lock()
	defer func() {
		villager.data.Unlock()
		hasTerminated <- true
	}()

	switch event.Name {
	case logic.NightStarts:
		if villager.data.TurnCount() == 1 {
			villager.communication.SendToPlayer(player, "You are a villager!")
		}
	case logic.DayStarts:
		villager.lynch.Reset()
	case logic.Lynch:
		villager.lynch.Handle(player)
	}
}
