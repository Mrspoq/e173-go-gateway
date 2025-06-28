package sip

import (
    "fmt"
    "log"
    "net"
    "strings"
    
    "github.com/e173-gateway/e173_go_gateway/pkg/validation"
)

// BasicSIPServer handles SIP calls with custom routing logic
type BasicSIPServer struct {
    listener   net.PacketConn
    port       int
    filterEng  *FilterEngine
    routingEng *RoutingEngine
    voiceAI    *VoiceAIService
    logger     *log.Logger
}

// FilterEngine processes calls through multiple validation layers
type FilterEngine struct {
    blacklistSvc    *BlacklistService
    whatsappAPI     *validation.PrivateWhatsAppValidator
    phoneValidator  *validation.GooglePhoneValidator
    historyAnalyzer *CallHistoryAnalyzer
    operatorDetector *OperatorPrefixDetector
}

// RoutingEngine handles intelligent call routing
type RoutingEngine struct {
    gatewayPool    *GatewayPool
    stickyRoutes   map[string]string
    loadBalancer   *LoadBalancer
    failoverMgr    *FailoverManager
}

// FilterResult contains the decision from filtering pipeline
type FilterResult struct {
    Allow       bool    `json:"allow"`
    Reason      string  `json:"reason"`
    RouteToAI   bool    `json:"route_to_ai"`
    Gateway     string  `json:"preferred_gateway"`
    Confidence  float64 `json:"confidence"`
}

// Gateway represents a remote Asterisk gateway
type Gateway struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    SIPEndpoint string `json:"sip_endpoint"`
    SIPPort     int    `json:"sip_port"`
    IsHealthy   func() bool
}



// Services with real implementations
type BlacklistService struct {
    // TODO: Implement database-backed blacklist
}

type CallHistoryAnalyzer struct {
    // TODO: Implement database-backed call pattern analysis
}

type OperatorPrefixDetector struct {
    // TODO: Implement database-backed operator detection
}

type GatewayPool struct {
    // TODO: Implement database-backed gateway management
}

type LoadBalancer struct {
    // TODO: Implement intelligent load balancing
}

type FailoverManager struct {
    // TODO: Implement failover logic
}

type VoiceAIService struct {
    // TODO: Implement AI voice agent integration
}

// NewBasicSIPServer creates a new SIP server instance
func NewBasicSIPServer(port int, whatsappAPIKey string) *BasicSIPServer {
    return &BasicSIPServer{
        port:       port,
        filterEng:  NewFilterEngine(whatsappAPIKey),
        routingEng: NewRoutingEngine(),
        voiceAI:    NewVoiceAIService(),
        logger:     log.New(log.Writer(), "[SIP] ", log.LstdFlags),
    }
}

// Start begins listening for SIP packets
func (s *BasicSIPServer) Start() error {
    addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.port))
    if err != nil {
        return fmt.Errorf("failed to resolve UDP address: %w", err)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        return fmt.Errorf("failed to listen on UDP: %w", err)
    }

    s.listener = conn
    s.logger.Printf("Starting SIP server on port %d", s.port)

    // Start handling packets
    buffer := make([]byte, 4096)
    for {
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            s.logger.Printf("Error reading UDP packet: %v", err)
            continue
        }

        // Process SIP message in goroutine
        go s.handleSIPMessage(buffer[:n], clientAddr)
    }
}

// handleSIPMessage processes incoming SIP messages
func (s *BasicSIPServer) handleSIPMessage(data []byte, clientAddr *net.UDPAddr) {
    message := string(data)
    s.logger.Printf("Received SIP message from %s:\n%s", clientAddr, message)

    // Parse basic SIP message
    if strings.HasPrefix(message, "INVITE") {
        s.handleInvite(message, clientAddr)
    } else if strings.HasPrefix(message, "BYE") {
        s.handleBye(message, clientAddr)
    } else if strings.HasPrefix(message, "CANCEL") {
        s.handleCancel(message, clientAddr)
    } else {
        s.logger.Printf("Unhandled SIP method: %s", strings.Split(message, " ")[0])
    }
}

