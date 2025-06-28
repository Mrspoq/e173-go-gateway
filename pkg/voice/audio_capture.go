package voice

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "io"
    "sync"
    "time"
)

// AudioFormat represents the format of audio data
type AudioFormat struct {
    SampleRate   int    // Hz (e.g., 8000, 16000)
    Channels     int    // 1 for mono, 2 for stereo
    BitsPerSample int   // 8, 16, 24, or 32
    Codec        string // "pcm", "ulaw", "alaw", etc.
}

// AudioStream represents a stream of audio data from a call
type AudioStream struct {
    CallID      string
    Direction   CallDirection
    Format      AudioFormat
    reader      io.Reader
    buffer      *bytes.Buffer
    mu          sync.Mutex
    closed      bool
}

// AudioCapture handles capturing audio from SIP/RTP streams
type AudioCapture struct {
    streams    map[string]*AudioStream
    recordings map[string]*Recording
    mu         sync.RWMutex
}

// Recording represents a saved audio recording
type Recording struct {
    CallID       string
    StartTime    time.Time
    EndTime      time.Time
    Duration     time.Duration
    Format       AudioFormat
    Size         int64
    Buffer       *bytes.Buffer
    Classification *Classification
}

// NewAudioCapture creates a new audio capture service
func NewAudioCapture() *AudioCapture {
    return &AudioCapture{
        streams:    make(map[string]*AudioStream),
        recordings: make(map[string]*Recording),
    }
}

// StartCapture begins capturing audio for a call
func (a *AudioCapture) StartCapture(callID string, direction CallDirection, format AudioFormat) (*AudioStream, error) {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    if _, exists := a.streams[callID]; exists {
        return nil, fmt.Errorf("capture already started for call %s", callID)
    }
    
    stream := &AudioStream{
        CallID:    callID,
        Direction: direction,
        Format:    format,
        buffer:    new(bytes.Buffer),
        closed:    false,
    }
    
    a.streams[callID] = stream
    
    // Start recording
    recording := &Recording{
        CallID:    callID,
        StartTime: time.Now(),
        Format:    format,
        Buffer:    new(bytes.Buffer),
    }
    a.recordings[callID] = recording
    
    return stream, nil
}

// StopCapture stops capturing audio for a call
func (a *AudioCapture) StopCapture(callID string) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    stream, exists := a.streams[callID]
    if !exists {
        return fmt.Errorf("no capture found for call %s", callID)
    }
    
    stream.mu.Lock()
    stream.closed = true
    stream.mu.Unlock()
    
    // Finalize recording
    if recording, exists := a.recordings[callID]; exists {
        recording.EndTime = time.Now()
        recording.Duration = recording.EndTime.Sub(recording.StartTime)
        recording.Size = int64(recording.Buffer.Len())
    }
    
    delete(a.streams, callID)
    
    return nil
}

// WriteAudio writes audio data to the stream
func (a *AudioCapture) WriteAudio(callID string, data []byte) error {
    a.mu.RLock()
    stream, exists := a.streams[callID]
    recording, hasRecording := a.recordings[callID]
    a.mu.RUnlock()
    
    if !exists {
        return fmt.Errorf("no stream found for call %s", callID)
    }
    
    stream.mu.Lock()
    defer stream.mu.Unlock()
    
    if stream.closed {
        return fmt.Errorf("stream closed for call %s", callID)
    }
    
    // Write to stream buffer
    _, err := stream.buffer.Write(data)
    if err != nil {
        return err
    }
    
    // Also write to recording
    if hasRecording {
        recording.Buffer.Write(data)
    }
    
    return nil
}

// GetStream returns the audio stream for a call
func (a *AudioCapture) GetStream(callID string) (*AudioStream, error) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    
    stream, exists := a.streams[callID]
    if !exists {
        return nil, fmt.Errorf("no stream found for call %s", callID)
    }
    
    return stream, nil
}

