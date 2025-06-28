package cache

import (
    "context"
    "fmt"
    "time"
    
    "github.com/e173-gateway/e173_go_gateway/pkg/models"
)

// CacheService provides high-level caching operations for the application
type CacheService struct {
    redis *RedisClient
}

// NewCacheService creates a new cache service
func NewCacheService(redis *RedisClient) *CacheService {
    return &CacheService{
        redis: redis,
    }
}

// Keys for different cache types
const (
    KeyWhatsAppValidation = "whatsapp:validation:%s"     // phone number
    KeyCallStats          = "stats:calls:%s"             // time period (daily, hourly)
    KeySIMStatus          = "sim:status:%s"              // sim id
    KeyModemStatus        = "modem:status:%s"            // modem id
    KeyUserSession        = "user:session:%s"            // session token
    KeyRateLimit          = "ratelimit:%s:%s"            // endpoint:ip
    KeySpamPattern        = "spam:pattern:%s"            // pattern hash
    KeyAPIResponse        = "api:response:%s"            // request hash
    KeyDashboardStats     = "dashboard:stats"
    KeySystemMetrics      = "system:metrics"
)

// CacheWhatsAppValidation caches a WhatsApp validation result
func (s *CacheService) CacheWhatsAppValidation(ctx context.Context, phoneNumber string, result *models.ValidationResult) error {
    key := fmt.Sprintf(KeyWhatsAppValidation, phoneNumber)
    ttl := 24 * time.Hour // 24-hour cache as per spec
    
    return s.redis.SetWithTTL(ctx, key, result, ttl)
}

// GetWhatsAppValidation retrieves a cached WhatsApp validation
func (s *CacheService) GetWhatsAppValidation(ctx context.Context, phoneNumber string) (*models.ValidationResult, error) {
    key := fmt.Sprintf(KeyWhatsAppValidation, phoneNumber)
    var result models.ValidationResult
    
    err := s.redis.Get(ctx, key, &result)
    if err != nil {
        return nil, err
    }
    
    return &result, nil
}

// CallStats represents aggregated call statistics
type CallStats struct {
    TotalCalls      int64     `json:"total_calls"`
    SuccessfulCalls int64     `json:"successful_calls"`
    FailedCalls     int64     `json:"failed_calls"`
    TotalDuration   int64     `json:"total_duration"`
    AvgDuration     float64   `json:"avg_duration"`
    UniqueNumbers   int64     `json:"unique_numbers"`
    LastUpdated     time.Time `json:"last_updated"`
}

// UpdateCallStats updates call statistics atomically
func (s *CacheService) UpdateCallStats(ctx context.Context, period string, success bool, duration int64) error {
    key := fmt.Sprintf(KeyCallStats, period)
    
    // Use pipeline for atomic updates
    pipe := s.redis.GetClient().Pipeline()
    
    pipe.HIncrBy(ctx, s.redis.buildKey(key), "total_calls", 1)
    if success {
        pipe.HIncrBy(ctx, s.redis.buildKey(key), "successful_calls", 1)
    } else {
        pipe.HIncrBy(ctx, s.redis.buildKey(key), "failed_calls", 1)
    }
    pipe.HIncrBy(ctx, s.redis.buildKey(key), "total_duration", duration)
    pipe.HSet(ctx, s.redis.buildKey(key), "last_updated", time.Now().Unix())
    
    // Set expiry for hourly stats (keep for 7 days)
    if period == "hourly" {
        pipe.Expire(ctx, s.redis.buildKey(key), 7*24*time.Hour)
    }
    
    _, err := pipe.Exec(ctx)
    return err
}

// GetCallStats retrieves call statistics
func (s *CacheService) GetCallStats(ctx context.Context, period string) (*CallStats, error) {
    key := fmt.Sprintf(KeyCallStats, period)
    
    data, err := s.redis.GetClient().HGetAll(ctx, s.redis.buildKey(key)).Result()
    if err != nil {
        return nil, err
    }
    
    if len(data) == 0 {
        return nil, ErrCacheMiss
    }
    
    stats := &CallStats{}
    // Parse the hash data
    if v, ok := data["total_calls"]; ok {
        fmt.Sscanf(v, "%d", &stats.TotalCalls)
    }
    if v, ok := data["successful_calls"]; ok {
        fmt.Sscanf(v, "%d", &stats.SuccessfulCalls)
    }
    if v, ok := data["failed_calls"]; ok {
        fmt.Sscanf(v, "%d", &stats.FailedCalls)
    }
    if v, ok := data["total_duration"]; ok {
        fmt.Sscanf(v, "%d", &stats.TotalDuration)
    }
    if v, ok := data["last_updated"]; ok {
        var ts int64
        fmt.Sscanf(v, "%d", &ts)
        stats.LastUpdated = time.Unix(ts, 0)
    }
    
    // Calculate average duration
    if stats.TotalCalls > 0 {
        stats.AvgDuration = float64(stats.TotalDuration) / float64(stats.TotalCalls)
    }
    
    return stats, nil
}

// SIMStatusCache represents cached SIM status
type SIMStatusCache struct {
    SIMID       string    `json:"sim_id"`
    Status      string    `json:"status"`
    Balance     float64   `json:"balance"`
    LastUsed    time.Time `json:"last_used"`
    SignalLevel int       `json:"signal_level"`
    Issues      []string  `json:"issues,omitempty"`
}

