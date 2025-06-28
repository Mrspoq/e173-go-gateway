package sip

import (
    "context"
    "fmt"
    "log"
    "net"
    "time"
    
    "github.com/e173-gateway/e173_go_gateway/pkg/voice"
    "github.com/e173-gateway/e173_go_gateway/pkg/ai"
    "github.com/jackc/pgx/v4/pgxpool"
)

// VoiceEnabledSIPServer extends BasicSIPServer with voice recognition
type VoiceEnabledSIPServer struct {
    *BasicSIPServer
    voiceService   *voice.RecognitionService
    agentManager   *ai.VoiceAgentManager
    audioCapture   *voice.AudioCapture
}

// NewVoiceEnabledSIPServer creates a SIP server with voice recognition
func NewVoiceEnabledSIPServer(port int, whatsappAPIKey string, dbPool *pgxpool.Pool, 
    sttProvider voice.STTProvider, llmProvider voice.LLMProvider) *VoiceEnabledSIPServer {
    
    // Create base SIP server
    baseSIPServer := NewBasicSIPServerWithDB(port, whatsappAPIKey, dbPool)
    
    // Create voice components
    classifier := voice.NewRuleBasedClassifier() // Start with rule-based, can upgrade to LLM
    audioRecorder := voice.NewSimpleAudioRecorder()
    audioCapture := voice.NewAudioCapture()
    
    // Create action engine with database
    aiRouter := voice.NewSimpleAIRouter(baseSIPServer.logger)
    simManager := voice.NewSimpleSIMManager(dbPool, baseSIPServer.logger)
    actionEngine := voice.NewDefaultActionEngine(dbPool, aiRouter, simManager, baseSIPServer.logger)
    
    // Create voice recognition service
    voiceService := voice.NewRecognitionService(sttProvider, classifier, actionEngine, audioRecorder)
    
    // Create AI agent manager (would need proper TTS/LLM providers in production)
    agentManager := ai.NewVoiceAgentManager(nil, nil, baseSIPServer.logger)
    
    return &VoiceEnabledSIPServer{
        BasicSIPServer: baseSIPServer,
        voiceService:   voiceService,
        agentManager:   agentManager,
        audioCapture:   audioCapture,
    }
}

// handleInviteWithVoice processes INVITE with voice recognition
func (s *VoiceEnabledSIPServer) handleInviteWithVoice(message string, clientAddr *net.UDPAddr) {
    // First, do standard processing
    caller := extractSIPHeader(message, "From:")
    destination := extractSIPHeader(message, "To:")
    callID := extractSIPHeader(message, "Call-ID:")
    
    callerNumber := extractPhoneNumber(caller)
    destNumber := extractPhoneNumber(destination)
    
    s.logger.Printf("Processing INVITE with voice recognition: %s -> %s (Call-ID: %s)", 
        callerNumber, destNumber, callID)
    
    // Apply standard filtering
    filterResult := s.filterEng.ProcessCall(callerNumber, destNumber)
    
    // Start audio capture for this call
    audioFormat := voice.AudioFormat{
        SampleRate:    8000,
        Channels:      1,
        BitsPerSample: 16,
        Codec:         "ulaw", // Standard for telephony
    }
    
    audioStream, err := s.audioCapture.StartCapture(callID, voice.DirectionIncoming, audioFormat)
    if err != nil {
        s.logger.Printf("Failed to start audio capture: %v", err)
    }
    
    // If call is allowed but suspicious, analyze voice
    if filterResult.Allow && filterResult.Confidence < 0.8 {
        go s.analyzeCallVoice(callID, audioStream)
    }
    
    // If already detected as spam, route to AI immediately
    if filterResult.RouteToAI {
        s.routeToAIWithVoice(message, clientAddr, filterResult, callID)
        return
    }
    
    // Continue with standard processing
    if !filterResult.Allow {
        s.rejectCall(message, clientAddr, filterResult.Reason)
        return
    }
    
    gateway := s.routingEng.SelectGateway(destNumber, filterResult.Gateway)
    if gateway == nil {
        s.rejectCall(message, clientAddr, "No available gateways")
        return
    }
    
    s.forwardToGateway(message, clientAddr, gateway)
}

