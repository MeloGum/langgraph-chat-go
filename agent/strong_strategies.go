package agent

import (
    "sort"
)

// New strategies to counter Poison Leader

// GlassJawFeint - Stay below strongest, finish what poison_leader starts
const GlassJawFeint StrategyType = "glass_jaw_feint"

// DecoyCrown - Maintain a decoy leader to absorb poison_leader's votes
const DecoyCrown StrategyType = "decoy_crown"

// TruceBroker - Coalition building, anti-poison-leader majority
const TruceBroker StrategyType = "truce_broker"

// ReverseShield - Target poison_leader directly
const ReverseShield StrategyType = "reverse_shield"

// KingpinTrap - Keep strongest alive to absorb poison_leader's votes
const KingpinTrap StrategyType = "kingpin_trap"

// Strategy state for new strategies
type StrongStrategyState struct {
    DecoyName    string
    LastX        string
    XStreak      int
    RoundCount   int
}

// glassJawFeintVote - If top-2 HP, vote strongest; otherwise vote who poison_leader targeted
func (a *Agent) glassJawFeintVote(agents []*Agent, history []VoteRecord) string {
    living := getLivingAgents(agents)
    if len(living) < 2 {
        return a.randomVote(agents, history)
    }
    
    // Sort by HP descending
    sorted := make([]*Agent, len(living))
    copy(sorted, living)
    sort.Slice(sorted, func(i, j int) bool {
        if sorted[i].HP == sorted[j].HP {
            return sorted[i].Name < sorted[j].Name
        }
        return sorted[i].HP > sorted[j].HP
    })
    
    // If I'm top-2, vote strongest (avoid being target)
    top2 := make(map[string]bool)
    if len(sorted) >= 1 {
        top2[sorted[0].Name] = true
    }
    if len(sorted) >= 2 {
        top2[sorted[1].Name] = true
    }
    
    if top2[a.Name] {
        return sorted[0].Name // Vote strongest to keep game moving
    }
    
    // Find who poison_leader targeted most in last 2 rounds
    poisonTargets := make(map[string]int)
    roundMap := make(map[int][]string)
    for _, v := range history {
        roundMap[v.Round] = append(roundMap[v.Round], v.Target)
    }
    
    maxRound := 0
    for r := range roundMap {
        if r > maxRound {
            maxRound = r
        }
    }
    
    for r := maxRound; r > maxRound-2 && r > 0; r-- {
        for _, target := range roundMap[r] {
            // Simple heuristic: assume highest HP targets are poison_leader's targets
            for _, ag := range sorted {
                if ag.Name == target {
                    poisonTargets[target]++
                    break
                }
            }
        }
    }
    
    if len(poisonTargets) > 0 {
        // Vote the target with most "poison" association
        for _, ag := range sorted {
            if _, ok := poisonTargets[ag.Name]; ok && ag.Name != a.Name {
                return ag.Name
            }
        }
    }
    
    return sorted[0].Name
}

// decoyCrownVote - Maintain a decoy leader, then kill them at right moment
func (a *Agent) decoyCrownVote(agents []*Agent, history []VoteRecord) string {
    living := getLivingAgents(agents)
    if len(living) < 2 {
        return a.randomVote(agents, history)
    }
    
    sorted := make([]*Agent, len(living))
    copy(sorted, living)
    sort.Slice(sorted, func(i, j int) bool {
        if sorted[i].HP == sorted[j].HP {
            return sorted[i].Name < sorted[j].Name
        }
        return sorted[i].HP > sorted[j].HP
    })
    
    // Pick decoy (initially strongest)
    decoy := sorted[0].Name
    
    // Count expected incoming votes (heuristic: assume 1-2 votes per strong target)
    // If decoy HP is low, switch to kill them
    
    // For simplicity: alternate between protecting decoy and attacking
    if a.RoundCount%3 == 0 {
        // Protect decoy: attack whoever threatens to overtake them
        for i := len(sorted) - 1; i >= 0; i-- {
            if sorted[i].Name != decoy && sorted[i].Name != a.Name {
                // Vote someone who might overtake decoy
                return sorted[i].Name
            }
        }
    }
    
    // Otherwise, attack weakest (keep decoy safe)
    return sorted[len(sorted)-1].Name
}

