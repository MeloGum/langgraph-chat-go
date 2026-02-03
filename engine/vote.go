package engine

// VoteSystem handles the voting mechanism
type VoteSystem struct {
    Votes     map[string][]string // target -> [voters]
    Records   []VoteRecord
    Threshold int // votes needed to trigger attack
}

// NewVoteSystem creates a new vote system
func NewVoteSystem() *VoteSystem {
    return &VoteSystem{
        Votes:     make(map[string][]string),
        Records:   make([]VoteRecord, 0),
        Threshold: 1,
    }
}

// CastVote registers a vote from a voter for a target
func (vs *VoteSystem) CastVote(voter, target string) {
    vs.Votes[target] = append(vs.Votes[target], voter)
    
    record := VoteRecord{
        Voter:  voter,
        Target: target,
    }
    vs.Records = append(vs.Records, record)
}

// TallyVotes returns the vote counts
func (vs *VoteSystem) TallyVotes() map[string]int {
    counts := make(map[string]int)
    for target, voters := range vs.Votes {
        counts[target] = len(voters)
    }
    return counts
}

// GetVotersFor returns the list of voters for a target
func (vs *VoteSystem) GetVotersFor(target string) []string {
    return vs.Votes[target]
}

// Reset clears all votes
func (vs *VoteSystem) Reset() {
    vs.Votes = make(map[string][]string)
}
