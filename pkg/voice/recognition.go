package voice

import (
    "context"
    "fmt"
    "io"
    "time"
)

// RecognitionService handles voice recognition for both incoming and outgoing calls
type RecognitionService struct {
    sttProvider    STTProvider
    classifier     CallClassifier
    actionEngine   ActionEngine
    audioRecorder  AudioRecorder
}

// STTProvider defines the interface for speech-to-text services
type STTProvider interface {
    TranscribeAudio(ctx context.Context, audio io.Reader) (*Transcript, error)
    GetLanguage() string
    IsAvailable() bool
}

// Transcript represents the result of speech-to-text conversion
type Transcript struct {
    Text       string    `json:"text"`
    Language   string    `json:"language"`
    Duration   float64   `json:"duration"`
    Confidence float64   `json:"confidence"`
    Timestamp  time.Time `json:"timestamp"`
    Segments   []Segment `json:"segments,omitempty"`
}

// Segment represents a time-segmented portion of the transcript
type Segment struct {
    Text       string  `json:"text"`
    StartTime  float64 `json:"start_time"`
    EndTime    float64 `json:"end_time"`
    Confidence float64 `json:"confidence"`
}

// CallClassifier analyzes transcripts to determine call type
type CallClassifier interface {
    ClassifyCall(ctx context.Context, transcript *Transcript) (*Classification, error)
    GetConfidenceThreshold() float64
}

// Classification represents the analysis result of a call
type Classification struct {
    Category    CallCategory `json:"category"`
    Confidence  float64      `json:"confidence"`
    Action      CallAction   `json:"action"`
    Reason      string       `json:"reason"`
    Keywords    []string     `json:"keywords,omitempty"`
    RiskScore   float64      `json:"risk_score"`
}

// CallCategory represents the type of call detected
type CallCategory string

const (
    CategorySpamRobocall   CallCategory = "SPAM_ROBOCALL"
    CategorySIMBlocked     CallCategory = "SIM_BLOCKED"
    CategoryVoicemail      CallCategory = "VOICEMAIL"
    CategoryNormalCall     CallCategory = "NORMAL_CALL"
    CategoryOperatorIVR    CallCategory = "OPERATOR_IVR"
    CategoryLowCredit      CallCategory = "LOW_CREDIT"
    CategoryUnknown        CallCategory = "UNKNOWN"
)

// CallAction represents what to do with the call
type CallAction string

const (
    ActionRouteToAI      CallAction = "ROUTE_TO_AI"
    ActionFlagSIM        CallAction = "FLAG_SIM"
    ActionNormalRouting  CallAction = "NORMAL_ROUTING"
    ActionBlockCall      CallAction = "BLOCK_CALL"
    ActionRecordForReview CallAction = "RECORD_FOR_REVIEW"
)

// ActionEngine executes decisions based on classification
type ActionEngine interface {
    ExecuteAction(ctx context.Context, callID string, classification *Classification) error
    RouteToAI(ctx context.Context, callID string) error
    FlagSIM(ctx context.Context, simID string, reason string) error
}

// AudioRecorder handles call recording for analysis
type AudioRecorder interface {
    StartRecording(callID string) error
    StopRecording(callID string) (io.ReadCloser, error)
    SaveRecording(callID string, classification *Classification) error
}

// RecognitionResult combines all analysis results
type RecognitionResult struct {
    CallID         string          `json:"call_id"`
    Direction      CallDirection   `json:"direction"`
    Transcript     *Transcript     `json:"transcript"`
    Classification *Classification `json:"classification"`
    ActionTaken    CallAction      `json:"action_taken"`
    Timestamp      time.Time       `json:"timestamp"`
}

// CallDirection indicates if call is incoming or outgoing
type CallDirection string

const (
    DirectionIncoming CallDirection = "INCOMING"
    DirectionOutgoing CallDirection = "OUTGOING"
)

// NewRecognitionService creates a new voice recognition service
func NewRecognitionService(stt STTProvider, classifier CallClassifier, action ActionEngine, recorder AudioRecorder) *RecognitionService {
    return &RecognitionService{
        sttProvider:   stt,
        classifier:    classifier,
        actionEngine:  action,
        audioRecorder: recorder,
    }
}

// AnalyzeCall processes a call through the recognition pipeline
func (r *RecognitionService) AnalyzeCall(ctx context.Context, callID string, audio io.Reader, direction CallDirection) (*RecognitionResult, error) {
    // Start recording for potential review
    if err := r.audioRecorder.StartRecording(callID); err != nil {
        return nil, fmt.Errorf("failed to start recording: %w", err)
    }
    
    // Speech to text conversion
    transcript, err := r.sttProvider.TranscribeAudio(ctx, audio)
    if err != nil {
        return nil, fmt.Errorf("transcription failed: %w", err)
    }
    
    // Classify the call based on transcript
    classification, err := r.classifier.ClassifyCall(ctx, transcript)
    if err != nil {
        return nil, fmt.Errorf("classification failed: %w", err)
    }
    
    // Execute action based on classification
    if err := r.actionEngine.ExecuteAction(ctx, callID, classification); err != nil {
        return nil, fmt.Errorf("action execution failed: %w", err)
    }
    
    // Save recording if needed
    if shouldSaveRecording(classification) {
        if err := r.audioRecorder.SaveRecording(callID, classification); err != nil {
            // Log error but don't fail the whole process
            fmt.Printf("Failed to save recording: %v\n", err)
        }
    }
    
    return &RecognitionResult{
        CallID:         callID,
        Direction:      direction,
        Transcript:     transcript,
        Classification: classification,
        ActionTaken:    classification.Action,
        Timestamp:      time.Now(),
    }, nil
}

// shouldSaveRecording determines if we should keep the recording
func shouldSaveRecording(classification *Classification) bool {
    // Save recordings for:
    // - Spam calls (for training)
    // - SIM issues (for debugging)
    // - Low confidence classifications (for review)
    return classification.Category == CategorySpamRobocall ||
           classification.Category == CategorySIMBlocked ||
           classification.Confidence < 0.7
}

// AnalyzeIncomingCall specifically handles incoming call analysis for spam detection
func (r *RecognitionService) AnalyzeIncomingCall(ctx context.Context, callID string, audio io.Reader) (*RecognitionResult, error) {
    return r.AnalyzeCall(ctx, callID, audio, DirectionIncoming)
}

// AnalyzeOutgoingCall specifically handles outgoing call analysis for SIM status
func (r *RecognitionService) AnalyzeOutgoingCall(ctx context.Context, callID string, simID string, audio io.Reader) (*RecognitionResult, error) {
    result, err := r.AnalyzeCall(ctx, callID, audio, DirectionOutgoing)
    if err != nil {
        return nil, err
    }
    
    // Additional handling for SIM-specific issues
    if result.Classification.Category == CategorySIMBlocked || 
       result.Classification.Category == CategoryLowCredit {
        if err := r.actionEngine.FlagSIM(ctx, simID, result.Classification.Reason); err != nil {
            fmt.Printf("Failed to flag SIM %s: %v\n", simID, err)
        }
    }
    
    return result, nil
}