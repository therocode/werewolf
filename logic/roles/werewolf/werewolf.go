package werewolf

import (
	"log"
	"strings"

	"github.com/therocode/werewolf/logic"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/roles"
	"github.com/therocode/werewolf/logic/timeline"
)

// Werewolves can identify other werewolves and can collectively kill one villager per night
type Werewolf struct {
	data          logic.Data
	communication logic.Communication
	killVote      *components.Vote
	lynch         *components.Lynch
}

func NewWerewolf(data logic.Data, communication logic.Communication, killVote *components.Vote, lynch *components.Lynch) *Werewolf {
	instance := &Werewolf{}
	instance.communication = communication
	instance.data = data
	instance.killVote = killVote
	instance.lynch = lynch
	return instance
}

/*
 * Role interface
 */

func (*Werewolf) Name() string {
	return roles.Werewolf
}

func (werewolf *Werewolf) HasComponent(component logic.Component) bool {
	return werewolf.lynch.Name() == component.Name() || werewolf.killVote.Name() == component.Name()
}

func (*Werewolf) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			Name:   logic.WerewolvesSeeEachOther,
			Before: map[string]bool{logic.NightStarts: true},
			After:  map[string]bool{logic.DayStarts: true},
		},
		timeline.Event{
			Name:   logic.WerewolvesKill,
			Before: map[string]bool{logic.WerewolvesSeeEachOther: true},
			After:  map[string]bool{logic.DayStarts: true},
		},
	}
}

func (werewolf *Werewolf) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	werewolf.data.Lock()
	defer func() {
		werewolf.data.Unlock()
		hasTerminated <- true
	}()

	switch event.Name {
	case logic.NightStarts:
		werewolf.killVote.Reset()
	case logic.WerewolvesSeeEachOther:
		werewolves := werewolf.data.GetPlayersWithRole(roles.Werewolf)
		werewolf.communication.SendToPlayer(player, "The werewolves are: %s", strings.Join(werewolves, ", "))
	case logic.WerewolvesKill:
		// Skip werewolf kill on first turn
		if werewolf.data.GetTurnCount() <= 1 {
			return
		}

		vote, timeout := werewolf.getKillVote(player)

		if timeout {
			werewolf.killVote.VoteBlank()
			werewolf.communication.SendToPlayer(player, "You took too long to decide, so your vote will not count!")
		} else {
			log.Printf("%s wanted to kill %s", player, vote)
			werewolf.killVote.Vote(vote)
		}

		totalKillVoteCount := werewolf.killVote.TotalVoteCount()
		neededKillVotes := werewolf.data.CountComponent(werewolf.killVote)
		log.Printf("%d people voted, need %d votes", totalKillVoteCount, neededKillVotes)
		if totalKillVoteCount == neededKillVotes {
			mostVoted, noVotes := werewolf.killVote.MostVoted()

			if noVotes {
				werewolf.communication.SendToChannel("The werewolves were so indecisive, nobody was killed tonight!")
			} else {
				werewolf.data.Kill(mostVoted)
				werewolf.communication.SendToChannel("%s was killed!", mostVoted)
			}
		}
	case logic.DayStarts:
		werewolf.lynch.Reset()
	case logic.Lynch:
		werewolf.lynch.Handle(player)
	}
}

func (werewolf *Werewolf) getKillVote(player string) (string, bool) {
	for {
		werewolf.data.Unlock()
		vote, timeout := werewolf.communication.Request(player, "%s, who do you want to kill?: ", player)
		werewolf.data.Lock()

		switch {
		case timeout:
			return "", true
		case vote == player:
			werewolf.communication.SendToPlayer(player, "You cannot kill yourself, sorry.")
		case !werewolf.data.IsRole(vote, roles.Villager):
			werewolf.communication.SendToPlayer(player, "That is not a living villager!")
		default:
			return vote, false
		}
	}
}
