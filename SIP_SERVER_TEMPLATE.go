package sip

import (
    "context"
    "fmt"
    "log"
    "net"
    "strings"
    "time"
    
    "github.com/ghettovoice/gosip/sip"
    "github.com/ghettovoice/gosip/transport"
)

// SIPServer handles incoming SIP calls and routes them intelligently
type SIPServer struct {
    server     sip.Server
    filterEng  *FilterEngine
    routingEng *RoutingEngine
    voiceAI    *VoiceAIService
    port       int
    logger     *log.Logger
}

// FilterEngine processes calls through multiple validation layers
type FilterEngine struct {
    blacklistSvc    *BlacklistService
    whatsappAPI     *WhatsAppValidator
    phoneValidator  *GooglePhoneValidator
    historyAnalyzer *CallHistoryAnalyzer
    operatorDetector *OperatorPrefixDetector
}

// RoutingEngine handles intelligent call routing
type RoutingEngine struct {
    gatewayPool    *GatewayPool
    stickyRoutes   map[string]string // caller -> gateway mapping
    loadBalancer   *LoadBalancer
    failoverMgr    *FailoverManager
}

// FilterResult contains the decision from filtering pipeline
type FilterResult struct {
    Allow       bool   `json:"allow"`
    Reason      string `json:"reason"`
    RouteToAI   bool   `json:"route_to_ai"`
    Gateway     string `json:"preferred_gateway"`
    Confidence  float64 `json:"confidence"`
}

// NewSIPServer creates a new SIP server instance
func NewSIPServer(port int) *SIPServer {
    return &SIPServer{
        port:       port,
        filterEng:  NewFilterEngine(),
        routingEng: NewRoutingEngine(),
        voiceAI:    NewVoiceAIService(),
        logger:     log.New(os.Stdout, "[SIP] ", log.LstdFlags),
    }
}

// Start begins listening for SIP calls
func (s *SIPServer) Start() error {
    logger := sip.NewDefaultLoggerFactory().
        WithLevel(sip.LogLevelError).
        CreateLogger()

    // Create server config
    serverConfig := &sip.ServerConfig{
        Host:      "0.0.0.0",
        Port:      s.port,
        Transport: []string{"udp", "tcp"},
    }

    // Create server
    server, err := sip.NewServer(serverConfig, nil, logger)
    if err != nil {
        return fmt.Errorf("failed to create SIP server: %w", err)
    }

    s.server = server

    // Register handlers
    s.server.OnInvite(s.handleInvite)
    s.server.OnBye(s.handleBye)
    s.server.OnCancel(s.handleCancel)
    s.server.OnAck(s.handleAck)

    s.logger.Printf("Starting SIP server on port %d", s.port)
    return s.server.Listen(context.Background())
}

// handleInvite processes incoming INVITE requests
func (s *SIPServer) handleInvite(req sip.Request, tx sip.ServerTransaction) {
    s.logger.Printf("Received INVITE: %s", req.String())

    // Extract call information
    from := req.From()
    to := req.To()
    callID := req.CallID()

    caller := extractPhoneNumber(from.Address.String())
    destination := extractPhoneNumber(to.Address.String())

    s.logger.Printf("Call: %s -> %s (Call-ID: %s)", caller, destination, callID)

    // Apply filtering pipeline
    filterResult := s.filterEng.ProcessCall(caller, destination)
    
    // Log filtering decision
    s.logger.Printf("Filter decision: %+v", filterResult)

    if !filterResult.Allow {
        if filterResult.RouteToAI {
            s.routeToAI(req, tx, filterResult)
        } else {
            s.rejectCall(req, tx, filterResult.Reason)
        }
        return
    }

    // Route to appropriate gateway
    gateway := s.routingEng.SelectGateway(destination, filterResult.Gateway)
    if gateway == nil {
        s.rejectCall(req, tx, "No available gateways")
        return
    }

    s.forwardToGateway(req, tx, gateway)
}

