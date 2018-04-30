package roles

import (
	"log"
	"sync"

	"github.com/therocode/werewolf/logic"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/timeline"
)

// Villager role
type Villager struct {
	data          logic.Data
	communication logic.Communication
	lynchVote     *components.Vote
	dataMutex     sync.Mutex
}

// NewVillager creates a new villager instance
func NewVillager(data logic.Data, communication logic.Communication, vote *components.Vote) *Villager {
	instance := &Villager{}
	instance.communication = communication
	instance.data = data
	instance.lynchVote = vote
	return instance
}

// Name implements Role interface
func (*Villager) Name() string {
	return "villager"
}

// HasComponent implements Role interface
func (villager *Villager) HasComponent(component logic.Component) bool {
	return villager.lynchVote.Name() == component.Name()
}

// Generate implements Role interface
func (*Villager) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			Name:   "lynch",
			Before: map[string]bool{"day_starts": true},
			After:  map[string]bool{},
		},
	}
}

// Handle implements Role interface
func (villager *Villager) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	villager.data.Lock()
	defer func() {
		villager.data.Unlock()
		hasTerminated <- true
	}()

	switch event.Name {
	case "night_starts":
	case "day_starts":
		villager.lynchVote.Reset()
	case "lynch":
		vote, timeout := villager.getLynchVote(player)

		if timeout {
			villager.communication.SendToPlayer(player, "Sorry, you took too long to decide.")
			villager.communication.SendToChannel("%s took too long to decide, and forfeited their vote", player)
			villager.lynchVote.VoteBlank()
		} else {
			villager.communication.SendToChannel("%s voted to lynch %s", player, vote)
			villager.lynchVote.Vote(vote)
		}

		voteCount := villager.lynchVote.TotalVoteCount()
		neededVotes := villager.data.CountComponent(villager.lynchVote)
		log.Printf("%d people voted, need %d votes", voteCount, neededVotes)
		if voteCount == neededVotes {
			mostVoted, noVotes := villager.lynchVote.MostVoted()

			if noVotes {
				villager.communication.SendToChannel("The villagers were so indecisive, nobody was lynched today!")
			} else {
				villager.data.Kill(mostVoted)
				villager.communication.SendToChannel("%s was lynched!", mostVoted)
			}
		}
	}
}

func (villager *Villager) getLynchVote(player string) (string, bool) {
	for {
		villager.data.Unlock()
		vote, timeout := villager.communication.Request(player, "%s, who do you want to lynch?: ", player)
		villager.data.Lock()

		switch {
		case timeout:
			return "", true
		case vote == player:
			villager.communication.SendToPlayer(player, "You cannot lynch yourself, sorry.")
		case !villager.data.IsPlayer(vote):
			villager.communication.SendToPlayer(player, "That is not a living player!")
		default:
			return vote, false
		}
	}
}
