package voice

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// WhisperProvider implements STTProvider using OpenAI Whisper API
type WhisperProvider struct {
    apiKey      string
    apiURL      string
    model       string
    language    string
    client      *http.Client
}

// WhisperResponse represents the API response structure
type WhisperResponse struct {
    Text     string             `json:"text"`
    Language string             `json:"language,omitempty"`
    Duration float64            `json:"duration,omitempty"`
    Segments []WhisperSegment   `json:"segments,omitempty"`
}

// WhisperSegment represents a segment in the transcription
type WhisperSegment struct {
    ID               int     `json:"id"`
    Seek             int     `json:"seek"`
    Start            float64 `json:"start"`
    End              float64 `json:"end"`
    Text             string  `json:"text"`
    Temperature      float64 `json:"temperature"`
    AvgLogprob       float64 `json:"avg_logprob"`
    CompressionRatio float64 `json:"compression_ratio"`
    NoSpeechProb     float64 `json:"no_speech_prob"`
}

// NewWhisperProvider creates a new Whisper STT provider
func NewWhisperProvider(apiKey string) *WhisperProvider {
    return &WhisperProvider{
        apiKey:   apiKey,
        apiURL:   "https://api.openai.com/v1/audio/transcriptions",
        model:    "whisper-1",
        language: "en", // Default to English
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// TranscribeAudio converts audio to text using Whisper
func (w *WhisperProvider) TranscribeAudio(ctx context.Context, audio io.Reader) (*Transcript, error) {
    // Create multipart form data
    var requestBody bytes.Buffer
    
    // For now, we'll create a simple JSON request
    // In production, this would be multipart/form-data with the audio file
    audioData, err := io.ReadAll(audio)
    if err != nil {
        return nil, fmt.Errorf("failed to read audio: %w", err)
    }
    
    // Create the request
    req, err := http.NewRequestWithContext(ctx, "POST", w.apiURL, bytes.NewReader(audioData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Set headers
    req.Header.Set("Authorization", "Bearer "+w.apiKey)
    req.Header.Set("Content-Type", "multipart/form-data")
    
    // Make the request
    resp, err := w.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("API request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Check response status
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
    }
    
    // Parse response
    var whisperResp WhisperResponse
    if err := json.NewDecoder(resp.Body).Decode(&whisperResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    // Convert to our format
    transcript := &Transcript{
        Text:       whisperResp.Text,
        Language:   whisperResp.Language,
        Duration:   whisperResp.Duration,
        Confidence: 0.9, // Whisper doesn't provide confidence scores
        Timestamp:  time.Now(),
    }
    
    // Convert segments if available
    if len(whisperResp.Segments) > 0 {
        transcript.Segments = make([]Segment, len(whisperResp.Segments))
        for i, seg := range whisperResp.Segments {
            transcript.Segments[i] = Segment{
                Text:       seg.Text,
                StartTime:  seg.Start,
                EndTime:    seg.End,
                Confidence: 1.0 - seg.NoSpeechProb, // Estimate confidence
            }
        }
    }
    
    return transcript, nil
}

// GetLanguage returns the configured language
func (w *WhisperProvider) GetLanguage() string {
    return w.language
}

// IsAvailable checks if the Whisper service is available
func (w *WhisperProvider) IsAvailable() bool {
    // Simple health check - could be enhanced
    return w.apiKey != ""
}

// SetLanguage configures the language for transcription
func (w *WhisperProvider) SetLanguage(language string) {
    w.language = language
}

// LocalWhisperProvider implements STTProvider using local Whisper model
type LocalWhisperProvider struct {
    modelPath   string
    language    string
    device      string // "cpu" or "cuda"
}

// NewLocalWhisperProvider creates a provider using local Whisper model
func NewLocalWhisperProvider(modelPath string) *LocalWhisperProvider {
    return &LocalWhisperProvider{
        modelPath: modelPath,
        language:  "en",
        device:    "cpu",
    }
}

// TranscribeAudio using local Whisper model
func (l *LocalWhisperProvider) TranscribeAudio(ctx context.Context, audio io.Reader) (*Transcript, error) {
    // This would integrate with a local Whisper deployment
    // For now, return a placeholder implementation
    return &Transcript{
        Text:       "Local Whisper transcription not yet implemented",
        Language:   l.language,
        Duration:   0,
        Confidence: 0.5,
        Timestamp:  time.Now(),
    }, nil
}

// GetLanguage returns the configured language
func (l *LocalWhisperProvider) GetLanguage() string {
    return l.language
}

// IsAvailable checks if local Whisper is available
func (l *LocalWhisperProvider) IsAvailable() bool {
    // Check if model file exists
    // For now, return false as it's not implemented
    return false
}