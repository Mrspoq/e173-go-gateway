# SIP Architecture Recommendation for E173 Gateway

## 🎯 Recommended Approach: Go-Based SIP Stack

### Why Go Instead of OpenSIPS?

**Advantages of Go SIP Implementation:**
1. **Unified Codebase**: Everything in Go, easier maintenance
2. **Custom Logic**: Direct integration of filtering rules
3. **Database Integration**: Native PostgreSQL connectivity
4. **Scalability**: Goroutines handle thousands of concurrent calls
5. **Maintainability**: Your team can modify and extend easily

### Recommended Go SIP Libraries:

1. **Primary Choice: `github.com/ghettovoice/gosip`**
   - Full SIP stack implementation
   - Support for UDP/TCP/TLS
   - Active development and community

2. **Alternative: `github.com/pion/webrtc` + Custom SIP**
   - More control over implementation
   - Better for custom routing logic

## 🏗️ Proposed SIP Architecture

```go
// SIP Server Component
type SIPServer struct {
    filterEngine  *FilterEngine
    routingEngine *RoutingEngine
    gatewayPool   *GatewayPool
    voiceAI       *VoiceAIService
}

// Advanced Filtering Pipeline
type FilterEngine struct {
    blacklistService    *BlacklistService
    whatsappValidator   *WhatsAppAPI
    phoneValidator      *GooglePhoneLib
    historyAnalyzer     *CallHistoryService
    operatorDetector    *OperatorPrefixService
}

// Intelligent Routing
type RoutingEngine struct {
    stickyRouting    map[string]string // number -> gateway mapping
    loadBalancer     *GatewayLoadBalancer
    failoverManager  *FailoverManager
}
```

## 🔄 Call Flow Architecture

```
Incoming SIP Call
        ↓
   Filter Engine
   ├── Blacklist Check
   ├── WhatsApp Validation  
   ├── Google Phone Validation
   ├── Call History Analysis
   └── Operator Detection
        ↓
   Routing Decision
   ├── Route to AI Agent (if spam)
   ├── Route to Specific Gateway
   └── Apply Sticky Routing
        ↓
   Gateway Selection
   ├── Load Balancing
   ├── Health Checking
   └── Failover Logic
        ↓
   Asterisk Gateway
   └── E173 Modem
```

## 🚀 Implementation Plan

### Phase 1: Basic SIP Server (Week 1)
```go
// Basic SIP server setup
func main() {
    sipServer := gosip.NewServer(
        gosip.WithTransport("udp", ":5060"),
        gosip.WithHandler(handleIncomingCall),
    )
    
    sipServer.Start()
}

func handleIncomingCall(req *gosip.Request) {
    // Extract caller/destination
    caller := extractCaller(req)
    destination := extractDestination(req)
    
    // Apply filters
    if shouldBlock(caller, destination) {
        routeToAI(req)
        return
    }
    
    // Route to gateway
    gateway := selectGateway(destination)
    forwardToGateway(req, gateway)
}
```

### Phase 2: Advanced Filtering (Week 2)
```go
type FilterResult struct {
    Allow       bool
    Reason      string
    RouteToAI   bool
    PreferredGW string
}

func (f *FilterEngine) ProcessCall(caller, dest string) FilterResult {
    // Multi-stage filtering
    checks := []FilterCheck{
        f.checkBlacklist,
        f.validateWhatsApp,
        f.validatePhoneNumber,
        f.analyzeHistory,
        f.detectOperator,
    }
    
    for _, check := range checks {
        if result := check(caller, dest); !result.Allow {
            return result
        }
    }
    
    return FilterResult{Allow: true}
}
```

## 🎨 Integration with Existing Code

### Database Extensions:
```sql
-- Add SIP-specific tables
CREATE TABLE sip_calls (
    id SERIAL PRIMARY KEY,
    call_id VARCHAR(255) UNIQUE,
    caller_number VARCHAR(50),
    destination_number VARCHAR(50),
    gateway_id INTEGER REFERENCES gateways(id),
    filter_result JSONB,
    routed_to_ai BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Extend gateways table
ALTER TABLE gateways ADD COLUMN sip_endpoint VARCHAR(255);
ALTER TABLE gateways ADD COLUMN health_status VARCHAR(50) DEFAULT 'healthy';
```

### Repository Extensions:
```go
type SIPCallRepository interface {
    CreateCall(*SIPCall) error
    GetCallHistory(number string, days int) ([]*SIPCall, error)
    GetGatewayStats(gatewayID int) (*GatewayStats, error)
    UpdateCallResult(callID string, result *CallResult) error
}
```

## 📊 Performance Considerations

### Concurrent Call Handling:
- **Target**: 1000+ concurrent calls per server
- **Implementation**: Goroutine per call + connection pooling
- **Monitoring**: Real-time call metrics

### Database Optimization:
- **Call History**: Partitioned by date
- **Indexing**: Caller number, destination, timestamp
- **Caching**: Redis for frequent lookups

### Failover & High Availability:
- **Health Checks**: Continuous gateway monitoring
- **Automatic Failover**: Sub-second switchover
- **Load Distribution**: Round-robin with weights

This Go-based approach gives you maximum control and integration while maintaining enterprise-grade performance.
