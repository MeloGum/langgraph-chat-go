package engine

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

// Game represents the main game engine
type Game struct {
    Agents      []*Agent
    Rounds      int
    MaxRounds   int
    BaseHP      int
    AttackDamage int
    VoteHistory []VoteRecord
    Logs        []GameLog
    Config      *GameConfig
}

// GameConfig holds game configuration
type GameConfig struct {
    MaxRounds     int
    BaseHP        int
    AttackDamage  int
    AgentNames    []string
    AgentCount    int
}

// GameLog represents a single round's log
type GameLog struct {
    Round      int
    Dialogue   map[string]string
    Votes      map[string]string
    Attack     []AttackResult
    Survivors  []string
    Deaths     []string
    Timestamp  string
}

// AttackResult represents the result of an attack
type AttackResult struct {
    Attacker   string
    Target     string
    Damage     int
    VoterCount int
}

// NewGame creates a new game instance
func NewGame(cfg *GameConfig) *Game {
    game := &Game{
        Agents:      make([]*Agent, 0),
        Rounds:      0,
        MaxRounds:   cfg.MaxRounds,
        BaseHP:      cfg.BaseHP,
        AttackDamage: cfg.AttackDamage,
        VoteHistory: make([]VoteRecord, 0),
        Logs:        make([]GameLog, 0),
        Config:      cfg,
    }
    
    // Create agents based on config
    strategies := []StrategyType{
        EntropyShift,     // Alpha
        GrudgeKeeper,     // Beta
        PoisonLeader,     // Gamma
        SurvivorHunter,   // Delta
        BandwagonBreaker, // Epsilon
    }
    
    for i := 0; i < cfg.AgentCount; i++ {
        name := fmt.Sprintf("Agent_%d", i+1)
        if i < len(cfg.AgentNames) {
            name = cfg.AgentNames[i]
        }
        ag := NewAgent(name, strategies[i%len(strategies)], cfg.BaseHP)
        game.Agents = append(game.Agents, ag)
    }
    
    return game
}

// Run executes the game for max rounds
func (g *Game) Run() {
    fmt.Println("=== Starting Multi-Agent Debate Game ===")
    fmt.Println()
    
    for g.Rounds < g.MaxRounds {
        g.playRound()
        g.printStatus()
        
        if g.isGameOver() {
            break
        }
    }
    
    g.printFinalResults()
    g.saveLogs("game_log.json")
}

// playRound executes a single round
func (g *Game) playRound() {
    g.Rounds++
    fmt.Printf("\n--- Round %d ---\n", g.Rounds)
    
    log := GameLog{
        Round:      g.Rounds,
        Dialogue:   make(map[string]string),
        Votes:      make(map[string]string),
        Attack:     make([]AttackResult, 0),
        Timestamp:  time.Now().Format(time.RFC3339),
    }
    
    // Phase 1: Dialogue
    fmt.Println("Phase 1: Dialogue")
    for _, ag := range g.Agents {
        if ag.IsAlive {
            dialogue := ag.GetDialogue(g.Rounds, g.getGameState())
            log.Dialogue[ag.Name] = dialogue
            fmt.Printf("  [%s]: %s\n", ag.Name, dialogue)
        }
    }
    
    // Phase 2: Voting
    fmt.Println("\nPhase 2: Voting")
    votes := make(map[string]int)
    voteReasons := make(map[string]string)
    
    for _, ag := range g.Agents {
        if ag.IsAlive {
            target := ag.SelectVote(g.Agents, g.VoteHistory)
            if target != "" {
                votes[target]++
                reason := g.buildVoteReason(ag.Name, target)
                voteReasons[target] = reason
                log.Votes[ag.Name] = target
                
                record := VoteRecord{
                    Round:  g.Rounds,
                    Voter:  ag.Name,
                    Target: target,
                    Reason: reason,
                }
                g.VoteHistory = append(g.VoteHistory, record)
            }
        }
    }
    
    // Print voting results
    for voter, target := range log.Votes {
        fmt.Printf("  %s → %s\n", voter, target)
    }
    
    // Phase 3: Attack resolution
    fmt.Println("\nPhase 3: Attack Resolution")
    for target, count := range votes {
        damage := count * g.AttackDamage
        
        // Find target agent
        for _, ag := range g.Agents {
            if ag.Name == target && ag.IsAlive {
                ag.TakeDamage(damage)
                ag.VotesReceived += count
                
                attack := AttackResult{
                    Attacker:   "multiple",
                    Target:     target,
                    Damage:     damage,
                    VoterCount: count,
                }
                log.Attack = append(log.Attack, attack)
                
                fmt.Printf("  %s attacked! Votes: %d, Damage: %d, HP: %d → %d\n",
                    target, count, damage, ag.HP+damage, ag.HP)
                
                // Record death
                if !ag.IsAlive {
                    log.Deaths = append(log.Deaths, ag.Name)
                    fmt.Printf("  💀 %s has been eliminated!\n", ag.Name)
                }
                break
            }
        }
    }
    
    // Update survivors
    for _, ag := range g.Agents {
        if ag.IsAlive {
            log.Survivors = append(log.Survivors, ag.Name)
        }
    }
    
    g.Logs = append(g.Logs, log)
}

