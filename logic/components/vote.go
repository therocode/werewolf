package components

type Vote struct {
	componentName  string
	vote           map[string]int
	blankVoteCount int
}

func NewVote(name string) *Vote {
	return &Vote{name, map[string]int{}, 0}
}

func (this *Vote) Name() string {
	return this.componentName
}

func (this *Vote) Vote(name string) {
	this.vote[name]++
}

func (this *Vote) VoteBlank() {
	this.blankVoteCount++
}

func (this *Vote) Reset() {
	this.vote = map[string]int{}
	this.blankVoteCount = 0
}

func (this *Vote) TotalVoteCount() int {
	totalVoteCount := this.blankVoteCount
	for _, count := range this.vote {
		totalVoteCount += count
	}
	return totalVoteCount
}

// MostVoted returns the name with the most votes. Second parameter is true if no votes were cast for anyone.
func (this *Vote) MostVoted() (string, bool) {
	maxVoteCount := 0
	var mostVoted string
	for name, count := range this.vote {
		if count > maxVoteCount {
			maxVoteCount = count
			mostVoted = name
		}
	}

	return mostVoted, maxVoteCount == 0
}
