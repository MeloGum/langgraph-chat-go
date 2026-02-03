# LangGraph Chat → Eino Go 迁移计划

## 任务分解

### 阶段1: 文档研究
- [ ] Eino 官方文档概览
- [ ] Agent Development Kit (ADK) 使用
- [ ] Graph/Workflow API 最佳实践
- [ ] 多 Agent 编排示例

### 阶段2: 核心代码
- [ ] config/config.go - 配置管理
- [ ] llm/deepseek.go - DeepSeek API 客户端
- [ ] agent/strategies.go - 5 种策略定义
- [ ] engine/game.go - 游戏引擎
- [ ] engine/vote.go - 投票系统
- [ ] graph/workflow.go - Eino 工作流编排
- [ ] cmd/main.go - 入口文件

### 阶段3: 测试与文档
- [ ] tests/ 单元测试
- [ ] README.md
- [ ] go.mod 依赖管理

## 架构设计

```
cmd/main.go
    ↓
graph/workflow.go (Eino Workflow)
    ├── agent/strategies.go (5 agents)
    ├── engine/game.go (game loop)
    ├── engine/vote.go (voting logic)
    └── llm/deepseek.go (LLM calls)
```

## 策略映射

| Python (LangChain) | Go (Eino) |
|-------------------|-----------|
| Agent | adk.ChatModelAgent |
| State | schema.Message |
| Chain | Chain API |
| Graph | Graph API |
| Tool | tool.BaseTool |