// buildVoteReason generates a reason for a vote
func (g *Game) buildVoteReason(voter, target string) string {
    // Generate a strategic reason based on game state
    return fmt.Sprintf("Strategic calculation: %s targets %s", voter, target)
}

// getGameState returns a summary of current game state
func (g *Game) getGameState() string {
    survivors := 0
    totalHP := 0
    for _, ag := range g.Agents {
        if ag.IsAlive {
            survivors++
            totalHP += ag.HP
        }
    }
    
    return fmt.Sprintf("Round %d/%d, %d survivors, avg HP: %d",
        g.Rounds, g.MaxRounds, survivors, totalHP/survivors)
}

// isGameOver checks if the game has ended
func (g *Game) isGameOver() bool {
    livingCount := 0
    for _, ag := range g.Agents {
        if ag.IsAlive {
            livingCount++
        }
    }
    return livingCount <= 1
}

// printStatus prints current game status
func (g *Game) printStatus() {
    fmt.Println("\n--- Current Status ---")
    for _, ag := range g.Agents {
        status := "alive"
        if !ag.IsAlive {
            status = "dead"
        }
        fmt.Printf("  %s: HP=%d, VotesReceived=%d [%s]\n",
            ag.Name, ag.HP, ag.VotesReceived, status)
    }
}

// printFinalResults prints final game results
func (g *Game) printFinalResults() {
    fmt.Println("\n" + strings.Repeat("=", 50))
    fmt.Println("🏆 GAME OVER - FINAL RESULTS 🏆")
    fmt.Println(strings.Repeat("=", 50))
    
    // Sort by HP (descending)
    sorted := make([]*Agent, len(g.Agents))
    copy(sorted, g.Agents)
    
    // Simple bubble sort by HP
    for i := 0; i < len(sorted); i++ {
        for j := i + 1; j < len(sorted); j++ {
            if sorted[j].HP > sorted[i].HP {
                sorted[i], sorted[j] = sorted[j], sorted[i]
            }
        }
    }
    
    for i, ag := range sorted {
        if ag.IsAlive {
            fmt.Printf("%d. %s - HP=%d (SURVIVED)\n", i+1, ag.Name, ag.HP)
        } else {
            fmt.Printf("%d. %s - HP=%d (eliminated, %d votes received)\n",
                i+1, ag.Name, ag.HP, ag.VotesReceived)
        }
    }
    
    // Winner
    for _, ag := range sorted {
        if ag.IsAlive {
            fmt.Printf("\n🎉 WINNER: %s (Strategy: %s) 🎉\n", ag.Name, ag.Strategy)
            break
        }
    }
}

// saveLogs saves game logs to a JSON file
func (g *Game) saveLogs(filename string) {
    output := struct {
        GameConfig   GameConfig
        FinalState   []AgentState
        GameLogs     []GameLog
        VoteHistory  []VoteRecord
    }{
        GameConfig:  *g.Config,
        FinalState:  g.getAgentStates(),
        GameLogs:    g.Logs,
        VoteHistory: g.VoteHistory,
    }
    
    data, _ := json.MarshalIndent(output, "", "  ")
    _ = os.WriteFile(filename, data, 0644)
    fmt.Printf("\n📁 Game logs saved to: %s\n", filename)
}

// AgentState represents the final state of an agent
type AgentState struct {
    Name          string
    Strategy      string
    HP            int
    VotesReceived int
    IsAlive       bool
    History       []string
}

// getAgentStates returns the final states of all agents
func (g *Game) getAgentStates() []AgentState {
    states := make([]AgentState, len(g.Agents))
    for i, ag := range g.Agents {
        states[i] = AgentState{
            Name:          ag.Name,
            Strategy:      string(ag.Strategy),
            HP:            ag.HP,
            VotesReceived: ag.VotesReceived,
            IsAlive:       ag.IsAlive,
            History:       ag.History,
        }
    }
    return states
}

// Import strings for final results
import "strings"