// ProcessCall applies all filtering rules
func (f *FilterEngine) ProcessCall(caller, destination string) FilterResult {
    ctx := context.Background()
    
    // Stage 1: Blacklist check (fastest)
    if blocked, reason := f.blacklistSvc.IsBlocked(ctx, caller); blocked {
        return FilterResult{
            Allow:  false,
            Reason: fmt.Sprintf("Blacklisted: %s", reason),
            RouteToAI: true, // Route spam to AI for monetization
        }
    }

    // Stage 2: Phone number validation
    if !f.phoneValidator.IsValid(destination) {
        return FilterResult{
            Allow:  false,
            Reason: "Invalid destination number",
        }
    }

    // Stage 3: WhatsApp validation (for real person check)
    if hasWhatsApp, confidence := f.whatsappAPI.HasWhatsApp(ctx, caller); !hasWhatsApp && confidence > 0.8 {
        return FilterResult{
            Allow:     false,
            Reason:    "Likely non-human caller",
            RouteToAI: true,
            Confidence: confidence,
        }
    }

    // Stage 4: Call history analysis
    history := f.historyAnalyzer.GetCallPattern(ctx, caller)
    if history.IsSpamPattern() {
        return FilterResult{
            Allow:     false,
            Reason:    "Spam calling pattern detected",
            RouteToAI: true,
            Confidence: history.SpamScore,
        }
    }

    // Stage 5: Operator detection and routing preference
    operator := f.operatorDetector.DetectOperator(destination)
    preferredGateway := f.operatorDetector.GetPreferredGateway(operator)

    return FilterResult{
        Allow:   true,
        Gateway: preferredGateway,
        Reason:  "Call approved",
    }
}

// routeToAI forwards spam calls to AI voice agents for monetization
func (s *SIPServer) routeToAI(req sip.Request, tx sip.ServerTransaction, result FilterResult) {
    s.logger.Printf("Routing call to AI agent: %s", result.Reason)
    
    // Get available AI agent
    aiAgent := s.voiceAI.GetAvailableAgent()
    if aiAgent == nil {
        s.rejectCall(req, tx, "No AI agents available")
        return
    }

    // Forward to AI service
    response := sip.NewResponseFromRequest("", req, 200, "OK", "")
    
    // Add AI agent contact
    contact := &sip.ContactHeader{
        Address: sip.Uri{
            Scheme: "sip",
            User:   "ai-agent",
            Host:   aiAgent.Endpoint,
        },
    }
    response.AppendHeader(contact.Name(), contact)
    
    tx.Respond(response)
    
    // Start billing for AI interaction
    s.voiceAI.StartBilling(req.CallID().Value(), result)
}

// forwardToGateway routes approved calls to Asterisk gateways
func (s *SIPServer) forwardToGateway(req sip.Request, tx sip.ServerTransaction, gateway *Gateway) {
    s.logger.Printf("Forwarding call to gateway: %s", gateway.Name)
    
    // Update sticky routing
    from := req.From()
    caller := extractPhoneNumber(from.Address.String())
    s.routingEng.UpdateStickyRoute(caller, gateway.ID)
    
    // Forward request to gateway
    gatewayURI := sip.Uri{
        Scheme: "sip",
        Host:   gateway.SIPEndpoint,
        Port:   gateway.SIPPort,
    }
    
    // Create new request for gateway
    forwardReq := req.Clone()
    forwardReq.SetDestination(gatewayURI.String())
    
    // Send to gateway
    response := sip.NewResponseFromRequest("", req, 100, "Trying", "")
    tx.Respond(response)
    
    // Handle gateway response (simplified)
    go s.handleGatewayResponse(req, tx, gateway)
}

// SelectGateway chooses the best gateway for a call
func (r *RoutingEngine) SelectGateway(destination, preferred string) *Gateway {
    // Check if we have a preferred gateway
    if preferred != "" {
        if gw := r.gatewayPool.GetByName(preferred); gw != nil && gw.IsHealthy() {
            return gw
        }
    }
    
    // Check sticky routing
    if sticky := r.stickyRoutes[destination]; sticky != "" {
        if gw := r.gatewayPool.GetByID(sticky); gw != nil && gw.IsHealthy() {
            return gw
        }
    }
    
    // Use load balancer
    return r.loadBalancer.GetNextGateway()
}

// Helper function to extract phone numbers from SIP URIs
func extractPhoneNumber(sipURI string) string {
    // Parse SIP URI and extract user part
    if idx := strings.Index(sipURI, "@"); idx > 0 {
        userPart := sipURI[:idx]
        if idx := strings.LastIndex(userPart, ":"); idx >= 0 {
            return userPart[idx+1:]
        }
        return userPart
    }
    return sipURI
}

// Additional handler methods...
func (s *SIPServer) handleBye(req sip.Request, tx sip.ServerTransaction) {
    s.logger.Printf("Call ended: %s", req.CallID())
    response := sip.NewResponseFromRequest("", req, 200, "OK", "")
    tx.Respond(response)
}

func (s *SIPServer) handleCancel(req sip.Request, tx sip.ServerTransaction) {
    s.logger.Printf("Call cancelled: %s", req.CallID())
    response := sip.NewResponseFromRequest("", req, 200, "OK", "")
    tx.Respond(response)
}

func (s *SIPServer) handleAck(req sip.Request, tx sip.ServerTransaction) {
    s.logger.Printf("Call established: %s", req.CallID())
}

// This template is ready to use - just need to implement the service interfaces!
