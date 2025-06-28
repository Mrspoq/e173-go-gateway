package ai

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "time"
)

// VoiceAgent represents an AI agent that handles spam calls
type VoiceAgent struct {
    ID           string       `json:"id"`
    Name         string       `json:"name"`
    Personality  string       `json:"personality"`
    Strategy     AgentStrategy `json:"strategy"`
    VoiceProfile VoiceProfile `json:"voice_profile"`
    IsActive     bool         `json:"is_active"`
}

// AgentStrategy defines how the agent handles calls
type AgentStrategy string

const (
    StrategyTimeWaster    AgentStrategy = "TIME_WASTER"     // Keep spammer on line as long as possible
    StrategyInfoCollector AgentStrategy = "INFO_COLLECTOR"  // Gather information about the scam
    StrategyConfuser      AgentStrategy = "CONFUSER"        // Confuse the spammer with nonsense
    StrategyEducator      AgentStrategy = "EDUCATOR"        // Educate about scams (for legitimate mistakes)
    StrategyRecorder      AgentStrategy = "RECORDER"        // Just record for evidence
)

// VoiceProfile defines the voice characteristics
type VoiceProfile struct {
    Gender      string  `json:"gender"`
    Age         string  `json:"age_range"`
    Accent      string  `json:"accent"`
    SpeechRate  float64 `json:"speech_rate"`
    Personality string  `json:"personality_traits"`
}

// VoiceAgentManager manages AI voice agents
type VoiceAgentManager struct {
    agents      []VoiceAgent
    ttsProvider TTSProvider
    llmProvider LLMProvider
    logger      *log.Logger
}

// TTSProvider handles text-to-speech conversion
type TTSProvider interface {
    Synthesize(ctx context.Context, text string, voice VoiceProfile) ([]byte, error)
    GetAvailableVoices() []VoiceProfile
}

// LLMProvider handles conversation generation
type LLMProvider interface {
    GenerateResponse(ctx context.Context, conversation []Message, agent VoiceAgent) (string, error)
    AnalyzeScam(ctx context.Context, transcript string) (*ScamAnalysis, error)
}

// Message represents a conversation message
type Message struct {
    Role      string    `json:"role"` // "agent" or "caller"
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}

// ScamAnalysis contains analysis of a scam call
type ScamAnalysis struct {
    ScamType     string   `json:"scam_type"`
    ThreatLevel  string   `json:"threat_level"`
    KeyPhrases   []string `json:"key_phrases"`
    TargetInfo   string   `json:"target_info"`
    Recommendations []string `json:"recommendations"`
}

// NewVoiceAgentManager creates a new voice agent manager
func NewVoiceAgentManager(tts TTSProvider, llm LLMProvider, logger *log.Logger) *VoiceAgentManager {
    return &VoiceAgentManager{
        ttsProvider: tts,
        llmProvider: llm,
        logger:      logger,
        agents: []VoiceAgent{
            {
                ID:          "agent-grandma",
                Name:        "Confused Grandma",
                Personality: "Sweet but confused elderly lady who can't hear well",
                Strategy:    StrategyTimeWaster,
                VoiceProfile: VoiceProfile{
                    Gender:      "female",
                    Age:         "70-80",
                    Accent:      "midwest_american",
                    SpeechRate:  0.8,
                    Personality: "confused, sweet, hard of hearing",
                },
                IsActive: true,
            },
            {
                ID:          "agent-techie",
                Name:        "Over-Enthusiastic Tech Support",
                Personality: "Extremely helpful tech person who overcomplicates everything",
                Strategy:    StrategyConfuser,
                VoiceProfile: VoiceProfile{
                    Gender:      "male",
                    Age:         "25-35",
                    Accent:      "silicon_valley",
                    SpeechRate:  1.2,
                    Personality: "enthusiastic, technical, verbose",
                },
                IsActive: true,
            },
            {
                ID:          "agent-investigator",
                Name:        "Curious Investigator",
                Personality: "Very interested person who asks lots of questions",
                Strategy:    StrategyInfoCollector,
                VoiceProfile: VoiceProfile{
                    Gender:      "female",
                    Age:         "30-40",
                    Accent:      "neutral",
                    SpeechRate:  1.0,
                    Personality: "curious, persistent, friendly",
                },
                IsActive: true,
            },
        },
    }
}

// HandleSpamCall assigns an agent to handle a spam call
func (m *VoiceAgentManager) HandleSpamCall(ctx context.Context, callID string, initialTranscript string) (*VoiceAgent, error) {
    // Select an appropriate agent
    agent := m.selectAgent(initialTranscript)
    if agent == nil {
        return nil, fmt.Errorf("no available agents")
    }
    
    m.logger.Printf("Assigning agent %s (%s) to handle call %s", agent.ID, agent.Name, callID)
    
    // Start the conversation
    go m.runConversation(ctx, callID, agent, initialTranscript)
    
    return agent, nil
}