// analyzeCallVoice performs real-time voice analysis
func (s *VoiceEnabledSIPServer) analyzeCallVoice(callID string, audioStream *voice.AudioStream) {
    ctx := context.Background()
    
    // Wait for enough audio data (e.g., 3 seconds)
    time.Sleep(3 * time.Second)
    
    // Analyze the voice
    result, err := s.voiceService.AnalyzeIncomingCall(ctx, callID, audioStream)
    if err != nil {
        s.logger.Printf("Voice analysis failed for call %s: %v", callID, err)
        return
    }
    
    s.logger.Printf("Voice analysis for call %s: Category=%s, Action=%s, Confidence=%.2f",
        callID, result.Classification.Category, result.Classification.Action, 
        result.Classification.Confidence)
    
    // Update call record with voice analysis
    s.updateCallWithVoiceAnalysis(ctx, callID, result)
    
    // If spam detected after call started, we can still take action
    if result.Classification.Category == voice.CategorySpamRobocall {
        s.logger.Printf("Late spam detection for call %s - routing to AI", callID)
        // In production, this would trigger a call transfer to AI
    }
}

// routeToAIWithVoice routes spam calls to AI agents with voice handling
func (s *VoiceEnabledSIPServer) routeToAIWithVoice(message string, clientAddr *net.UDPAddr, 
    result FilterResult, callID string) {
    
    s.logger.Printf("Routing call %s to AI agent with voice handling", callID)
    
    // Send 200 OK to accept the call
    response := buildSIPResponse("200 OK", message)
    s.sendSIPResponse(response, clientAddr)
    
    // Get initial transcript if available
    initialTranscript := "Automated call detected"
    
    // Assign an AI agent
    ctx := context.Background()
    agent, err := s.agentManager.HandleSpamCall(ctx, callID, initialTranscript)
    if err != nil {
        s.logger.Printf("Failed to assign AI agent: %v", err)
        return
    }
    
    s.logger.Printf("Call %s assigned to AI agent: %s", callID, agent.Name)
    
    // Update database
    s.updateCallWithAIAgent(ctx, callID, agent.ID)
}

// updateCallWithVoiceAnalysis updates the database with voice analysis results
func (s *VoiceEnabledSIPServer) updateCallWithVoiceAnalysis(ctx context.Context, 
    callID string, result *voice.RecognitionResult) error {
    
    // Log the voice analysis result
    s.logger.Printf("Voice analysis completed for call %s: Category=%s, Action=%s",
        callID, result.Classification.Category, result.Classification.Action)
    
    // In production, this would update the database
    return nil
}

// updateCallWithAIAgent updates the call record with AI agent assignment
func (s *VoiceEnabledSIPServer) updateCallWithAIAgent(ctx context.Context, 
    callID string, agentID string) error {
    
    // Since BasicSIPServer doesn't have db field, we need to get it from somewhere
    // For now, we'll skip this implementation
    s.logger.Printf("Would update call %s with AI agent %s in database", callID, agentID)
    
    return nil
}

// MonitorOutgoingSIM monitors outgoing calls for SIM status messages
func (s *VoiceEnabledSIPServer) MonitorOutgoingSIM(ctx context.Context, 
    callID string, simID string, audioStream *voice.AudioStream) {
    
    // Analyze outgoing call for operator messages
    result, err := s.voiceService.AnalyzeOutgoingCall(ctx, callID, simID, audioStream)
    if err != nil {
        s.logger.Printf("Failed to analyze outgoing call %s: %v", callID, err)
        return
    }
    
    // Log any SIM issues detected
    if result.Classification.Category == voice.CategorySIMBlocked ||
       result.Classification.Category == voice.CategoryLowCredit {
        
        s.logger.Printf("SIM issue detected for %s: %s", simID, result.Classification.Reason)
        
        // Create alert
        s.createSIMAlert(ctx, simID, result.Classification)
    }
}

// createSIMAlert creates an alert for SIM issues
func (s *VoiceEnabledSIPServer) createSIMAlert(ctx context.Context, 
    simID string, classification *voice.Classification) error {
    
    severity := "medium"
    if classification.Category == voice.CategorySIMBlocked {
        severity = "high"
    }
    
    s.logger.Printf("SIM Alert: SIM=%s, Type=%s, Severity=%s, Reason=%s",
        simID, classification.Category, severity, classification.Reason)
    
    // In production, this would create a database alert
    return nil
}