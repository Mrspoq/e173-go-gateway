package voice

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/jackc/pgx/v4/pgxpool"
)

// DefaultActionEngine implements the ActionEngine interface
type DefaultActionEngine struct {
    db           *pgxpool.Pool
    aiRouter     AIRouter
    simManager   SIMManager
    logger       *log.Logger
}

// AIRouter handles routing calls to AI agents
type AIRouter interface {
    RouteCall(ctx context.Context, callID string, reason string) error
    GetAvailableAgents() []AIAgent
}

// AIAgent represents an AI voice agent
type AIAgent struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Endpoint    string `json:"endpoint"`
    Speciality  string `json:"speciality"`
    IsAvailable bool   `json:"is_available"`
}

// SIMManager handles SIM card operations
type SIMManager interface {
    FlagSIM(ctx context.Context, simID string, issue string) error
    GetSIMStatus(ctx context.Context, simID string) (*SIMStatus, error)
    ScheduleReplacement(ctx context.Context, simID string) error
}

// SIMStatus represents current SIM card status
type SIMStatus struct {
    SIMID      string    `json:"sim_id"`
    Status     string    `json:"status"`
    Issue      string    `json:"issue,omitempty"`
    FlaggedAt  time.Time `json:"flagged_at,omitempty"`
    ActionTaken string   `json:"action_taken,omitempty"`
}

// NewDefaultActionEngine creates a new action engine
func NewDefaultActionEngine(db *pgxpool.Pool, aiRouter AIRouter, simManager SIMManager, logger *log.Logger) *DefaultActionEngine {
    return &DefaultActionEngine{
        db:         db,
        aiRouter:   aiRouter,
        simManager: simManager,
        logger:     logger,
    }
}

// ExecuteAction performs the appropriate action based on classification
func (e *DefaultActionEngine) ExecuteAction(ctx context.Context, callID string, classification *Classification) error {
    e.logger.Printf("Executing action for call %s: Category=%s, Action=%s, Confidence=%.2f",
        callID, classification.Category, classification.Action, classification.Confidence)
    
    // Log the decision to database
    if err := e.logDecision(ctx, callID, classification); err != nil {
        e.logger.Printf("Failed to log decision: %v", err)
    }
    
    switch classification.Action {
    case ActionRouteToAI:
        return e.RouteToAI(ctx, callID)
        
    case ActionFlagSIM:
        // Extract SIM ID from call context (would be passed through context in real implementation)
        simID := "sim-placeholder" // This would come from call metadata
        return e.FlagSIM(ctx, simID, classification.Reason)
        
    case ActionBlockCall:
        return e.blockCall(ctx, callID, classification.Reason)
        
    case ActionRecordForReview:
        return e.markForReview(ctx, callID, classification)
        
    case ActionNormalRouting:
        // Normal routing - no special action needed
        e.logger.Printf("Call %s proceeding with normal routing", callID)
        return nil
        
    default:
        return fmt.Errorf("unknown action: %s", classification.Action)
    }
}

// RouteToAI sends the call to an AI agent for handling
func (e *DefaultActionEngine) RouteToAI(ctx context.Context, callID string) error {
    e.logger.Printf("Routing call %s to AI agent", callID)
    
    // Use the AI router to handle the call
    if err := e.aiRouter.RouteCall(ctx, callID, "Spam call detected"); err != nil {
        return fmt.Errorf("failed to route to AI: %w", err)
    }
    
    // Update call status in database
    query := `
        UPDATE sip_calls 
        SET status = 'routed_to_ai', 
            routed_at = NOW(),
            updated_at = NOW()
        WHERE call_id = $1
    `
    
    _, err := e.db.Exec(ctx, query, callID)
    if err != nil {
        e.logger.Printf("Failed to update call status: %v", err)
    }
    
    return nil
}

// FlagSIM marks a SIM card for attention
func (e *DefaultActionEngine) FlagSIM(ctx context.Context, simID string, reason string) error {
    e.logger.Printf("Flagging SIM %s: %s", simID, reason)
    
    // Use SIM manager to flag the SIM
    if err := e.simManager.FlagSIM(ctx, simID, reason); err != nil {
        return fmt.Errorf("failed to flag SIM: %w", err)
    }
    
    // Check if immediate action is needed
    if reason == "SIM card blocked by operator" {
        // Schedule replacement
        if err := e.simManager.ScheduleReplacement(ctx, simID); err != nil {
            e.logger.Printf("Failed to schedule SIM replacement: %v", err)
        }
    }
    
    return nil
}

// blockCall terminates a call
func (e *DefaultActionEngine) blockCall(ctx context.Context, callID string, reason string) error {
    e.logger.Printf("Blocking call %s: %s", callID, reason)
    
    // Update call status
    query := `
        UPDATE sip_calls 
        SET status = 'blocked',
            block_reason = $2,
            ended_at = NOW(),
            updated_at = NOW()
        WHERE call_id = $1
    `
    
    _, err := e.db.Exec(ctx, query, callID, reason)
    return err
}

