package seer

import (
	"github.com/therocode/werewolf/logic"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/roles"
	"github.com/therocode/werewolf/logic/timeline"
)

// Seer is a role that can determine whether another player is a werewolf once per day
type Seer struct {
	data          logic.Data
	communication logic.Communication
	lynch         *components.Lynch
}

func NewSeer(data logic.Data, communication logic.Communication, lynch *components.Lynch) *Seer {
	instance := &Seer{}
	instance.communication = communication
	instance.data = data
	instance.lynch = lynch
	return instance
}

/*
 * Role interface
 */

func (*Seer) Name() string {
	return roles.Seer
}

func (seer *Seer) HasComponent(component logic.Component) bool {
	return seer.lynch.Name() == component.Name()
}

func (*Seer) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			Name:   logic.SeerIdentifies,
			Before: map[string]bool{logic.NightStarts: true},
			After:  map[string]bool{logic.DayStarts: true},
		},
	}
}

func (seer *Seer) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	seer.data.Lock()
	defer func() {
		seer.data.Unlock()
		hasTerminated <- true
	}()

	switch event.Name {
	case logic.SeerIdentifies:
		target, timeout := seer.requestIdentificationTarget(player)

		if timeout {
			seer.communication.SendToPlayer(player, "Sorry, you took too long to decide")
			return
		}

		if seer.data.PlayerRole(target) == roles.Werewolf {
			seer.communication.SendToPlayer(player, "That is a werewolf!")
		} else {
			seer.communication.SendToPlayer(player, "That is not a werewolf.")
		}
	case logic.DayStarts:
		seer.lynch.Reset()
	case logic.Lynch:
		seer.lynch.Handle(player)
	}
}

func (seer *Seer) requestIdentificationTarget(player string) (string, bool) {
	for {
		seer.data.Unlock()
		vote, timeout := seer.communication.Request(player, "%s, who do you want to identify?: ", player)
		seer.data.Lock()

		switch {
		case timeout:
			return "", true
		case vote == player:
			seer.communication.SendToPlayer(player, "You cannot identify yourself, sorry.")
		case !seer.data.IsPlayer(vote):
			seer.communication.SendToPlayer(player, "That is not a living player!")
		default:
			return vote, false
		}
	}
}
