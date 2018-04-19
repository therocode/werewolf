package components

type Vote struct {
	componentName string
	vote          map[string]int
}

func NewVote(name string) *Vote {
	return &Vote{name, map[string]int{}}
}

func (this *Vote) Name() string {
	return this.componentName
}

func (this *Vote) Vote(name string) {
	this.vote[name]++
}

func (this *Vote) Reset() {
	this.vote = map[string]int{}
}

func (this *Vote) TotalVoteCount() int {
	totalVoteCount := 0
	for _, count := range this.vote {
		totalVoteCount += count
	}
	return totalVoteCount
}

func (this *Vote) MostVoted() string {
	maxVoteCount := 0
	var mostVoted string
	for name, count := range this.vote {
		if count > maxVoteCount {
			maxVoteCount = count
			mostVoted = name
		}
	}

	return mostVoted
}
