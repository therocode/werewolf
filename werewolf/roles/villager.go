package roles

import (
	"log"

	"github.com/therocode/werewolf/werewolf/timeline"
)

type Villager struct {
	vote          map[string]int
	communication Communication
	data          Data
}

func NewVillager(communication Communication, data Data) *Villager {
	instance := &Villager{}
	instance.communication = communication
	instance.data = data
	return instance
}

func (*Villager) Name() string {
	return "villager"
}

func (*Villager) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			"lynch",
			map[string]bool{"day_starts": true},
			map[string]bool{"night_starts": true},
		},
	}
}

func (instance *Villager) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	switch event.Name {
	case "night_starts":
		instance.communication.SendToChannel("A villager (%s) sleeps.", player)
	case "day_starts":
		instance.vote = map[string]int{}
	case "lynch":
		vote := instance.communication.RequestName(player, "%s, who do you want to lynch?: ", player)
		instance.communication.SendToChannel("%s wanted to kill %s", player, vote)

		instance.vote[vote]++

		log.Printf("%d people voted, need %d votes", instance.totalVoteCount(), instance.data.CountRole(instance.Name()))
		if instance.totalVoteCount() == instance.data.CountRole(instance.Name()) {
			instance.data.Kill(instance.mostVoted())
		}
	}

	hasTerminated <- true
}

func (instance *Villager) totalVoteCount() int {
	totalVoteCount := 0
	for _, count := range instance.vote {
		totalVoteCount += count
	}
	return totalVoteCount
}

func (instance *Villager) mostVoted() string {
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
