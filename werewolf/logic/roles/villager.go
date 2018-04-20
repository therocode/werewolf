package roles

import (
	"log"

	"github.com/therocode/werewolf/werewolf/logic"
	"github.com/therocode/werewolf/werewolf/logic/components"
	"github.com/therocode/werewolf/werewolf/timeline"
)

type Villager struct {
	communication logic.Communication
	data          logic.Data
	lynchVote     *components.Vote
}

func NewVillager(communication logic.Communication, data logic.Data, vote *components.Vote) *Villager {
	instance := &Villager{}
	instance.communication = communication
	instance.data = data
	instance.lynchVote = vote
	return instance
}

func (*Villager) Name() string {
	return "villager"
}

func (this *Villager) HasComponent(component components.Component) bool {
	return this.lynchVote.Name() == component.Name()
}

func (*Villager) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			"lynch",
			map[string]bool{"day_starts": true},
			map[string]bool{},
		},
	}
}

func (instance *Villager) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	switch event.Name {
	case "night_starts":
		instance.communication.SendToChannel("A villager (%s) sleeps.", player)
	case "day_starts":
		instance.lynchVote.Reset()
	case "lynch":
		vote := instance.getLynchVote(player)
		instance.communication.SendToChannel("%s voted to lynch %s", player, vote)

		instance.lynchVote.Vote(vote)

		voteCount := instance.lynchVote.TotalVoteCount()
		neededVotes := instance.data.CountComponent(instance.lynchVote)
		log.Printf("%d people voted, need %d votes", voteCount, neededVotes)
		if voteCount == neededVotes {
			mostVoted := instance.lynchVote.MostVoted()
			instance.data.Kill(instance.lynchVote.MostVoted())
			instance.communication.SendToChannel("%s was lynched!", mostVoted)
		}
	}

	hasTerminated <- true
}

func (instance *Villager) getLynchVote(player string) string {
	for {
		vote := instance.communication.RequestName(player, "%s, who do you want to lynch?: ", player)

		if vote == player {
			instance.communication.SendToPlayer(player, "You cannot lynch yourself, sorry.")
		} else if !instance.data.IsPlayer(vote) {
			instance.communication.SendToPlayer(player, "That is not a living player!")
		} else {
			return vote
		}
	}
}
