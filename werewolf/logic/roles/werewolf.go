package roles

import (
	"log"
	"strings"

	"github.com/therocode/werewolf/werewolf/logic"
	"github.com/therocode/werewolf/werewolf/logic/components"

	"github.com/therocode/werewolf/werewolf/logic/timeline"
)

// Werewolf role
type Werewolf struct {
	data          logic.Data
	communication logic.Communication
	killVote      *components.Vote
	lynchVote     *components.Vote
}

// NewWerewolf creates a new werewolf instance
func NewWerewolf(data logic.Data, communication logic.Communication, killVote *components.Vote, lynchVote *components.Vote) *Werewolf {
	instance := &Werewolf{}
	instance.communication = communication
	instance.data = data
	instance.killVote = killVote
	instance.lynchVote = lynchVote
	return instance
}

// Name implements Role interface
func (*Werewolf) Name() string {
	return "werewolf"
}

// HasComponent implements Role interface
func (werewolf *Werewolf) HasComponent(component components.Component) bool {
	return werewolf.lynchVote.Name() == component.Name() || werewolf.killVote.Name() == component.Name()
}

// Generate implements Role interface
func (*Werewolf) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			"werewolves_see_each_other",
			map[string]bool{"night_starts": true},
			map[string]bool{},
		},
		timeline.Event{
			"werewolves_kill",
			map[string]bool{"werewolves_see_each_other": true},
			map[string]bool{"day_starts": true},
		},
	}
}

// Handle implements Role interface
func (werewolf *Werewolf) Handle(player string, event timeline.Event, hasTerminated chan bool) {
	switch event.Name {
	case "night_starts":
		werewolf.killVote.Reset()
	case "werewolves_see_each_other":
		werewolves := werewolf.data.GetPlayersWithRole(werewolf.Name())
		werewolf.communication.SendToPlayer(player, "The werewolves are: %s", strings.Join(werewolves, ", "))
	case "werewolves_kill":
		vote := werewolf.getKillVote(player)
		log.Printf("%s wanted to kill %s", player, vote)

		werewolf.killVote.Vote(vote)

		totalKillVoteCount := werewolf.killVote.TotalVoteCount()
		neededKillVotes := werewolf.data.CountComponent(werewolf.killVote)
		log.Printf("%d people voted, need %d votes", totalKillVoteCount, neededKillVotes)
		if totalKillVoteCount == neededKillVotes {
			mostVoted := werewolf.killVote.MostVoted()
			werewolf.data.Kill(werewolf.killVote.MostVoted())
			werewolf.communication.SendToChannel("%s was killed!", mostVoted)
		}
	case "day_starts":
		werewolf.lynchVote.Reset()
	case "lynch":
		vote := werewolf.getLynchVote(player)
		werewolf.communication.SendToChannel("%s voted to lynch %s", player, vote)

		werewolf.lynchVote.Vote(vote)

		totalLynchVoteCount := werewolf.lynchVote.TotalVoteCount()
		neededLynchVotes := werewolf.data.CountComponent(werewolf.lynchVote)
		log.Printf("%d people voted, need %d votes", totalLynchVoteCount, neededLynchVotes)
		if totalLynchVoteCount == neededLynchVotes {
			mostVoted := werewolf.lynchVote.MostVoted()
			werewolf.data.Kill(werewolf.lynchVote.MostVoted())
			werewolf.communication.SendToChannel("%s was lynched!", mostVoted)
		}
	}
	hasTerminated <- true
}

func (werewolf *Werewolf) getKillVote(player string) string {
	for {
		vote := werewolf.communication.RequestName(player, "%s, who do you want to kill?: ", player)

		if vote == player {
			werewolf.communication.SendToPlayer(player, "You cannot kill yourself, sorry.")
		} else if !werewolf.data.IsRole(vote, "villager") {
			werewolf.communication.SendToPlayer(player, "That is not a living villager!")
		} else {
			return vote
		}
	}
}

func (werewolf *Werewolf) getLynchVote(player string) string {
	for {
		vote := werewolf.communication.RequestName(player, "%s, who do you want to lynch?: ", player)

		if vote == player {
			werewolf.communication.SendToPlayer(player, "You cannot lynch yourself, sorry.")
		} else if !werewolf.data.IsPlayer(vote) {
			werewolf.communication.SendToPlayer(player, "That is not a living player!")
		} else if werewolf.data.IsRole(vote, "werewolf") {
			werewolf.communication.SendToPlayer(player, "You cannot lynch another werewolf!")
		} else {
			return vote
		}
	}
}
