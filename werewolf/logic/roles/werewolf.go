package roles

import (
	"log"

	"github.com/therocode/werewolf/werewolf/logic"
	"github.com/therocode/werewolf/werewolf/logic/components"

	"github.com/therocode/werewolf/werewolf/timeline"
)

type Werewolf struct {
	communication logic.Communication
	data          logic.Data
	killVote      *components.Vote
	lynchVote     *components.Vote
}

func NewWerewolf(communication logic.Communication, data logic.Data, killVote *components.Vote, lynchVote *components.Vote) *Werewolf {
	instance := &Werewolf{}
	instance.communication = communication
	instance.data = data
	instance.killVote = killVote
	instance.lynchVote = lynchVote
	return instance
}

func (role *Werewolf) Name() string {
	return "werewolf"
}

func (this *Werewolf) HasComponent(component components.Component) bool {
	return this.lynchVote.Name() == component.Name() || this.killVote.Name() == component.Name()
}

func (*Werewolf) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			"werewolves_kill",
			map[string]bool{"night_starts": true},
			map[string]bool{"day_starts": true},
		},
	}
}

func (instance *Werewolf) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	switch event.Name {
	case "night_starts":
		instance.killVote.Reset()
		instance.communication.SendToChannel("A werewolf (%s) prowls.", player)
	case "werewolves_kill":
		vote := instance.getKillVote(player)
		instance.communication.SendToChannel("%s wanted to kill %s", player, vote)

		instance.killVote.Vote(vote)

		totalKillVoteCount := instance.killVote.TotalVoteCount()
		neededKillVotes := instance.data.CountComponent(instance.killVote)
		log.Printf("%d people voted, need %d votes", totalKillVoteCount, neededKillVotes)
		if totalKillVoteCount == neededKillVotes {
			mostVoted := instance.killVote.MostVoted()
			instance.data.Kill(instance.killVote.MostVoted())
			instance.communication.SendToChannel("%s was killed!", mostVoted)
		}
	case "day_starts":
		instance.lynchVote.Reset()
	case "lynch":
		vote := instance.getLynchVote(player)
		instance.communication.SendToChannel("%s voted to lynch %s", player, vote)

		instance.lynchVote.Vote(vote)

		totalLynchVoteCount := instance.lynchVote.TotalVoteCount()
		neededLynchVotes := instance.data.CountComponent(instance.lynchVote)
		log.Printf("%d people voted, need %d votes", totalLynchVoteCount, neededLynchVotes)
		if totalLynchVoteCount == neededLynchVotes {
			mostVoted := instance.lynchVote.MostVoted()
			instance.data.Kill(instance.lynchVote.MostVoted())
			instance.communication.SendToChannel("%s was lynched!", mostVoted)
		}
	}
	hasTerminated <- true
}

func (instance *Werewolf) getKillVote(player string) string {
	for {
		vote := instance.communication.RequestName(player, "%s, who do you want to kill?: ", player)

		if vote == player {
			instance.communication.SendToPlayer(player, "You cannot kill yourself, sorry.")
		} else if !instance.data.IsRole(vote, "villager") {
			instance.communication.SendToPlayer(player, "That is not a living villager!")
		} else {
			return vote
		}
	}
}

func (instance *Werewolf) getLynchVote(player string) string {
	for {
		vote := instance.communication.RequestName(player, "%s, who do you want to lynch?: ", player)

		if vote == player {
			instance.communication.SendToPlayer(player, "You cannot lynch yourself, sorry.")
		} else if !instance.data.IsPlayer(vote) {
			instance.communication.SendToPlayer(player, "That is not a living player!")
		} else if instance.data.IsRole(vote, "werewolf") {
			instance.communication.SendToPlayer(player, "You cannot lynch another werewolf!")
		} else {
			return vote
		}
	}
}