// handleInvite processes INVITE requests
func (s *BasicSIPServer) handleInvite(message string, clientAddr *net.UDPAddr) {
    // Extract call information from SIP headers
    caller := extractSIPHeader(message, "From:")
    destination := extractSIPHeader(message, "To:")
    callID := extractSIPHeader(message, "Call-ID:")

    callerNumber := extractPhoneNumber(caller)
    destNumber := extractPhoneNumber(destination)

    s.logger.Printf("Processing INVITE: %s -> %s (Call-ID: %s)", callerNumber, destNumber, callID)

    // Apply filtering pipeline
    filterResult := s.filterEng.ProcessCall(callerNumber, destNumber)
    s.logger.Printf("Filter decision: %+v", filterResult)

    if !filterResult.Allow {
        if filterResult.RouteToAI {
            s.routeToAI(message, clientAddr, filterResult)
        } else {
            s.rejectCall(message, clientAddr, filterResult.Reason)
        }
        return
    }

    // Route to appropriate gateway
    gateway := s.routingEng.SelectGateway(destNumber, filterResult.Gateway)
    if gateway == nil {
        s.rejectCall(message, clientAddr, "No available gateways")
        return
    }

    s.forwardToGateway(message, clientAddr, gateway)
}

// ProcessCall applies all filtering rules
func (f *FilterEngine) ProcessCall(caller, destination string) FilterResult {
    // Stage 1: Basic validation
    if caller == "" || destination == "" {
        return FilterResult{
            Allow:  false,
            Reason: "Invalid caller or destination",
        }
    }

    // Stage 2: Phone number validation using Google Phone Lib
    if !f.phoneValidator.IsValidMobile(caller) {
        return FilterResult{
            Allow:  false,
            Reason: "Invalid caller phone number format",
        }
    }

    if !f.phoneValidator.IsValidMobile(destination) {
        return FilterResult{
            Allow:  false,
            Reason: "Invalid destination phone number format",
        }
    }

    // Stage 3: WhatsApp validation for real person check
    isRealPerson, confidence, err := f.whatsappAPI.IsLikelyRealPerson(caller)
    if err == nil && !isRealPerson && confidence > 0.8 {
        return FilterResult{
            Allow:     false,
            Reason:    "Caller likely not a real person",
            RouteToAI: true,
            Confidence: confidence,
        }
    }

    // Stage 4: Blacklist check (simulate for now)
    if strings.Contains(caller, "spam") {
        return FilterResult{
            Allow:     false,
            Reason:    "Caller is blacklisted",
            RouteToAI: true,
            Confidence: 0.95,
        }
    }

    // Stage 5: Operator detection and routing
    operator := f.detectOperator(destination)
    
    return FilterResult{
        Allow:   true,
        Gateway: operator,
        Reason:  "Call approved",
    }
}

// detectOperator determines the operator from the phone number
func (f *FilterEngine) detectOperator(number string) string {
    // Simple operator detection based on prefixes
    if strings.HasPrefix(number, "+234803") || strings.HasPrefix(number, "+234806") {
        return "MTN"
    } else if strings.HasPrefix(number, "+234802") {
        return "Airtel"
    } else if strings.HasPrefix(number, "+234805") {
        return "Glo"
    }
    return "Unknown"
}

// SelectGateway chooses the best gateway for a call
func (r *RoutingEngine) SelectGateway(destination, preferred string) *Gateway {
    // Simulate gateway selection
    return &Gateway{
        ID:          "gw-001",
        Name:        "Gateway 1",
        SIPEndpoint: "192.168.1.100",
        SIPPort:     5060,
        IsHealthy:   func() bool { return true },
    }
}