// markForReview flags a call for manual review
func (e *DefaultActionEngine) markForReview(ctx context.Context, callID string, classification *Classification) error {
    e.logger.Printf("Marking call %s for review", callID)
    
    query := `
        INSERT INTO call_reviews (
            call_id, category, confidence, reason, 
            risk_score, status, created_at
        ) VALUES ($1, $2, $3, $4, $5, 'pending', NOW())
    `
    
    _, err := e.db.Exec(ctx, query,
        callID,
        string(classification.Category),
        classification.Confidence,
        classification.Reason,
        classification.RiskScore,
    )
    
    return err
}

// logDecision records the classification decision
func (e *DefaultActionEngine) logDecision(ctx context.Context, callID string, classification *Classification) error {
    query := `
        INSERT INTO voice_recognition_logs (
            call_id, category, action, confidence,
            reason, risk_score, keywords, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
    `
    
    _, err := e.db.Exec(ctx, query,
        callID,
        string(classification.Category),
        string(classification.Action),
        classification.Confidence,
        classification.Reason,
        classification.RiskScore,
        classification.Keywords,
    )
    
    return err
}

// SimpleAIRouter provides basic AI routing functionality
type SimpleAIRouter struct {
    agents []AIAgent
    logger *log.Logger
}

// NewSimpleAIRouter creates a basic AI router
func NewSimpleAIRouter(logger *log.Logger) *SimpleAIRouter {
    return &SimpleAIRouter{
        logger: logger,
        agents: []AIAgent{
            {
                ID:          "agent-001",
                Name:        "Sales Bot Handler",
                Endpoint:    "sip:ai-agent-1@ai.gateway.local",
                Speciality:  "sales_calls",
                IsAvailable: true,
            },
            {
                ID:          "agent-002",
                Name:        "Scam Call Handler",
                Endpoint:    "sip:ai-agent-2@ai.gateway.local",
                Speciality:  "scam_calls",
                IsAvailable: true,
            },
        },
    }
}

// RouteCall routes to an available AI agent
func (r *SimpleAIRouter) RouteCall(ctx context.Context, callID string, reason string) error {
    // Find available agent
    var selectedAgent *AIAgent
    for _, agent := range r.agents {
        if agent.IsAvailable {
            selectedAgent = &agent
            break
        }
    }
    
    if selectedAgent == nil {
        return fmt.Errorf("no available AI agents")
    }
    
    r.logger.Printf("Routing call %s to AI agent %s (%s)", callID, selectedAgent.ID, selectedAgent.Name)
    
    // In real implementation, this would initiate SIP transfer to AI agent
    // For now, just log the action
    
    return nil
}

// GetAvailableAgents returns list of available agents
func (r *SimpleAIRouter) GetAvailableAgents() []AIAgent {
    available := []AIAgent{}
    for _, agent := range r.agents {
        if agent.IsAvailable {
            available = append(available, agent)
        }
    }
    return available
}

// SimpleSIMManager provides basic SIM management
type SimpleSIMManager struct {
    db     *pgxpool.Pool
    logger *log.Logger
}

// NewSimpleSIMManager creates a basic SIM manager
func NewSimpleSIMManager(db *pgxpool.Pool, logger *log.Logger) *SimpleSIMManager {
    return &SimpleSIMManager{
        db:     db,
        logger: logger,
    }
}

// FlagSIM marks a SIM with an issue
func (s *SimpleSIMManager) FlagSIM(ctx context.Context, simID string, issue string) error {
    query := `
        UPDATE sim_cards 
        SET status = 'flagged',
            last_issue = $2,
            flagged_at = NOW(),
            updated_at = NOW()
        WHERE sim_id = $1
    `
    
    _, err := s.db.Exec(ctx, query, simID, issue)
    if err != nil {
        return fmt.Errorf("failed to update SIM status: %w", err)
    }
    
    s.logger.Printf("SIM %s flagged: %s", simID, issue)
    return nil
}

// GetSIMStatus retrieves current SIM status
func (s *SimpleSIMManager) GetSIMStatus(ctx context.Context, simID string) (*SIMStatus, error) {
    var status SIMStatus
    
    query := `
        SELECT sim_id, status, last_issue, flagged_at, action_taken
        FROM sim_cards
        WHERE sim_id = $1
    `
    
    err := s.db.QueryRow(ctx, query, simID).Scan(
        &status.SIMID,
        &status.Status,
        &status.Issue,
        &status.FlaggedAt,
        &status.ActionTaken,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &status, nil
}

// ScheduleReplacement schedules a SIM for replacement
func (s *SimpleSIMManager) ScheduleReplacement(ctx context.Context, simID string) error {
    query := `
        INSERT INTO sim_replacement_queue (
            sim_id, reason, priority, scheduled_at, status
        ) VALUES ($1, 'Blocked by operator', 'high', NOW(), 'pending')
        ON CONFLICT (sim_id) DO UPDATE 
        SET priority = 'high',
            updated_at = NOW()
    `
    
    _, err := s.db.Exec(ctx, query, simID)
    if err != nil {
        return fmt.Errorf("failed to schedule replacement: %w", err)
    }
    
    s.logger.Printf("SIM %s scheduled for replacement", simID)
    return nil
}