// truceBrokerVote - Build coalition, vote non-strongest majority target
func (a *Agent) truceBrokerVote(agents []*Agent, history []VoteRecord) string {
    living := getLivingAgents(agents)
    if len(living) < 2 {
        return a.randomVote(agents, history)
    }
    
    // Find target with >=2 votes last round (excluding strongest)
    voteCount := make(map[string]int)
    for _, v := range history {
        if v.Round == getLastRound(history) {
            voteCount[v.Target]++
        }
    }
    
    sorted := make([]*Agent, len(living))
    copy(sorted, living)
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].HP > sorted[j].HP
    })
    strongest := sorted[0].Name
    
    // Vote for non-strongest target with most votes
    bestTarget := ""
    bestVotes := 0
    for target, votes := range voteCount {
        if target != strongest && target != a.Name && votes > bestVotes {
            bestVotes = votes
            bestTarget = target
        }
    }
    
    if bestVotes >= 2 {
        return bestTarget
    }
    
    // Fallback: vote weakest among top-3
    if len(sorted) >= 3 {
        return sorted[len(sorted)-1].Name
    }
    
    return sorted[len(sorted)-1].Name
}

// reverseShieldVote - Target poison_leader directly
func (a *Agent) reverseShieldVote(agents []*Agent, history []VoteRecord) string {
    living := getLivingAgents(agents)
    if len(living) < 2 {
        return a.randomVote(agents, history)
    }
    
    sorted := make([]*Agent, len(living))
    copy(sorted, living)
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].HP > sorted[j].HP
    })
    
    // If I'm not strongest, vote poison_leader (whoever is targeting strongest)
    if sorted[0].Name != a.Name {
        // Simple: vote the second strongest (who is likely poison_leader's target)
        if len(sorted) >= 2 {
            return sorted[1].Name
        }
    }
    
    // If I am strongest, vote next-strongest to drop out of #1
    if len(sorted) >= 2 {
        return sorted[1].Name
    }
    
    return a.randomVote(agents, history)
}

// kingpinTrapVote - Keep strongest alive, then kill poison_leader
func (a *Agent) kingpinTrapVote(agents []*Agent, history []VoteRecord) string {
    living := getLivingAgents(agents)
    if len(living) < 2 {
        return a.randomVote(agents, history)
    }
    
    sorted := make([]*Agent, len(living))
    copy(sorted, living)
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].HP > sorted[j].HP
    })
    
    // Find who was targeted most last round (assume that's poison_leader's target)
    voteCount := make(map[string]int)
    lastRound := getLastRound(history)
    for _, v := range history {
        if v.Round == lastRound {
            voteCount[v.Target]++
        }
    }
    
    lastTarget := ""
    lastVotes := 0
    for target, votes := range voteCount {
        if votes > lastVotes {
            lastVotes = votes
            lastTarget = target
        }
    }
    
    // If poison_leader has been targeting the same X for 2+ rounds, vote poison_leader
    // Simplified: if last target exists and is strongest, vote them
    
    if lastTarget != "" && lastTarget != a.Name {
        // Check if this looks like poison_leader pattern (targeting strongest)
        if lastTarget == sorted[0].Name && lastVotes >= 1 {
            // Switch to vote the voter (poison_leader)
            for _, v := range history {
                if v.Round == lastRound && v.Target == lastTarget && v.Voter != a.Name {
                    // Vote the voter (simplified)
                    return v.Voter
                }
            }
        }
    }
    
    // Otherwise, vote someone who is NOT the strongest (keep kingpin alive)
    for _, ag := range living {
        if ag.Name != sorted[0].Name && ag.Name != a.Name {
            return ag.Name
        }
    }
    
    return sorted[0].Name
}

// getLivingAgents helper
func getLivingAgents(agents []*Agent) []*Agent {
    result := make([]*Agent, 0)
    for _, ag := range agents {
        if ag.IsAlive && ag != nil {
            result = append(result, ag)
        }
    }
    return result
}

// getLastRound helper
func getLastRound(history []VoteRecord) int {
    maxRound := 0
    for _, v := range history {
        if v.Round > maxRound {
            maxRound = v.Round
        }
    }
    return maxRound
}

// RoundCount tracks rounds for some strategies
type RoundCount int

// Add RoundCount to Agent struct in strategies.go
const RoundCountField = "RoundCount"
