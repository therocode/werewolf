package roles

import (
	"log"

	"github.com/therocode/werewolf/logic"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/timeline"
)

// Seer role
type Seer struct {
	data          logic.Data
	communication logic.Communication
	lynchVote     *components.Vote
}

// NewSeer creates a new seer instance
func NewSeer(data logic.Data, communication logic.Communication, lynchVote *components.Vote) *Seer {
	instance := &Seer{}
	instance.communication = communication
	instance.data = data
	instance.lynchVote = lynchVote
	return instance
}

// Name implements Role interface
func (*Seer) Name() string {
	return "seer"
}

// HasComponent implements Role interface
func (seer *Seer) HasComponent(component logic.Component) bool {
	return seer.lynchVote.Name() == component.Name()
}

// Generate implements Role interface
func (*Seer) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			Name:   "seer_identifies",
			Before: map[string]bool{"night_starts": true},
			After:  map[string]bool{"day_starts": true},
		},
	}
}

// Handle implements Role interface
func (seer *Seer) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	seer.data.Lock()
	defer func() {
		seer.data.Unlock()
		hasTerminated <- true
	}()

	switch event.Name {
	case "seer_identifies":
		target, timeout := seer.getIdentificationTarget(player)

		if timeout {
			seer.communication.SendToPlayer(player, "Sorry, you took too long to decide")
			return
		}

		if seer.data.GetPlayerRole(target) == "werewolf" {
			seer.communication.SendToPlayer(player, "That is a werewolf!")
		} else {
			seer.communication.SendToPlayer(player, "That is not a werewolf.")
		}
	case "day_starts":
		seer.lynchVote.Reset()
	case "lynch":
		vote, timeout := seer.getLynchVote(player)

		if timeout {
			seer.communication.SendToPlayer(player, "Sorry, you took too long to decide.")
			seer.communication.SendToChannel("%s took too long to decide and forfeited their vote", player)
			seer.lynchVote.VoteBlank()
		} else {
			seer.communication.SendToChannel("%s voted to lynch %s", player, vote)
			seer.lynchVote.Vote(vote)
		}

		totalLynchVoteCount := seer.lynchVote.TotalVoteCount()
		neededLynchVotes := seer.data.CountComponent(seer.lynchVote)
		log.Printf("%d people voted, need %d votes", totalLynchVoteCount, neededLynchVotes)
		if totalLynchVoteCount == neededLynchVotes {
			mostVoted, noVotes := seer.lynchVote.MostVoted()

			if noVotes {
				seer.communication.SendToChannel("The villagers were so indecisive, nobody was lynched today!")
			} else {
				seer.data.Kill(mostVoted)
				seer.communication.SendToChannel("%s was lynched!", mostVoted)
			}
		}
	}
}

func (seer *Seer) getLynchVote(player string) (string, bool) {
	for {
		seer.data.Unlock()
		vote, timeout := seer.communication.Request(player, "%s, who do you want to lynch?: ", player)
		seer.data.Lock()

		switch {
		case timeout:
			return "", true
		case vote == player:
			seer.communication.SendToPlayer(player, "You cannot lynch yourself, sorry.")
		case !seer.data.IsPlayer(vote):
			seer.communication.SendToPlayer(player, "That is not a living player!")
		default:
			return vote, false
		}
	}
}

func (seer *Seer) getIdentificationTarget(player string) (string, bool) {
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
