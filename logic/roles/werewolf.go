package roles

import (
	"log"
	"strings"

	"github.com/therocode/werewolf/logic"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/timeline"
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
func (werewolf *Werewolf) HasComponent(component logic.Component) bool {
	return werewolf.lynchVote.Name() == component.Name() || werewolf.killVote.Name() == component.Name()
}

// Generate implements Role interface
func (*Werewolf) Generate() []timeline.Event {
	return []timeline.Event{
		timeline.Event{
			Name:   "werewolves_see_each_other",
			Before: map[string]bool{"night_starts": true},
			After:  map[string]bool{},
		},
		timeline.Event{
			Name:   "werewolves_kill",
			Before: map[string]bool{"werewolves_see_each_other": true},
			After:  map[string]bool{"day_starts": true},
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
	case "day_starts":
		werewolf.lynchVote.Reset()
	case "lynch":
		vote, timeout := werewolf.getLynchVote(player)

		if timeout {
			werewolf.communication.SendToChannel("%s took too long to decide and forfeited their vote", player)
			werewolf.lynchVote.VoteBlank()
		} else {
			werewolf.communication.SendToChannel("%s voted to lynch %s", player, vote)
			werewolf.lynchVote.Vote(vote)
		}

		totalLynchVoteCount := werewolf.lynchVote.TotalVoteCount()
		neededLynchVotes := werewolf.data.CountComponent(werewolf.lynchVote)
		log.Printf("%d people voted, need %d votes", totalLynchVoteCount, neededLynchVotes)
		if totalLynchVoteCount == neededLynchVotes {
			mostVoted, noVotes := werewolf.lynchVote.MostVoted()

			if noVotes {
				werewolf.communication.SendToChannel("The villagers were so indecisive, nobody was lynched today!")
			} else {
				werewolf.data.Kill(mostVoted)
				werewolf.communication.SendToChannel("%s was lynched!", mostVoted)
			}
		}
	}
	hasTerminated <- true
}

func (werewolf *Werewolf) getKillVote(player string) (string, bool) {
	for {
		vote, timeout := werewolf.communication.Request(player, "%s, who do you want to kill?: ", player)

		switch {
		case timeout:
			return "", true
		case vote == player:
			werewolf.communication.SendToPlayer(player, "You cannot kill yourself, sorry.")
		case !werewolf.data.IsRole(vote, "villager"):
			werewolf.communication.SendToPlayer(player, "That is not a living villager!")
		default:
			return vote, false
		}
	}
}

func (werewolf *Werewolf) getLynchVote(player string) (string, bool) {
	for {
		vote, timeout := werewolf.communication.Request(player, "%s, who do you want to lynch?: ", player)

		switch {
		case timeout:
			return "", true
		case vote == player:
			werewolf.communication.SendToPlayer(player, "You cannot lynch yourself, sorry.")
		case !werewolf.data.IsPlayer(vote):
			werewolf.communication.SendToPlayer(player, "That is not a living player!")
		case werewolf.data.IsRole(vote, "werewolf"):
			werewolf.communication.SendToPlayer(player, "You cannot lynch another werewolf!")
		default:
			return vote, false
		}
	}
}
