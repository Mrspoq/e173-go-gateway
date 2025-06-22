package ami

import (
	"context"
	"encoding/json"
	"fmt"
	"net" // Added import
	"regexp" // Added import
	"strconv"
	"strings"
	"sync"
	"time"

	cfg "github.com/e173-gateway/e173_go_gateway/pkg/config"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/e173-gateway/e173_go_gateway/pkg/repository"
	"github.com/sirupsen/logrus"
	goami2 "github.com/staskobzar/goami2"
)

const (
	reconnectDelay = 5 * time.Second
)

var (
	// Example: Dongle/dongle0-0100000000 or DAHDI/i1/12345-1
	// This regex attempts to capture the device name (e.g., "dongle0", "i1")
	channelDeviceRegex = regexp.MustCompile(`^(?:Dongle|DAHDI)/([^/-]+)`)
)

// Helper function to parse string to *int, returning nil if empty or error
func parseOptionalInt(valStr string) *int {
	if valStr == "" {
		return nil
	}
	if val, err := strconv.Atoi(valStr); err == nil {
		return &val
	}
	return nil
}

// Helper function to parse string to *int64, returning nil if empty or error
func parseOptionalInt64(valStr string) *int64 {
	if valStr == "" {
		return nil
	}
	if val, err := strconv.ParseInt(valStr, 10, 64); err == nil {
		return &val
	}
	return nil
}

// AMIService handles the connection to Asterisk Manager Interface
// and processes events to generate Call Detail Records (CDRs).
type AMIService struct {
	config         *cfg.AppConfig
	amiClient      *goami2.Client // Changed AMIClient to Client
	cdrRepo        repository.CdrRepository
	logger         *logrus.Logger
	mu             sync.Mutex
	connected      bool
	lastEventTime  time.Time
	reconnectMutex sync.Mutex
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewAMIService creates a new AMIService.
func NewAMIService(appConfig *cfg.AppConfig, cdrRepo repository.CdrRepository, logger *logrus.Logger) (*AMIService, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &AMIService{
		config:  appConfig,
		cdrRepo: cdrRepo,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

// Start initiates the AMI connection and event processing loop.
func (s *AMIService) Start() {
	s.logger.Info("Starting AMIService...")
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				s.logger.Info("AMIService shutting down.")
				if s.amiClient != nil {
					s.amiClient.Close()
				}
				return
			default:
				err := s.connectAndListen()
				if err != nil {
					s.logger.Errorf("AMI connection or listener error: %v. Reconnecting in %s...", err, reconnectDelay)
					if s.amiClient != nil {
						s.amiClient.Close()
						s.amiClient = nil
					}
					time.Sleep(reconnectDelay)
				} else {
					s.logger.Info("connectAndListen returned without error (likely due to context cancellation), AMIService stopping.")
					return
				}
			}
		}
	}()
}

// Stop gracefully shuts down the AMIService.
func (s *AMIService) Stop() {
	s.logger.Info("Stopping AMIService...")
	s.cancel()
}

func (s *AMIService) connectAndListen() error {
	amiAddress := fmt.Sprintf("%s:%s", s.config.AsteriskAMIHost, s.config.AsteriskAMIPort)
	s.logger.Infof("Attempting to connect to AMI at %s", amiAddress)

	conn, err := net.DialTimeout("tcp", amiAddress, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to dial AMI: %w", err)
	}

	// Pass s.ctx to NewClientWithContext for graceful shutdown propagation
	client, err := goami2.NewClientWithContext(s.ctx, conn, s.config.AsteriskAMIUser, s.config.AsteriskAMIPass)
	if err != nil {
		conn.Close() // Ensure connection is closed on client creation failure
		return fmt.Errorf("failed to create AMI client or login: %w", err)
	}
	s.amiClient = client
	s.logger.Info("Successfully connected and logged in to AMI.")

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Context cancelled during listen, closing AMI client.")
			s.amiClient.Close() // Close client on context cancellation
			return s.ctx.Err()
		case msg, ok := <-s.amiClient.AllMessages():
			if !ok {
				s.logger.Warn("AMI AllMessages channel closed. Connection likely lost.")
				return fmt.Errorf("AllMessages channel closed")
			}
			s.handleEvent(msg)
		case err, ok := <-s.amiClient.Err():
			if !ok {
				s.logger.Warn("AMI Err channel closed.")
				// This might happen during graceful shutdown if the client closes Err before AllMessages
				return fmt.Errorf("Err channel closed")
			}
			// If context is done, this error might be a result of closing the connection.
			if s.ctx.Err() != nil {
				s.logger.Infof("AMI client error received after context cancellation: %v", err)
				return s.ctx.Err()
			}
			if err != nil {
				s.logger.Errorf("AMI client run error: %v", err)
				return fmt.Errorf("AMI client run error: %w", err) // Propagate error to trigger reconnect
			}
		}
	}
}