// CacheSIMStatus caches SIM card status
func (s *CacheService) CacheSIMStatus(ctx context.Context, simID string, status *SIMStatusCache) error {
    key := fmt.Sprintf(KeySIMStatus, simID)
    ttl := 5 * time.Minute // Short TTL for real-time status
    
    return s.redis.SetWithTTL(ctx, key, status, ttl)
}

// GetSIMStatus retrieves cached SIM status
func (s *CacheService) GetSIMStatus(ctx context.Context, simID string) (*SIMStatusCache, error) {
    key := fmt.Sprintf(KeySIMStatus, simID)
    var status SIMStatusCache
    
    err := s.redis.Get(ctx, key, &status)
    if err != nil {
        return nil, err
    }
    
    return &status, nil
}

// ModemStatusCache represents cached modem status
type ModemStatusCache struct {
    ModemID     string    `json:"modem_id"`
    Status      string    `json:"status"`
    CurrentCall string    `json:"current_call,omitempty"`
    Temperature float64   `json:"temperature"`
    Uptime      int64     `json:"uptime"`
    LastPing    time.Time `json:"last_ping"`
}

// CacheModemStatus caches modem status
func (s *CacheService) CacheModemStatus(ctx context.Context, modemID string, status *ModemStatusCache) error {
    key := fmt.Sprintf(KeyModemStatus, modemID)
    ttl := 1 * time.Minute // Very short TTL for real-time monitoring
    
    return s.redis.SetWithTTL(ctx, key, status, ttl)
}

// GetModemStatus retrieves cached modem status
func (s *CacheService) GetModemStatus(ctx context.Context, modemID string) (*ModemStatusCache, error) {
    key := fmt.Sprintf(KeyModemStatus, modemID)
    var status ModemStatusCache
    
    err := s.redis.Get(ctx, key, &status)
    if err != nil {
        return nil, err
    }
    
    return &status, nil
}

// RateLimitCheck checks and updates rate limit
func (s *CacheService) RateLimitCheck(ctx context.Context, endpoint, identifier string, limit int64, window time.Duration) (bool, error) {
    key := fmt.Sprintf(KeyRateLimit, endpoint, identifier)
    
    // Increment counter
    count, err := s.redis.Increment(ctx, key)
    if err != nil {
        return false, err
    }
    
    // Set expiry on first request
    if count == 1 {
        if err := s.redis.Expire(ctx, key, window); err != nil {
            return false, err
        }
    }
    
    return count <= limit, nil
}

// GetRateLimitRemaining returns remaining requests in the current window
func (s *CacheService) GetRateLimitRemaining(ctx context.Context, endpoint, identifier string, limit int64) (int64, time.Duration, error) {
    key := fmt.Sprintf(KeyRateLimit, endpoint, identifier)
    
    // Get current count
    var count int64
    err := s.redis.Get(ctx, key, &count)
    if err != nil {
        if err == ErrCacheMiss {
            return limit, 0, nil
        }
        return 0, 0, err
    }
    
    // Get TTL
    ttl, err := s.redis.GetTTL(ctx, key)
    if err != nil {
        return 0, 0, err
    }
    
    remaining := limit - count
    if remaining < 0 {
        remaining = 0
    }
    
    return remaining, ttl, nil
}

// DashboardStats represents cached dashboard statistics
type DashboardStats struct {
    ModemStats    map[string]interface{} `json:"modem_stats"`
    SIMStats      map[string]interface{} `json:"sim_stats"`
    CallStats     map[string]interface{} `json:"call_stats"`
    SystemStats   map[string]interface{} `json:"system_stats"`
    LastUpdated   time.Time              `json:"last_updated"`
}

// CacheDashboardStats caches dashboard statistics
func (s *CacheService) CacheDashboardStats(ctx context.Context, stats *DashboardStats) error {
    ttl := 30 * time.Second // Short TTL for dashboard freshness
    return s.redis.SetWithTTL(ctx, KeyDashboardStats, stats, ttl)
}

// GetDashboardStats retrieves cached dashboard statistics
func (s *CacheService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
    var stats DashboardStats
    err := s.redis.Get(ctx, KeyDashboardStats, &stats)
    if err != nil {
        return nil, err
    }
    return &stats, nil
}

// InvalidatePattern invalidates all cache entries matching a pattern
func (s *CacheService) InvalidatePattern(ctx context.Context, pattern string) error {
    // This would need to be implemented based on Redis SCAN command
    // For now, return nil
    return nil
}

// WarmupCache pre-loads frequently accessed data
func (s *CacheService) WarmupCache(ctx context.Context) error {
    // This would be called on startup to pre-populate cache
    // with frequently accessed data from the database
    return nil
}

// GetCacheStats returns cache statistics
func (s *CacheService) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
    info, err := s.redis.GetClient().Info(ctx, "stats").Result()
    if err != nil {
        return nil, err
    }
    
    // Parse Redis INFO stats
    stats := map[string]interface{}{
        "redis_info": info,
        "connected":  true,
    }
    
    return stats, nil
}