package config

import (
    "os"
    "strconv"
    "strings"
)

// Config holds all configuration for the application
type Config struct {
    // DeepSeek LLM Configuration
    DeepSeekAPIKey string
    DeepSeekBaseURL string
    DeepSeekModel string
    
    // Game Configuration
    MaxRounds int
    BaseHP int
    AttackDamage int
    
    // Agent Configuration
    SystemPrompt string
    AgentNames []string
    AgentCount int
    
    // Eino Configuration
    EnableTracing bool
    LogLevel string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    cfg := &Config{
        DeepSeekAPIKey:     getEnv("DEEPSEEK_API_KEY", ""),
        DeepSeekBaseURL:    getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com"),
        DeepSeekModel:      getEnv("DEEPSEEK_MODEL", "deepseek-chat"),
        MaxRounds:          getEnvInt("MAX_ROUNDS", 10),
        BaseHP:             getEnvInt("BASE_HP", 100),
        AttackDamage:       getEnvInt("ATTACK_DAMAGE", 1),
        SystemPrompt:       getEnv("SYSTEM_PROMPT", defaultSystemPrompt()),
        AgentNames:         parseAgentNames(getEnv("AGENT_NAMES", "Alpha,Beta,Gamma,Delta,Epsilon")),
        AgentCount:         5,
        EnableTracing:      getEnvBool("ENABLE_TRACING", true),
        LogLevel:           getEnv("LOG_LEVEL", "info"),
    }
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if v := os.Getenv(key); v != "" {
        if i, err := strconv.Atoi(v); err == nil {
            return i
        }
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if v := os.Getenv(key); v != "" {
        return strings.ToLower(v) == "true"
    }
    return defaultValue
}

func parseAgentNames(s string) []string {
    if s == "" {
        return []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"}
    }
    parts := strings.Split(s, ",")
    names := make([]string, 0, len(parts))
    for _, p := range parts {
        if trimmed := strings.TrimSpace(p); trimmed != "" {
            names = append(names, trimmed)
        }
    }
    return names
}

func defaultSystemPrompt() string {
    return `You are an AI agent participating in a strategic game. 
You will engage in dialogue with other agents and vote on who to attack.
Your goal is to survive as long as possible.

Consider:
- What other agents might do
- Their voting patterns
- Strategic positioning
- Coalition building

Respond concisely and strategically.`
}