// routeToAI forwards spam calls to AI voice agents
func (s *BasicSIPServer) routeToAI(message string, clientAddr *net.UDPAddr, result FilterResult) {
    s.logger.Printf("Routing call to AI agent: %s", result.Reason)
    
    // Send 200 OK response
    response := buildSIPResponse("200 OK", message)
    s.sendSIPResponse(response, clientAddr)
    
    // TODO: Integrate with actual AI voice service
}

// forwardToGateway routes approved calls to Asterisk gateways
func (s *BasicSIPServer) forwardToGateway(message string, clientAddr *net.UDPAddr, gateway *Gateway) {
    s.logger.Printf("Forwarding call to gateway: %s", gateway.Name)
    
    // Send 100 Trying response
    response := buildSIPResponse("100 Trying", message)
    s.sendSIPResponse(response, clientAddr)
    
    // TODO: Forward to actual gateway
}

// rejectCall sends rejection response
func (s *BasicSIPServer) rejectCall(message string, clientAddr *net.UDPAddr, reason string) {
    s.logger.Printf("Rejecting call: %s", reason)
    
    response := buildSIPResponse("403 Forbidden", message)
    s.sendSIPResponse(response, clientAddr)
}

// sendSIPResponse sends a SIP response back to the client
func (s *BasicSIPServer) sendSIPResponse(response string, clientAddr *net.UDPAddr) {
    conn := s.listener.(*net.UDPConn)
    _, err := conn.WriteToUDP([]byte(response), clientAddr)
    if err != nil {
        s.logger.Printf("Error sending SIP response: %v", err)
    }
}

// Helper functions

func extractSIPHeader(message, header string) string {
    lines := strings.Split(message, "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, header) {
            return strings.TrimSpace(line[len(header):])
        }
    }
    return ""
}

func extractPhoneNumber(sipHeader string) string {
    // Extract phone number from SIP header (simplified)
    if idx := strings.Index(sipHeader, "sip:"); idx >= 0 {
        start := idx + 4
        if end := strings.Index(sipHeader[start:], "@"); end >= 0 {
            return sipHeader[start : start+end]
        }
    }
    return sipHeader
}

func buildSIPResponse(status, originalMessage string) string {
    callID := extractSIPHeader(originalMessage, "Call-ID:")
    from := extractSIPHeader(originalMessage, "From:")
    to := extractSIPHeader(originalMessage, "To:")
    via := extractSIPHeader(originalMessage, "Via:")
    
    return fmt.Sprintf(`SIP/2.0 %s
Via: %s
From: %s
To: %s
Call-ID: %s
Content-Length: 0

`, status, via, from, to, callID)
}

// NewFilterEngine creates filter engine with real validation services
func NewFilterEngine(whatsappAPIKey string) *FilterEngine {
    return &FilterEngine{
        blacklistSvc:     &BlacklistService{},
        whatsappAPI:      validation.NewPrivateWhatsAppValidator(whatsappAPIKey), // Use your private API
        phoneValidator:   validation.NewGooglePhoneValidator("NG"), // Default to Nigeria
        historyAnalyzer:  &CallHistoryAnalyzer{},
        operatorDetector: &OperatorPrefixDetector{},
    }
}

func NewRoutingEngine() *RoutingEngine {
    return &RoutingEngine{
        gatewayPool:   &GatewayPool{},
        stickyRoutes:  make(map[string]string),
        loadBalancer:  &LoadBalancer{},
        failoverMgr:   &FailoverManager{},
    }
}

func NewVoiceAIService() *VoiceAIService {
    return &VoiceAIService{}
}

func (s *BasicSIPServer) handleBye(message string, clientAddr *net.UDPAddr) {
    s.logger.Printf("Call ended from %s", clientAddr)
    response := buildSIPResponse("200 OK", message)
    s.sendSIPResponse(response, clientAddr)
}

func (s *BasicSIPServer) handleCancel(message string, clientAddr *net.UDPAddr) {
    s.logger.Printf("Call cancelled from %s", clientAddr)
    response := buildSIPResponse("200 OK", message)
    s.sendSIPResponse(response, clientAddr)
}
