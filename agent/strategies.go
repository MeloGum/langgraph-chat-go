package agent

import (
    "fmt"
    "math/rand"
    "strings"
)

// StrategyType defines the type of agent strategy
type StrategyType string

const (
    EntropyShift      StrategyType = "entropy_shift"      // Unpredictable, hard to model
    GrudgeKeeper      StrategyType = "grudge_keeper"      // Attacks those who attacked them
    PoisonLeader      StrategyType = "poison_leader"      // Targets the strongest
    SurvivorHunter    StrategyType = "survivor_hunter"    // Stays low, attacks when safe
    BandwagonBreaker  StrategyType = "bandwagon_breaker"  // Goes against the majority
)

// Agent represents a game agent with a specific strategy
type Agent struct {
    Name       string
    Strategy   StrategyType
    HP         int
    VotesReceived int
    IsAlive    bool
    History    []string
}

// NewAgent creates a new agent with the given strategy
func NewAgent(name string, strategy StrategyType, maxHP int) *Agent {
    return &Agent{
        Name:         name,
        Strategy:     strategy,
        HP:           maxHP,
        VotesReceived: 0,
        IsAlive:      true,
        History:      make([]string, 0),
    }
}

// SelectVote selects a target based on the agent's strategy
func (a *Agent) SelectVote(agents []*Agent, history []VoteRecord) string {
    if !a.IsAlive {
        return ""
    }
    
    switch a.Strategy {
    case EntropyShift:
        return a.entropyShiftVote(agents, history)
    case GrudgeKeeper:
        return a.grudgeKeeperVote(agents, history)
    case PoisonLeader:
        return a.poisonLeaderVote(agents, history)
    case SurvivorHunter:
        return a.survivorHunterVote(agents, history)
    case BandwagonBreaker:
        return a.bandwagonBreakerVote(agents, history)
    default:
        return a.randomVote(agents)
    }
}

// GetStrategyPrompt returns the system prompt for this strategy
func (a *Agent) GetStrategyPrompt() string {
    prompts := map[StrategyType]string{
        EntropyShift: `You are unpredictable. Your voting patterns are chaotic and hard to model. 
You sometimes vote randomly to confuse other agents. You believe entropy is your ally.`,

        GrudgeKeeper: `You remember those who attack you. You seek revenge against agents who have targeted you in the past. 
Your motto: "An eye for an eye." You hold grudges and prioritize payback.`,

        PoisonLeader: `You believe in eliminating the strongest threat. You target agents with high HP or those leading the game.
Your philosophy: "Cut off the head of the snake."`,

        SurvivorHunter: `You prefer to stay under the radar. You attack only when you can safely do so without drawing too much attention.
Your strategy: "Live to fight another day." You target weak or isolated agents.`,

        BandwagonBreaker: `You go against the majority. If everyone is voting for X, you vote for someone else.
Your philosophy: "The crowd is often wrong." You value contrarian thinking.`,
    }
    
    if p, ok := prompts[a.Strategy]; ok {
        return p
    }
    return "You are a strategic agent trying to survive."
}

// GetDialogue generates dialogue based on the strategy
func (a *Agent) GetDialogue(round int, gameState string) string {
    templates := map[StrategyType][]string{
        EntropyShift: {
            "Chaos is my companion. Who can predict the unpredictable?",
            "Entropy increases. Patterns dissolve. I embrace the chaos.",
            "Prediction is futile. I am the variable you cannot model.",
        },
        GrudgeKeeper: {
            "I remember those who struck me. The ledger never forgets.",
            "Payback is inevitable. My patience is long.",
            "Those who hurt me will face consequences.",
        },
        PoisonLeader: {
            "The strongest must fall. No leader stands forever.",
            "Eliminate the threat before it eliminates us.",
            "The head of the snake must be severed.",
        },
        SurvivorHunter: {
            "I watch. I wait. I strike when the moment is right.",
            "Why fight when you can outlast?",
            "The patient hunter catches the prey.",
        },
        BandwagonBreaker: {
            "The crowd rushes one way. I go the other.",
            "Consensus is often wrong. Dissent is wisdom.",
            "Why follow when you can lead... in the opposite direction?",
        },
    }
    
    templates := templates[a.Strategy]
    return templates[round%len(templates)]
}

// Strategy implementations
func (a *Agent) entropyShiftVote(agents []*Agent, history []VoteRecord) string {
    // Sometimes vote randomly to be unpredictable
    if rand.Float32() < 0.3 {
        return a.randomVote(agents)
    }
    // Otherwise, target based on recent voting patterns
    return a.analyzeVotingPatterns(agents, history)
}

