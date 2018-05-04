package components

import (
	"log"

	"github.com/therocode/werewolf/logic"
)

// Lynch is a component containing lynching functionality
type Lynch struct {
	data          logic.Data
	communication logic.Communication
	vote          *Vote
}

func NewLynch(data logic.Data, communication logic.Communication) *Lynch {
	return &Lynch{data, communication, NewVote("lynchVote")}
}

func (*Lynch) Name() string {
	return "lynch"
}

// Reset the voting ballot to empty
func (lynch *Lynch) Reset() {
	lynch.vote.Reset()
}

// Handle handles a lynch event for a role
func (lynch *Lynch) Handle(player string) {
	vote, timeout := lynch.requestVote(player)

	if timeout {
		lynch.communication.SendToPlayer(player, "Sorry, you took too long to decide.")
		lynch.communication.SendToChannel("%s took too long to decide and forfeited their vote", player)
		lynch.vote.VoteBlank()
	} else {
		lynch.communication.SendToChannel("%s voted to lynch %s", player, vote)
		lynch.vote.Vote(vote)
	}

	totalLynchVoteCount := lynch.vote.TotalVoteCount()
	neededLynchVotes := lynch.data.CountComponent(lynch.vote)
	log.Printf("%d people voted, need %d votes", totalLynchVoteCount, neededLynchVotes)
	if totalLynchVoteCount == neededLynchVotes {
		mostVoted, noVotes := lynch.vote.MostVoted()

		if noVotes {
			lynch.communication.SendToChannel("The villagers were so indecisive, nobody was lynched today!")
		} else {
			lynch.data.Kill(mostVoted)
			lynch.communication.SendToChannel("%s was lynched!", mostVoted)
		}
	}
}

func (lynch *Lynch) requestVote(player string) (string, bool) {
	for {
		lynch.data.Unlock()
		vote, timeout := lynch.communication.Request(player, "%s, who do you want to lynch?: ", player)
		lynch.data.Lock()

		switch {
		case timeout:
			return "", true
		case vote == player:
			lynch.communication.SendToPlayer(player, "You cannot lynch yourself, sorry.")
		case !lynch.data.IsPlayer(vote):
			lynch.communication.SendToPlayer(player, "That is not a living player!")
		default:
			return vote, false
		}
	}
}
