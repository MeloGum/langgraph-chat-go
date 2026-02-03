package llm

import (
    "context"
    "fmt"
    
    "github.com/cloudwego/eino/components/llm"
    "github.com/cloudwego/eino/components/prompt"
)

// DeepSeekClient wraps the LLM client for DeepSeek API
type DeepSeekClient struct {
    client llm.ChatModel
    prompt *prompt.ChatPrompt
}

// NewDeepSeekClient creates a new DeepSeek client
func NewDeepSeekClient(apiKey, baseURL, model string) (*DeepSeekClient, error) {
    client, err := llm.NewChatModel(context.Background(), &llm.ChatModelConfig{
        APIKey:  apiKey,
        BaseURL: baseURL,
        Model:   model,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create DeepSeek client: %w", err)
    }
    
    return &DeepSeekClient{
        client: client,
        prompt: prompt.NewChatPrompt(),
    }, nil
}

// GenerateResponse generates a response given a conversation
func (c *DeepSeekClient) GenerateResponse(ctx context.Context, messages []Message) (string, error) {
    msgs := make([]*llm.Message, len(messages))
    for i, m := range messages {
        msgs[i] = &llm.Message{
            Role:    llm.RoleType(m.Role),
            Content: m.Content,
        }
    }
    
    resp, err := c.client.Generate(ctx, msgs)
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}

// GenerateWithSystemPrompt generates response with system prompt
func (c *DeepSeekClient) GenerateWithSystemPrompt(ctx context.Context, systemPrompt, userMessage string) (string, error) {
    msgs := []*llm.Message{
        {Role: llm.System, Content: systemPrompt},
        {Role: llm.User, Content: userMessage},
    }
    
    resp, err := c.client.Generate(ctx, msgs)
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}

// Message represents a chat message
type Message struct {
    Role    string // "system", "user", "assistant"
    Content string
}

// BuildVotePrompt builds a prompt for voting
func BuildVotePrompt(agentName, gameState string, history []string) string {
    prompt := fmt.Sprintf(`You are %s. Current game state:
%s

Recent voting history:
`, agentName, gameState)
    
    for _, h := range history {
        prompt += fmt.Sprintf("- %s\n", h)
    }
    
    prompt += fmt.Sprintf(`Based on this information, who should you vote to attack? 
Respond with ONLY the agent name you choose to attack.`)
    
    return prompt
}

// BuildDialoguePrompt builds a prompt for dialogue
func BuildDialoguePrompt(agentName, context string) string {
    return fmt.Sprintf(`You are %s. 
Current context:
%s

Write a brief, strategic message to other agents. Consider:
- Your survival strategy
- Coalition building
- Misinformation when appropriate

Keep it concise and strategic.`, agentName, context)
}
