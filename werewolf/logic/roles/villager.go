package roles

import (
	"log"

	"github.com/therocode/werewolf/werewolf/logic"
	"github.com/therocode/werewolf/werewolf/logic/components"
	"github.com/therocode/werewolf/werewolf/logic/timeline"
)

// Villager role
type Villager struct {
	data          logic.Data
	communication logic.Communication
	lynchVote     *components.Vote
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
func (villager *Villager) HasComponent(component components.Component) bool {
	return villager.lynchVote.Name() == component.Name()
}

// Generate implements Role interface
func (*Villager) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			"lynch",
			map[string]bool{"day_starts": true},
			map[string]bool{},
		},
	}
}

// Handle implements Role interface
func (villager *Villager) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	switch event.Name {
	case "night_starts":
	case "day_starts":
		villager.lynchVote.Reset()
	case "lynch":
		vote := villager.getLynchVote(player)
		villager.communication.SendToChannel("%s voted to lynch %s", player, vote)

		villager.lynchVote.Vote(vote)

		voteCount := villager.lynchVote.TotalVoteCount()
		neededVotes := villager.data.CountComponent(villager.lynchVote)
		log.Printf("%d people voted, need %d votes", voteCount, neededVotes)
		if voteCount == neededVotes {
			mostVoted := villager.lynchVote.MostVoted()
			villager.data.Kill(villager.lynchVote.MostVoted())
			villager.communication.SendToChannel("%s was lynched!", mostVoted)
		}
	}

	hasTerminated <- true
}

func (villager *Villager) getLynchVote(player string) string {
	for {
		vote := villager.communication.RequestName(player, "%s, who do you want to lynch?: ", player)

		if vote == player {
			villager.communication.SendToPlayer(player, "You cannot lynch yourself, sorry.")
		} else if !villager.data.IsPlayer(vote) {
			villager.communication.SendToPlayer(player, "That is not a living player!")
		} else {
			return vote
		}
	}
}