// selectAgent chooses the best agent for the call
func (m *VoiceAgentManager) selectAgent(transcript string) *VoiceAgent {
    activeAgents := []VoiceAgent{}
    for _, agent := range m.agents {
        if agent.IsActive {
            activeAgents = append(activeAgents, agent)
        }
    }
    
    if len(activeAgents) == 0 {
        return nil
    }
    
    // For now, random selection. Could be enhanced with analysis
    return &activeAgents[rand.Intn(len(activeAgents))]
}

// runConversation manages the agent's conversation
func (m *VoiceAgentManager) runConversation(ctx context.Context, callID string, agent *VoiceAgent, initialTranscript string) {
    conversation := []Message{
        {
            Role:      "caller",
            Content:   initialTranscript,
            Timestamp: time.Now(),
        },
    }
    
    // Generate initial response
    response, err := m.generateAgentResponse(ctx, agent, conversation)
    if err != nil {
        m.logger.Printf("Failed to generate response: %v", err)
        return
    }
    
    // Convert to speech
    audio, err := m.ttsProvider.Synthesize(ctx, response, agent.VoiceProfile)
    if err != nil {
        m.logger.Printf("Failed to synthesize speech: %v", err)
        return
    }
    
    m.logger.Printf("Agent %s responding to call %s: %s", agent.Name, callID, response)
    
    // In real implementation, this would send audio to the call
    _ = audio
}

// generateAgentResponse creates a response based on agent personality
func (m *VoiceAgentManager) generateAgentResponse(ctx context.Context, agent *VoiceAgent, conversation []Message) (string, error) {
    // Build prompt based on agent strategy
    prompt := m.buildAgentPrompt(agent, conversation)
    
    // Generate response using LLM
    response, err := m.llmProvider.GenerateResponse(ctx, conversation, *agent)
    if err != nil {
        // Fallback responses based on agent type
        return m.getFallbackResponse(agent), nil
    }
    
    return response, nil
}

// buildAgentPrompt creates the LLM prompt for the agent
func (m *VoiceAgentManager) buildAgentPrompt(agent *VoiceAgent, conversation []Message) string {
    basePrompt := fmt.Sprintf(`You are %s. Your personality: %s

Your goal is to %s

Current conversation:
`, agent.Name, agent.Personality, m.getStrategyGoal(agent.Strategy))
    
    for _, msg := range conversation {
        basePrompt += fmt.Sprintf("\n%s: %s", msg.Role, msg.Content)
    }
    
    basePrompt += "\n\nGenerate your next response that fits your personality and advances your goal:"
    
    return basePrompt
}

// getStrategyGoal returns the goal description for a strategy
func (m *VoiceAgentManager) getStrategyGoal(strategy AgentStrategy) string {
    switch strategy {
    case StrategyTimeWaster:
        return "keep the caller on the line as long as possible by being confused, asking for repetition, and going off on tangents"
    case StrategyInfoCollector:
        return "gather as much information as possible about the scam by asking detailed questions while seeming interested"
    case StrategyConfuser:
        return "confuse the caller with technical jargon, contradictions, and nonsensical responses"
    case StrategyEducator:
        return "politely inform the caller about the harm of scam calls and encourage them to pursue legitimate work"
    default:
        return "engage with the caller appropriately"
    }
}

// getFallbackResponse provides a fallback when LLM fails
func (m *VoiceAgentManager) getFallbackResponse(agent *VoiceAgent) string {
    switch agent.ID {
    case "agent-grandma":
        responses := []string{
            "What's that dear? I can't hear you very well. Can you speak up?",
            "Oh my, that sounds complicated. Let me get my glasses...",
            "I'm sorry, I was just feeding my cats. What were you saying?",
            "Is this about my computer? My grandson usually helps me with that.",
        }
        return responses[rand.Intn(len(responses))]
        
    case "agent-techie":
        responses := []string{
            "Oh fascinating! But first, have you tried turning it off and on again? And by 'it' I mean your entire network infrastructure.",
            "Before we proceed, I need to explain the OSI model to you. It's crucial for understanding what we're about to do.",
            "Excellent question! Let me explain using a car analogy, but first, do you understand quantum computing?",
        }
        return responses[rand.Intn(len(responses))]
        
    case "agent-investigator":
        responses := []string{
            "That's very interesting! Can you tell me more about how this works exactly?",
            "I see. And what company did you say you were calling from again?",
            "Fascinating! How long have you been doing this type of work?",
        }
        return responses[rand.Intn(len(responses))]
        
    default:
        return "I'm sorry, could you repeat that?"
    }
}

// AnalyzeCompletedCall analyzes a completed scam call
func (m *VoiceAgentManager) AnalyzeCompletedCall(ctx context.Context, callID string, fullTranscript string) (*ScamAnalysis, error) {
    return m.llmProvider.AnalyzeScam(ctx, fullTranscript)
}

// GetAgentStats returns statistics for an agent
func (m *VoiceAgentManager) GetAgentStats(agentID string) map[string]interface{} {
    // This would query the database for real stats
    return map[string]interface{}{
        "agent_id":          agentID,
        "calls_handled":     42,
        "total_time_wasted": "3h 27m",
        "scams_documented":  15,
        "average_call_time": "4m 52s",
    }
}