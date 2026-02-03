package graph

import (
    "context"
    
    "github.com/cloudwego/eino/flow"
    "github.com/cloudwego/eino/schema"
    
    "langgraph-chat-go/agent"
)

// GameWorkflow represents the Eino workflow for the game
type GameWorkflow struct {
    workflow *flow.Workflow[*schema.Message, *schema.Message]
    agents   map[string]*agent.Agent
}

// NewGameWorkflow creates a new game workflow
func NewGameWorkflow(agents []*agent.Agent) *GameWorkflow {
    wf := flow.NewWorkflow[*schema.Message, *schema.Message]()
    
    agentMap := make(map[string]*agent.Agent)
    for _, ag := range agents {
        agentMap[ag.Name] = ag
    }
    
    return &GameWorkflow{
        workflow: wf,
        agents:   agentMap,
    }
}

// Build constructs the workflow graph
func (gw *GameWorkflow) Build() error {
    // Add dialogue generation node for each agent
    for name := range gw.agents {
        gw.workflow.AddLambdaNode(
            "dialogue_"+name,
            gw.createDialogueNode(name),
        )
    }
    
    // Add voting node for each agent
    for name := range gw.agents {
        gw.workflow.AddLambdaNode(
            "vote_"+name,
            gw.createVoteNode(name),
        )
    }
    
    // Add attack resolution node
    gw.workflow.AddLambdaNode(
        "attack_resolution",
        gw.createAttackNode(),
    )
    
    return nil
}

// createDialogueNode creates a lambda node for dialogue generation
func (gw *GameWorkflow) createDialogueNode(agentName string) 
    func(ctx context.Context, input *schema.Message) (*schema.Message, error) {
    
    return func(ctx context.Context, input *schema.Message) (*schema.Message, error) {
        ag := gw.agents[agentName]
        dialogue := ag.GetDialogue(0, "")
        return schema.UserMessage(dialogue), nil
    }
}

// createVoteNode creates a lambda node for voting
func (gw *GameWorkflow) createVoteNode(agentName string)
    func(ctx context.Context, input *schema.Message) (*schema.Message, error) {
    
    return func(ctx context.Context, input *schema.Message) (*schema.Message, error) {
        ag := gw.agents[agentName]
        target := ag.SelectVote(nil, nil)
        return schema.UserMessage(target), nil
    }
}

// createAttackNode creates a lambda node for attack resolution
func (gw *GameWorkflow) createAttackNode()
    func(ctx context.Context, input *schema.Message) (*schema.Message, error) {
    
    return func(ctx context.Context, input *schema.Message) (*schema.Message, error) {
        // Attack resolution logic
        return schema.UserMessage("Attack resolved"), nil
    }
}

// Execute runs the workflow
func (gw *GameWorkflow) Execute(ctx context.Context, input *schema.Message) error {
    compiled, err := gw.workflow.Compile(ctx)
    if err != nil {
        return err
    }
    _, err = compiled.Invoke(ctx, input)
    return err
}

// GetWorkflow returns the underlying workflow
func (gw *GameWorkflow) GetWorkflow() *flow.Workflow[*schema.Message, *schema.Message] {
    return gw.workflow
}