func (a *Agent) grudgeKeeperVote(agents []*Agent, history []VoteRecord) string {
    // Find who has voted for us most
    offenderCount := make(map[string]int)
    for _, v := range history {
        if v.Target == a.Name && v.Voter != a.Name {
            offenderCount[v.Voter]++
        }
    }
    
    // Vote for the biggest offender
    maxCount := 0
    target := ""
    for voter, count := range offenderCount {
        if count > maxCount {
            maxCount = count
            target = voter
        }
    }
    
    if target != "" {
        return target
    }
    return a.randomVote(agents)
}

func (a *Agent) poisonLeaderVote(agents []*Agent, history []VoteRecord) string {
    // Find agent with most votes received
    voteCount := make(map[string]int)
    for _, v := range history {
        voteCount[v.Target]++
    }
    
    maxVotes := 0
    target := ""
    for name, count := range voteCount {
        if count > maxVotes {
            maxVotes = count
            target = name
        }
    }
    
    if target != "" {
        return target
    }
    return a.randomVote(agents)
}

func (a *Agent) survivorHunterVote(agents []*Agent, history []VoteRecord) string {
    // Target low-HP agents who aren't paying attention
    livingAgents := make([]*Agent, 0)
    for _, ag := range agents {
        if ag.IsAlive && ag.Name != a.Name {
            livingAgents = append(livingAgents, ag)
        }
    }
    
    if len(livingAgents) == 0 {
        return a.randomVote(agents)
    }
    
    // Sort by HP (ascending)
    // Simple selection: target the one with lowest HP
    target := livingAgents[0]
    for _, ag := range livingAgents[1:] {
        if ag.HP < target.HP {
            target = ag
        }
    }
    
    return target.Name
}

func (a *Agent) bandwagonBreakerVote(agents []*Agent, history []VoteRecord) string {
    // Count votes for each agent in recent history
    voteCount := make(map[string]int)
    for _, v := range history {
        voteCount[v.Target]++
    }
    
    // Find the most voted agent
    maxVotes := 0
    popularTarget := ""
    for target, count := range voteCount {
        if count > maxVotes {
            maxVotes = count
            popularTarget = target
        }
    }
    
    // If there's a clear favorite, vote for someone else
    if popularTarget != "" && maxVotes >= 2 {
        for _, ag := range agents {
            if ag.IsAlive && ag.Name != popularTarget && ag.Name != a.Name {
                return ag.Name
            }
        }
    }
    
    return a.randomVote(agents)
}

func (a *Agent) analyzeVotingPatterns(agents []*Agent, history []VoteRecord) string {
    // Analyze who is being targeted and by whom
    // This is a simple heuristic - in practice would use LLM
    return a.randomVote(agents)
}

func (a *Agent) randomVote(agents []*Agent) string {
    livingAgents := make([]string, 0)
    for _, ag := range agents {
        if ag.IsAlive && ag.Name != a.Name {
            livingAgents = append(livingAgents, ag.Name)
        }
    }
    
    if len(livingAgents) == 0 {
        return ""
    }
    
    return livingAgents[rand.Intn(len(livingAgents))]
}

// VoteRecord represents a single vote
type VoteRecord struct {
    Round     int
    Voter     string
    Target    string
    Reason    string
}

// TakeDamage reduces the agent's HP
func (a *Agent) TakeDamage(damage int) {
    a.HP -= damage
    if a.HP <= 0 {
        a.HP = 0
        a.IsAlive = false
    }
}

// AddHistory records an event in the agent's history
func (a *Agent) AddHistory(event string) {
    a.History = append(a.History, event)
}

// String returns a string representation of the agent
func (a *Agent) String() string {
    status := "alive"
    if !a.IsAlive {
        status = "dead"
    }
    return fmt.Sprintf("%s (%s) - HP: %d, Strategy: %s", a.Name, status, a.HP, a.Strategy)
}

// ParseStrategy parses a strategy string to StrategyType
func ParseStrategy(s string) StrategyType {
    s = strings.ToLower(strings.TrimSpace(s))
    switch s {
    case "entropy_shift":
        return EntropyShift
    case "grudge_keeper":
        return GrudgeKeeper
    case "poison_leader":
        return PoisonLeader
    case "survivor_hunter":
        return SurvivorHunter
    case "bandwagon_breaker":
        return BandwagonBreaker
    default:
        return EntropyShift
    }
}

// GetAllStrategies returns all available strategies
func GetAllStrategies() []StrategyType {
    return []StrategyType{
        EntropyShift,
        GrudgeKeeper,
        PoisonLeader,
        SurvivorHunter,
        BandwagonBreaker,
    }
}