func (s *AMIService) handleEvent(msg *goami2.Message) {
	eventName := getHeader(msg, "Event")
	// Log relevant events, but can be noisy. Consider DEBUG level for all.
	if eventName == "Hangup" || eventName == "Newchannel" || eventName == "Cdr" {
		s.logger.Infof("Received AMI Event: %s, UniqueID: %s, LinkedID: %s, Channel: %s",
			eventName, getHeader(msg, "Uniqueid"), getHeader(msg, "Linkedid"), getHeader(msg, "Channel"))
		s.logger.Debugf("Full Event Data for %s (%s): %v", eventName, getHeader(msg, "Uniqueid"), getAllHeadersAsMap(msg))
	}

	switch eventName {
	case "Hangup":
		s.processHangupEvent(msg)
	// Future:
	// case "Newchannel":
	// s.processNewChannelEvent(msg) // To capture CallStartTime accurately
	// case "Cdr":
	// s.processCdrEvent(msg) // If Asterisk is configured to send full Cdr events
	default:
		// s.logger.Debugf("Unhandled AMI Event: %s, Data: %v", eventName, msg.Headers)
	}
}

func (s *AMIService) processHangupEvent(msg *goami2.Message) {
	uniqueID := getHeader(msg, "Uniqueid") 
	if uniqueID == "" {
		s.logger.Infof("Hangup event received: %+v", getAllHeadersAsMap(msg)) 
		return
	}
	s.logger.Infof("Processing Hangup event for UniqueID: %s", uniqueID)

	// --- Critical Identifiers (often need dialplan setup) ---
	var modemIDForCdr, simCardIDForCdr *int

	// ModemID: Can be inferred from channel or set as a variable e.g., CDR(modem_id) or DONGLE_NAME
	rawModemIDStr := getHeader(msg, "CDR(modem_id)") 
	if rawModemIDStr == "" {
		rawModemIDStr = getHeader(msg, "DONGLE_NAME") 
	}
	if rawModemIDStr == "" { // Fallback: try to parse from channel string
		channel := getHeader(msg, "Channel") 
		matches := channelDeviceRegex.FindStringSubmatch(channel)
		if len(matches) > 1 {
			rawModemIDStr = matches[1] // Get the string value
			s.logger.Debugf("CDR for %s: Inferred ModemID string '%s' from Channel '%s'", uniqueID, rawModemIDStr, channel)
		}
	}
	modemIDForCdr = parseOptionalInt(rawModemIDStr)
	if modemIDForCdr == nil && rawModemIDStr != "" { 
		s.logger.Warnf("CDR for %s: Failed to parse ModemID string '%s' as int64.", uniqueID, rawModemIDStr)
	} else if modemIDForCdr == nil {
		s.logger.Warnf("CDR for %s: ModemID is missing or unparsable. Ensure 'CDR(modem_id)', 'DONGLE_NAME', or channel name part is a valid int64 string.", uniqueID)
	}

	// SIMCardID: Best sourced from a channel variable set in dialplan e.g., CDR(sim_iccid) or SIM_ICCID
	rawSimCardIDStr := getHeader(msg, "CDR(sim_iccid)") 
	if rawSimCardIDStr == "" {
		rawSimCardIDStr = getHeader(msg, "SIM_ICCID") 
	}
	simCardIDForCdr = parseOptionalInt(rawSimCardIDStr)
	if simCardIDForCdr == nil && rawSimCardIDStr != "" { 
		s.logger.Warnf("CDR for %s: Failed to parse SimCardID string '%s' as int64.", uniqueID, rawSimCardIDStr)
	} else if simCardIDForCdr == nil {
		s.logger.Warnf("CDR for %s: SIMCardID is missing or unparsable. Ensure 'CDR(sim_iccid)' or 'SIM_ICCID' is a valid int64 string.", uniqueID)
	}

	// --- Call Timings ---
	callStartTimeStr := getHeader(msg, "CDR(start)") 
	callStartTime := parseAMITime(callStartTimeStr, s.logger)
	if callStartTime.IsZero() {
		s.logger.Warnf("CDR for %s: CallStartTime is MISSING or invalid from CDR(start). This CDR will be inaccurate. Event Timestamp (Hangup Time): %s", uniqueID, getHeader(msg, "Timestamp")) 
	}
	callEndTime := parseAMITime(getHeader(msg, "Timestamp"), s.logger) 

	// --- Durations ---
	var durationSeconds, billableDurationSeconds *int
	if priorityStr := getHeader(msg, "Priority"); priorityStr != "" { 
		if val, err := strconv.ParseInt(priorityStr, 10, 0); err == nil { 
			parsedVal := int(val)
			durationSeconds = &parsedVal
		} else {
			s.logger.Warnf("CDR for %s: Failed to parse Priority (duration) '%s': %v", uniqueID, priorityStr, err)
		}
	}
	if billableSecondsStr := getHeader(msg, "BillableSeconds"); billableSecondsStr != "" { 
		if val, err := strconv.ParseInt(billableSecondsStr, 10, 0); err == nil { 
			parsedVal := int(val)
			billableDurationSeconds = &parsedVal
		} else {
			s.logger.Warnf("CDR for %s: Failed to parse BillableSeconds '%s': %v", uniqueID, billableSecondsStr, err)
		}
	}

	// --- Disposition ---
	disposition := determineDisposition(msg) 

	// Marshal all event fields to JSON for RawEventData
	var rawEventDataJSON []byte
	rawEventFields := getAllHeadersAsMap(msg)
	if len(rawEventFields) > 0 { 
		var err error
		rawEventDataJSON, err = json.Marshal(rawEventFields)
		if err != nil {
			s.logger.Warnf("CDR for %s: Failed to marshal raw event data to JSON: %v", uniqueID, err)
		}
	} // Corrected brace placement

	cdr := &models.Cdr{
		UniqueID:            uniqueID,
		Channel:             models.StringPtr(getHeader(msg, "Channel")), 
		CallerIDNum:         models.StringPtr(getHeader(msg, "CallerIDNum")), 
		CallerIDName:        models.StringPtr(getHeader(msg, "CallerIDName")), 
		ConnectedLineNum:    models.StringPtr(getHeader(msg, "ConnectedLineNum")), 
		ConnectedLineName:   models.StringPtr(getHeader(msg, "ConnectedLineName")), 
		AccountCode:         models.StringPtr(getHeader(msg, "AccountCode")), 
		Cause:               models.StringPtr(getHeader(msg, "Cause")), 
		CauseTxt:            models.StringPtr(getHeader(msg, "CauseTxt")), 
		Disposition:         getOptionalString(disposition),
		StartTime:           &callStartTime,
		AnswerTime:          nil,
		EndTime:             &callEndTime,
		Duration:            durationSeconds,
		BillableSeconds:     billableDurationSeconds,
		ModemID:             modemIDForCdr,
		SimCardID:           simCardIDForCdr,
		CallDirection:       getOptionalString(determineCallDirection(msg, s.logger)), 
		IsSpam:              new(bool),
		Context:             models.StringPtr(getHeader(msg, "Context")), 
		Extension:           models.StringPtr(getHeader(msg, "Exten")), 
		Priority:            parseOptionalInt(getHeader(msg, "Priority")), 
		RawEventData:        rawEventDataJSON,
	}

	// Log the populated CDR before saving
	s.logger.Debugf("Populated CDR for %s: %+v", uniqueID, cdr)

	err := s.cdrRepo.CreateCdr(s.ctx, cdr)
	if err != nil {
		s.logger.Errorf("Failed to create CDR for UniqueID %s: %v", uniqueID, err)
		return
	}
	s.logger.Infof("Successfully created CDR for UniqueID %s (DB ID: %s)", uniqueID, cdr.ID)
}

func getOptionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func parseAMITime(timestampStr string, logger *logrus.Logger) time.Time { // Changed logger type
	if timestampStr == "" {
		return time.Time{}
	}
	// Try parsing as Unix epoch (integer or float)
	if sec, err := strconv.ParseFloat(timestampStr, 64); err == nil {
		return time.Unix(int64(sec), int64((sec-float64(int64(sec)))*1e9)).UTC()
	}
	// Add other parsing logic if Asterisk sends timestamps in different formats
	logger.Warnf("Could not parse AMI timestamp: '%s' as Unix epoch. Returning zero time.", timestampStr)
	return time.Time{}
}

func determineCallDirection(msg *goami2.Message, logger *logrus.Logger) string { // Changed logger type
	// This is highly dependent on your dialplan.
	// Examples:
	// - Check if CDR(direction) variable is set.
	// - Check channel name patterns (e.g., if Channel starts with "SIP/your_trunk_provider" it's likely Inbound).
	// - Check dialplan context (e.g., Context "from-pstn" vs "from-internal").
	// CDR(direction) is the most reliable if set consistently.
	if direction := getHeader(msg, "CDR(direction)"); direction != "" {
		return strings.ToLower(direction) // e.g., "inbound", "outbound"
	}
	logger.Warnf("Call direction for UniqueID %s is 'unknown'. Set CDR(direction) in dialplan for accuracy.", getHeader(msg, "Uniqueid"))
	return models.CallDirectionUnknown // Default from models
}

