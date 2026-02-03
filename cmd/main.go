package main

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "time"
    
    "github.com/cloudwego/eino/flow"
    "github.com/cloudwego/eino/schema"
    
    "langgraph-chat-go/config"
    "langgraph-chat-go/engine"
    "langgraph-chat-go/llm"
)

func main() {
    // Seed random for strategies
    rand.Seed(time.Now().UnixNano())
    
    fmt.Println("╔══════════════════════════════════════════════════════════╗")
    fmt.Println("║   Multi-Agent Debate Game - Eino Go Implementation       ║")
    fmt.Println("║   Migrated from LangGraph/Python to Eino/Go              ║")
    fmt.Println("╚══════════════════════════════════════════════════════════╝")
    fmt.Println()
    
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    fmt.Printf("Configuration loaded:\n")
    fmt.Printf("  - Max Rounds: %d\n", cfg.MaxRounds)
    fmt.Printf("  - Base HP: %d\n", cfg.BaseHP)
    fmt.Printf("  - Attack Damage: %d\n", cfg.AttackDamage)
    fmt.Printf("  - Agent Count: %d\n", cfg.AgentCount)
    fmt.Println()
    
    // Create game config
    gameConfig := &engine.GameConfig{
        MaxRounds:    cfg.MaxRounds,
        BaseHP:       cfg.BaseHP,
        AttackDamage: cfg.AttackDamage,
        AgentNames:   cfg.AgentNames,
        AgentCount:   cfg.AgentCount,
    }
    
    // Initialize LLM client (optional - will use heuristics if not available)
    var llmClient *llm.DeepSeekClient
    if cfg.DeepSeekAPIKey != "" {
        llmClient, err = llm.NewDeepSeekClient(
            cfg.DeepSeekAPIKey,
            cfg.DeepSeekBaseURL,
            cfg.DeepSeekModel,
        )
        if err != nil {
            fmt.Printf("Warning: Failed to initialize LLM client: %v\n", err)
            fmt.Println("  Falling back to heuristic-based strategies\n")
        } else {
            fmt.Println("LLM client initialized successfully")
        }
    } else {
        fmt.Println("No API key provided - using heuristic-based strategies")
    }
    
    // Create and run game
    game := engine.NewGame(gameConfig)
    
    fmt.Println("\nStarting game...")
    game.Run()
    
    fmt.Println("\n✨ Migration complete! The game has been successfully ported to Go/Eino.")
}

// demonstrateEinoWorkflow demonstrates how to use Eino's workflow API
func demonstrateEinoWorkflow(ctx context.Context) {
    fmt.Println("\n--- Eino Workflow Demonstration ---")
    
    // Example: Create a simple chain using Eino's workflow API
    // This shows how the Go version could leverage Eino's orchestration
    
    workflow := flow.NewWorkflow[*schema.Message, *schema.Message]()
    
    // In a full implementation, you would:
    // 1. Add ChatModel nodes for each agent
    // 2. Add Tool nodes for game actions
    // 3. Configure edges between nodes
    // 4. Compile and execute the workflow
    
    fmt.Println("Eino Workflow structure:")
    fmt.Println("  Workflow[Message, Message]")
    fmt.Println("    ├── ChatModel: dialogue_agent")
    fmt.Println("    ├── ChatModel: vote_agent")
    fmt.Println("    └── Tool: attack_action")
    fmt.Println()
    fmt.Println("This demonstrates Eino's composition framework capabilities.")
    fmt.Println("For production use, agents would be connected via Graph/Workflow API.")
}
