package roles

import (
	"log"

	"github.com/therocode/werewolf/werewolf/timeline"
)

type Werewolf struct {
	vote          map[string]int
	villager      *Villager
	communication Communication
	data          Data
}

func NewWerewolf(communication Communication, data Data, villager *Villager) *Werewolf {
	instance := &Werewolf{}
	instance.communication = communication
	instance.data = data
	instance.villager = villager
	return instance
}

func (role *Werewolf) Name() string {
	return "werewolf"
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
	if event.Name == "night_starts" {
		instance.vote = map[string]int{}
		instance.communication.SendToChannel("A werewolf (%s) prowls.", player)
	} else if event.Name == "werewolves_kill" {
		vote := instance.getVote(player)
		instance.communication.SendToChannel("%s wanted to kill %s", player, vote)

		instance.vote[vote]++

		log.Printf("%d people voted, need %d votes", instance.totalVoteCount(), instance.data.CountRole(instance.Name()))
		if instance.totalVoteCount() == instance.data.CountRole(instance.Name()) {
			instance.data.Kill(instance.mostVoted())
		}
	}
	hasTerminated <- true
}

func (instance *Werewolf) getVote(player string) string {
	validVote := false
	for !validVote {
		vote := instance.communication.RequestName(player, "%s, who do you want to kill?: ", player)

		validVote = instance.data.IsRole(vote, "villager")
		if validVote {
			return vote
		}

		instance.communication.SendToChannel("That is not a living villager!")
	}

	return ""
}

func (instance *Werewolf) totalVoteCount() int {
	totalVoteCount := 0
	for _, count := range instance.vote {
		totalVoteCount += count
	}
	return totalVoteCount
}

func (instance *Werewolf) mostVoted() string {
	maxVoteCount := 0
	var mostVoted string
	for name, count := range instance.vote {
		if count > maxVoteCount {
			maxVoteCount = count
			mostVoted = name
		}
	}

	return mostVoted
}