func determineDisposition(msg *goami2.Message) string {
	// AMAFlags: 0 = UNKNOWN, 1 = ANSWERED, 2 = NOANSWER, 3 = BUSY
	// (Note: Some docs say DOCUMENTATION instead of UNKNOWN for 0)
	amaFlagsStr := getHeader(msg, "AMAFlags")
	if amaFlagsStr != "" {
		switch amaFlagsStr {
		case "1":
			return models.CallDispositionAnswered
		case "2":
			return models.CallDispositionNoAnswer
		case "3":
			return models.CallDispositionBusy
		default:
			// Use Cause-txt if AMAFlags is UNKNOWN or not one of the main ones
			if causeTxt := getHeader(msg, "Cause-txt"); causeTxt != "" { // Will be replaced by getHeader
				return causeTxt
			}
			return models.CallDispositionFailed // A general fallback
		}
	}
	// Fallback to Cause-txt if AMAFlags is not present
	if causeTxt := getHeader(msg, "Cause-txt"); causeTxt != "" { // Will be replaced by getHeader
		return causeTxt
	}
	return models.CallDispositionFailed // General fallback
}

// getHeader retrieves a specific header value from an AMI message.
// It iterates through the headers slice provided by msg.Headers().
func getHeader(msg *goami2.Message, key string) string {
	for _, h := range msg.Headers() {
		if h.Name == key {
			return h.Value
		}
	}
	return ""
}

// getAllHeadersAsMap converts all headers from an AMI message into a map[string]string.
func getAllHeadersAsMap(msg *goami2.Message) map[string]string {
	headersMap := make(map[string]string)
	for _, h := range msg.Headers() {
		headersMap[h.Name] = h.Value
	}
	return headersMap
}

