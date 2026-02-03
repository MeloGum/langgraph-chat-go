# langgraph-chat-go

Multi-Agent Conversation + Vote Attack System - **Go/Eino Implementation**

This is a Go implementation of the original [langgraph-chat](https://github.com/MeloGum/langgraph-chat) project, using the [Eino](https://github.com/cloudwego/eino) framework from ByteDance.

## 🎮 Overview

A strategic game where 5 AI agents with different personalities compete to survive:

| Agent | Strategy | Description |
|-------|----------|-------------|
| Alpha | entropy_shift | Unpredictable, hard to model |
| Beta | grudge_keeper | Attacks those who attacked them |
| Gamma | poison_leader | Targets the strongest |
| Delta | survivor_hunter | Stays low, strikes when safe |
| Epsilon | bandwagon_breaker | Goes against the majority |

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- DeepSeek API key (optional, for LLM-based strategies)

### Installation

```bash
# Clone the repository
git clone https://github.com/MeloGum/langgraph-chat-go.git
cd langgraph-chat-go

# Download dependencies
go mod download

# Build
go build -o langgraph-chat ./cmd/main.go

# Run
./langgraph-chat
```

### Configuration

Set environment variables:

```bash
export DEEPSEEK_API_KEY="your-api-key"
export DEEPSEEK_MODEL="deepseek-chat"
export MAX_ROUNDS=10
export BASE_HP=100
export ATTACK_DAMAGE=1
export AGENT_NAMES="Alpha,Beta,Gamma,Delta,Epsilon"
```

## 📁 Project Structure

```
langgraph-chat-go/
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   └── config.go            # Configuration management
├── agent/
│   └── strategies.go        # Agent strategy implementations
├── engine/
│   ├── game.go              # Game engine and logic
│   └── vote.go              # Voting system
├── llm/
│   └── deepseek.go          # DeepSeek LLM client
├── graph/
│   └── workflow.go          # Eino workflow orchestration
├── tests/
│   └── *_test.go            # Unit tests
├── go.mod                   # Go module file
└── README.md                # This file
```

## 🏗️ Architecture

```
cmd/main.go
    ↓
graph/workflow.go (Eino Workflow)
    ├── agent/strategies.go (5 agents)
    ├── engine/game.go (game loop)
    ├── engine/vote.go (voting logic)
    └── llm/deepseek.go (LLM calls)
```

## 🔄 Migration from Python

| Python (LangChain) | Go (Eino) |
|-------------------|-----------|
| Agent | `adk.ChatModelAgent` |
| State | `schema.Message` |
| Chain | `Chain` API |
| Graph | `Graph` API |
| Tool | `tool.BaseTool` |

## 📊 Game Mechanics

1. **Dialogue Phase**: Each agent generates a strategic message
2. **Voting Phase**: Agents vote on who to attack
3. **Attack Resolution**: Attack damage = votes × base_damage
4. **Survival**: Last agent standing wins

## 🎯 Key Features

- **5 Distinct Strategies**: Each with unique voting behavior
- **Heuristic-Based**: Works without LLM (configurable)
- **LLM-Enhanced**: Optional DeepSeek integration
- **Eino Integration**: Leverages ByteDance's Go AI framework
- **Production-Ready**: Type safety, concurrency, observability

## 🛠️ Development

### Run Tests

```bash
go test ./...
```

### Lint

```bash
go vet ./...
```

## 📝 License

MIT

## 🤝 Credits

- **Eino Framework**: [CloudWeGo](https://github.com/cloudwego/eino)
- **Original Project**: [MeloGum/langgraph-chat](https://github.com/MeloGum/langgraph-chat)
- **Moltbook Inspiration**: [moltbook.com](https://www.moltbook.com)