// GetRecording returns the recording for a call
func (a *AudioCapture) GetRecording(callID string) (*Recording, error) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    
    recording, exists := a.recordings[callID]
    if !exists {
        return nil, fmt.Errorf("no recording found for call %s", callID)
    }
    
    return recording, nil
}

// Read implements io.Reader for AudioStream
func (s *AudioStream) Read(p []byte) (n int, err error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if s.closed && s.buffer.Len() == 0 {
        return 0, io.EOF
    }
    
    return s.buffer.Read(p)
}

// ConvertToPCM converts audio data to PCM format for STT
func ConvertToPCM(data []byte, format AudioFormat) ([]byte, error) {
    switch format.Codec {
    case "pcm":
        return data, nil
        
    case "ulaw":
        return convertULawToPCM(data), nil
        
    case "alaw":
        return convertALawToPCM(data), nil
        
    default:
        return nil, fmt.Errorf("unsupported codec: %s", format.Codec)
    }
}

// convertULawToPCM converts μ-law encoded audio to PCM
func convertULawToPCM(ulaw []byte) []byte {
    pcm := make([]byte, len(ulaw)*2)
    
    for i, sample := range ulaw {
        // μ-law to PCM conversion
        sign := (sample & 0x80) >> 7
        exponent := (sample & 0x70) >> 4
        mantissa := sample & 0x0F
        
        value := int16(mantissa) << (exponent + 3)
        if sign == 0 {
            value = -value
        }
        
        // Write as 16-bit PCM
        binary.LittleEndian.PutUint16(pcm[i*2:], uint16(value))
    }
    
    return pcm
}

// convertALawToPCM converts A-law encoded audio to PCM
func convertALawToPCM(alaw []byte) []byte {
    pcm := make([]byte, len(alaw)*2)
    
    for i, sample := range alaw {
        // A-law to PCM conversion
        sample ^= 0x55 // Toggle even bits
        
        sign := (sample & 0x80) >> 7
        exponent := (sample & 0x70) >> 4
        mantissa := sample & 0x0F
        
        value := int16(mantissa) << (exponent + 4)
        if sign == 0 {
            value = -value
        }
        
        // Write as 16-bit PCM
        binary.LittleEndian.PutUint16(pcm[i*2:], uint16(value))
    }
    
    return pcm
}

// SimpleAudioRecorder implements the AudioRecorder interface
type SimpleAudioRecorder struct {
    capture *AudioCapture
}

// NewSimpleAudioRecorder creates a simple audio recorder
func NewSimpleAudioRecorder() *SimpleAudioRecorder {
    return &SimpleAudioRecorder{
        capture: NewAudioCapture(),
    }
}

// StartRecording begins recording a call
func (r *SimpleAudioRecorder) StartRecording(callID string) error {
    // Default format for telephony
    format := AudioFormat{
        SampleRate:    8000,
        Channels:      1,
        BitsPerSample: 16,
        Codec:         "pcm",
    }
    
    _, err := r.capture.StartCapture(callID, DirectionIncoming, format)
    return err
}

// StopRecording stops recording and returns the audio
func (r *SimpleAudioRecorder) StopRecording(callID string) (io.ReadCloser, error) {
    err := r.capture.StopCapture(callID)
    if err != nil {
        return nil, err
    }
    
    recording, err := r.capture.GetRecording(callID)
    if err != nil {
        return nil, err
    }
    
    return io.NopCloser(bytes.NewReader(recording.Buffer.Bytes())), nil
}

// SaveRecording saves the recording with classification data
func (r *SimpleAudioRecorder) SaveRecording(callID string, classification *Classification) error {
    recording, err := r.capture.GetRecording(callID)
    if err != nil {
        return err
    }
    
    recording.Classification = classification
    
    // In a real implementation, this would save to disk or cloud storage
    // For now, just log it
    fmt.Printf("Saved recording for call %s: %d bytes, category: %s\n",
        callID, recording.Size, classification.Category)
    
    return nil
}