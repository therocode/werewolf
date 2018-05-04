package components

// Vote is a component containing voting functionality
type Vote struct {
	componentName  string
	vote           map[string]int
	blankVoteCount int
}

func NewVote(name string) *Vote {
	return &Vote{name, map[string]int{}, 0}
}

func (vote *Vote) Name() string {
	return vote.componentName
}

// Vote for a particular name
func (vote *Vote) Vote(name string) {
	vote.vote[name]++
}

// VoteBlank i. e. increase total vote count but don't vote for a particular name
func (vote *Vote) VoteBlank() {
	vote.blankVoteCount++
}

// Reset the ballot to empty
func (vote *Vote) Reset() {
	vote.vote = map[string]int{}
	vote.blankVoteCount = 0
}

// TotalVoteCount returns the total number of votes cast so far
func (vote *Vote) TotalVoteCount() int {
	totalVoteCount := vote.blankVoteCount
	for _, count := range vote.vote {
		totalVoteCount += count
	}
	return totalVoteCount
}

// MostVoted returns the name with the most votes. Second parameter is true if no votes were cast for anyone.
func (vote *Vote) MostVoted() (string, bool) {
	maxVoteCount := 0
	var mostVoted string
	for name, count := range vote.vote {
		if count > maxVoteCount {
			maxVoteCount = count
			mostVoted = name
		}
	}

	return mostVoted, maxVoteCount == 0
}
