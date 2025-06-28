package analytics

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    "github.com/jackc/pgx/v4/pgxpool"
)

// Service provides analytics and reporting functionality
type Service struct {
    db    *pgxpool.Pool
    cache CacheProvider
}

// CacheProvider interface for analytics caching
type CacheProvider interface {
    Get(ctx context.Context, key string, dest interface{}) error
    SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error
}

// NewService creates a new analytics service
func NewService(db *pgxpool.Pool, cache CacheProvider) *Service {
    return &Service{
        db:    db,
        cache: cache,
    }
}

// CallAnalytics represents call analytics data
type CallAnalytics struct {
    Period          string                 `json:"period"`
    TotalCalls      int64                  `json:"total_calls"`
    SuccessfulCalls int64                  `json:"successful_calls"`
    FailedCalls     int64                  `json:"failed_calls"`
    SpamCalls       int64                  `json:"spam_calls"`
    TotalMinutes    float64                `json:"total_minutes"`
    AvgCallDuration float64                `json:"avg_call_duration"`
    ASR             float64                `json:"asr"` // Answer Seizure Ratio
    ACD             float64                `json:"acd"` // Average Call Duration
    ByHour          []HourlyStats          `json:"by_hour,omitempty"`
    ByOperator      map[string]OperatorStats `json:"by_operator,omitempty"`
    TopDestinations []DestinationStats     `json:"top_destinations,omitempty"`
}

// HourlyStats represents hourly call statistics
type HourlyStats struct {
    Hour       int     `json:"hour"`
    Calls      int64   `json:"calls"`
    Minutes    float64 `json:"minutes"`
    ASR        float64 `json:"asr"`
}

// OperatorStats represents per-operator statistics
type OperatorStats struct {
    Operator   string  `json:"operator"`
    Calls      int64   `json:"calls"`
    Minutes    float64 `json:"minutes"`
    ASR        float64 `json:"asr"`
    AvgQuality float64 `json:"avg_quality"`
}

// DestinationStats represents top destination statistics
type DestinationStats struct {
    Destination string  `json:"destination"`
    Country     string  `json:"country"`
    Calls       int64   `json:"calls"`
    Minutes     float64 `json:"minutes"`
}

// GetCallAnalytics retrieves call analytics for a time period
func (s *Service) GetCallAnalytics(ctx context.Context, startTime, endTime time.Time) (*CallAnalytics, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("analytics:calls:%d:%d", startTime.Unix(), endTime.Unix())
    var analytics CallAnalytics
    if err := s.cache.Get(ctx, cacheKey, &analytics); err == nil {
        return &analytics, nil
    }
    
    // Main query for overall stats
    query := `
        SELECT 
            COUNT(*) as total_calls,
            COUNT(CASE WHEN disposition = 'ANSWERED' THEN 1 END) as successful_calls,
            COUNT(CASE WHEN disposition != 'ANSWERED' THEN 1 END) as failed_calls,
            COUNT(CASE WHEN voice_category = 'SPAM_ROBOCALL' THEN 1 END) as spam_calls,
            COALESCE(SUM(duration) / 60.0, 0) as total_minutes,
            COALESCE(AVG(duration), 0) as avg_duration
        FROM call_detail_records
        WHERE start_time >= $1 AND start_time < $2
    `
    
    err := s.db.QueryRow(ctx, query, startTime, endTime).Scan(
        &analytics.TotalCalls,
        &analytics.SuccessfulCalls,
        &analytics.FailedCalls,
        &analytics.SpamCalls,
        &analytics.TotalMinutes,
        &analytics.AvgCallDuration,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get call analytics: %w", err)
    }
    
    // Calculate ASR and ACD
    if analytics.TotalCalls > 0 {
        analytics.ASR = float64(analytics.SuccessfulCalls) / float64(analytics.TotalCalls) * 100
    }
    if analytics.SuccessfulCalls > 0 {
        analytics.ACD = analytics.TotalMinutes / float64(analytics.SuccessfulCalls) * 60 // in seconds
    }
    
    analytics.Period = fmt.Sprintf("%s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
    
    // Get hourly breakdown if period is <= 7 days
    if endTime.Sub(startTime) <= 7*24*time.Hour {
        analytics.ByHour, _ = s.getHourlyStats(ctx, startTime, endTime)
    }
    
    // Get operator breakdown
    analytics.ByOperator, _ = s.getOperatorStats(ctx, startTime, endTime)
    
    // Get top destinations
    analytics.TopDestinations, _ = s.getTopDestinations(ctx, startTime, endTime, 10)
    
    // Cache the result
    s.cache.SetWithTTL(ctx, cacheKey, analytics, 5*time.Minute)
    
    return &analytics, nil
}

// getHourlyStats retrieves hourly statistics
func (s *Service) getHourlyStats(ctx context.Context, startTime, endTime time.Time) ([]HourlyStats, error) {
    query := `
        SELECT 
            EXTRACT(HOUR FROM start_time) as hour,
            COUNT(*) as calls,
            COALESCE(SUM(duration) / 60.0, 0) as minutes,
            COUNT(CASE WHEN disposition = 'ANSWERED' THEN 1 END) * 100.0 / COUNT(*) as asr
        FROM call_detail_records
        WHERE start_time >= $1 AND start_time < $2
        GROUP BY hour
        ORDER BY hour
    `
    
    rows, err := s.db.Query(ctx, query, startTime, endTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var stats []HourlyStats
    for rows.Next() {
        var stat HourlyStats
        err := rows.Scan(&stat.Hour, &stat.Calls, &stat.Minutes, &stat.ASR)
        if err != nil {
            continue
        }
        stats = append(stats, stat)
    }
    
    return stats, nil
}

// getOperatorStats retrieves per-operator statistics
func (s *Service) getOperatorStats(ctx context.Context, startTime, endTime time.Time) (map[string]OperatorStats, error) {
    query := `
        SELECT 
            COALESCE(s.operator_name, 'Unknown') as operator,
            COUNT(*) as calls,
            COALESCE(SUM(c.duration) / 60.0, 0) as minutes,
            COUNT(CASE WHEN c.disposition = 'ANSWERED' THEN 1 END) * 100.0 / COUNT(*) as asr
        FROM call_detail_records c
        LEFT JOIN sim_cards s ON c.sim_id = s.id
        WHERE c.start_time >= $1 AND c.start_time < $2
        GROUP BY operator
    `
    
    rows, err := s.db.Query(ctx, query, startTime, endTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    stats := make(map[string]OperatorStats)
    for rows.Next() {
        var operator string
        var stat OperatorStats
        err := rows.Scan(&operator, &stat.Calls, &stat.Minutes, &stat.ASR)
        if err != nil {
            continue
        }
        stat.Operator = operator
        stats[operator] = stat
    }
    
    return stats, nil
}

// getTopDestinations retrieves top call destinations
func (s *Service) getTopDestinations(ctx context.Context, startTime, endTime time.Time, limit int) ([]DestinationStats, error) {
    query := `
        SELECT 
            SUBSTRING(destination_number FROM 1 FOR 5) as prefix,
            COUNT(*) as calls,
            COALESCE(SUM(duration) / 60.0, 0) as minutes
        FROM call_detail_records
        WHERE start_time >= $1 AND start_time < $2
            AND destination_number IS NOT NULL
        GROUP BY prefix
        ORDER BY calls DESC
        LIMIT $3
    `
    
    rows, err := s.db.Query(ctx, query, startTime, endTime, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var stats []DestinationStats
    for rows.Next() {
        var stat DestinationStats
        err := rows.Scan(&stat.Destination, &stat.Calls, &stat.Minutes)
        if err != nil {
            continue
        }
        // Map prefix to country (simplified)
        stat.Country = getCountryFromPrefix(stat.Destination)
        stats = append(stats, stat)
    }
    
    return stats, nil
}

// SIMAnalytics represents SIM card analytics
type SIMAnalytics struct {
    TotalSIMs        int64             `json:"total_sims"`
    ActiveSIMs       int64             `json:"active_sims"`
    BlockedSIMs      int64             `json:"blocked_sims"`
    LowCreditSIMs    int64             `json:"low_credit_sims"`
    AvgDailyUsage    float64           `json:"avg_daily_usage"`
    ByOperator       map[string]int64  `json:"by_operator"`
    ReplacementQueue int64             `json:"replacement_queue"`
    CreditStatus     CreditDistribution `json:"credit_status"`
}

// CreditDistribution shows SIM credit distribution
type CreditDistribution struct {
    Under10  int64 `json:"under_10"`
    Under50  int64 `json:"under_50"`
    Under100 int64 `json:"under_100"`
    Over100  int64 `json:"over_100"`
}

// GetSIMAnalytics retrieves SIM card analytics
func (s *Service) GetSIMAnalytics(ctx context.Context) (*SIMAnalytics, error) {
    // Try cache first
    cacheKey := "analytics:sims:current"
    var analytics SIMAnalytics
    if err := s.cache.Get(ctx, cacheKey, &analytics); err == nil {
        return &analytics, nil
    }
    
    // Main query
    query := `
        SELECT 
            COUNT(*) as total,
            COUNT(CASE WHEN status = 'active' THEN 1 END) as active,
            COUNT(CASE WHEN status = 'blocked' THEN 1 END) as blocked,
            COUNT(CASE WHEN current_credit < 10 THEN 1 END) as low_credit
        FROM sim_cards
    `
    
    err := s.db.QueryRow(ctx, query).Scan(
        &analytics.TotalSIMs,
        &analytics.ActiveSIMs,
        &analytics.BlockedSIMs,
        &analytics.LowCreditSIMs,
    )
    if err != nil {
        return nil, err
    }
    
    // Get operator distribution
    analytics.ByOperator, _ = s.getSIMsByOperator(ctx)
    
    // Get replacement queue count
    s.db.QueryRow(ctx, 
        "SELECT COUNT(*) FROM sim_replacement_queue WHERE status = 'pending'",
    ).Scan(&analytics.ReplacementQueue)
    
    // Get credit distribution
    analytics.CreditStatus, _ = s.getCreditDistribution(ctx)
    
    // Cache the result
    s.cache.SetWithTTL(ctx, cacheKey, analytics, 1*time.Minute)
    
    return &analytics, nil
}

// getSIMsByOperator gets SIM count by operator
func (s *Service) getSIMsByOperator(ctx context.Context) (map[string]int64, error) {
    query := `
        SELECT operator_name, COUNT(*) 
        FROM sim_cards 
        WHERE operator_name IS NOT NULL
        GROUP BY operator_name
    `
    
    rows, err := s.db.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    result := make(map[string]int64)
    for rows.Next() {
        var operator sql.NullString
        var count int64
        if err := rows.Scan(&operator, &count); err == nil && operator.Valid {
            result[operator.String] = count
        }
    }
    
    return result, nil
}

// getCreditDistribution gets SIM credit distribution
func (s *Service) getCreditDistribution(ctx context.Context) (CreditDistribution, error) {
    var dist CreditDistribution
    
    query := `
        SELECT 
            COUNT(CASE WHEN current_credit < 10 THEN 1 END) as under_10,
            COUNT(CASE WHEN current_credit >= 10 AND current_credit < 50 THEN 1 END) as under_50,
            COUNT(CASE WHEN current_credit >= 50 AND current_credit < 100 THEN 1 END) as under_100,
            COUNT(CASE WHEN current_credit >= 100 THEN 1 END) as over_100
        FROM sim_cards
        WHERE current_credit IS NOT NULL
    `
    
    err := s.db.QueryRow(ctx, query).Scan(
        &dist.Under10,
        &dist.Under50,
        &dist.Under100,
        &dist.Over100,
    )
    
    return dist, err
}

// SpamAnalytics represents spam detection analytics
type SpamAnalytics struct {
    TotalSpamCalls    int64                    `json:"total_spam_calls"`
    SpamCallsBlocked  int64                    `json:"spam_calls_blocked"`
    SpamCallsToAI     int64                    `json:"spam_calls_to_ai"`
    TimeWasted        float64                  `json:"time_wasted_minutes"`
    TopSpamSources    []SpamSource             `json:"top_spam_sources"`
    DetectionMethods  map[string]int64         `json:"detection_methods"`
    AIAgentStats      map[string]AgentStats    `json:"ai_agent_stats"`
}

// SpamSource represents a source of spam calls
type SpamSource struct {
    Number    string `json:"number"`
    Calls     int64  `json:"calls"`
    LastSeen  time.Time `json:"last_seen"`
}

// AgentStats represents AI agent performance
type AgentStats struct {
    CallsHandled  int64   `json:"calls_handled"`
    AvgDuration   float64 `json:"avg_duration"`
    TotalMinutes  float64 `json:"total_minutes"`
}

// GetSpamAnalytics retrieves spam detection analytics
func (s *Service) GetSpamAnalytics(ctx context.Context, days int) (*SpamAnalytics, error) {
    startTime := time.Now().AddDate(0, 0, -days)
    
    // Try cache first
    cacheKey := fmt.Sprintf("analytics:spam:%d", days)
    var analytics SpamAnalytics
    if err := s.cache.Get(ctx, cacheKey, &analytics); err == nil {
        return &analytics, nil
    }
    
    // Get spam call counts
    query := `
        SELECT 
            COUNT(*) as total_spam,
            COUNT(CASE WHEN voice_action = 'BLOCK_CALL' THEN 1 END) as blocked,
            COUNT(CASE WHEN routed_to_ai = true THEN 1 END) as to_ai
        FROM sip_calls
        WHERE voice_category = 'SPAM_ROBOCALL'
            AND created_at >= $1
    `
    
    err := s.db.QueryRow(ctx, query, startTime).Scan(
        &analytics.TotalSpamCalls,
        &analytics.SpamCallsBlocked,
        &analytics.SpamCallsToAI,
    )
    if err != nil {
        return nil, err
    }
    
    // Get time wasted by AI agents
    s.db.QueryRow(ctx,
        `SELECT COALESCE(SUM(duration_seconds) / 60.0, 0) 
         FROM ai_agent_interactions 
         WHERE start_time >= $1`,
        startTime,
    ).Scan(&analytics.TimeWasted)
    
    // Get top spam sources
    analytics.TopSpamSources, _ = s.getTopSpamSources(ctx, startTime, 10)
    
    // Get detection methods breakdown
    analytics.DetectionMethods, _ = s.getDetectionMethods(ctx, startTime)
    
    // Get AI agent stats
    analytics.AIAgentStats, _ = s.getAIAgentStats(ctx, startTime)
    
    // Cache the result
    s.cache.SetWithTTL(ctx, cacheKey, analytics, 10*time.Minute)
    
    return &analytics, nil
}

// getTopSpamSources retrieves top spam sources
func (s *Service) getTopSpamSources(ctx context.Context, since time.Time, limit int) ([]SpamSource, error) {
    query := `
        SELECT 
            caller_id_num,
            COUNT(*) as calls,
            MAX(created_at) as last_seen
        FROM sip_calls
        WHERE voice_category = 'SPAM_ROBOCALL'
            AND created_at >= $1
            AND caller_id_num IS NOT NULL
        GROUP BY caller_id_num
        ORDER BY calls DESC
        LIMIT $2
    `
    
    rows, err := s.db.Query(ctx, query, since, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var sources []SpamSource
    for rows.Next() {
        var source SpamSource
        var number sql.NullString
        err := rows.Scan(&number, &source.Calls, &source.LastSeen)
        if err == nil && number.Valid {
            source.Number = number.String
            sources = append(sources, source)
        }
    }
    
    return sources, nil
}

// getDetectionMethods gets spam detection method breakdown
func (s *Service) getDetectionMethods(ctx context.Context, since time.Time) (map[string]int64, error) {
    // This would analyze which detection methods caught spam
    // For now, return sample data
    return map[string]int64{
        "whatsapp_validation": 145,
        "voice_recognition":   89,
        "blacklist":          34,
        "pattern_matching":   67,
    }, nil
}

// getAIAgentStats gets AI agent performance stats
func (s *Service) getAIAgentStats(ctx context.Context, since time.Time) (map[string]AgentStats, error) {
    query := `
        SELECT 
            agent_id,
            COUNT(*) as calls,
            AVG(duration_seconds) as avg_duration,
            SUM(duration_seconds) / 60.0 as total_minutes
        FROM ai_agent_interactions
        WHERE start_time >= $1
        GROUP BY agent_id
    `
    
    rows, err := s.db.Query(ctx, query, since)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    stats := make(map[string]AgentStats)
    for rows.Next() {
        var agentID string
        var stat AgentStats
        err := rows.Scan(&agentID, &stat.CallsHandled, &stat.AvgDuration, &stat.TotalMinutes)
        if err == nil {
            stats[agentID] = stat
        }
    }
    
    return stats, nil
}

// getCountryFromPrefix maps phone prefixes to countries (simplified)
func getCountryFromPrefix(prefix string) string {
    prefixMap := map[string]string{
        "+1":   "USA/Canada",
        "+44":  "UK",
        "+234": "Nigeria",
        "+27":  "South Africa",
        "+254": "Kenya",
        "+91":  "India",
        "+86":  "China",
    }
    
    for p, country := range prefixMap {
        if len(prefix) >= len(p) && prefix[:len(p)] == p {
            return country
        }
    }
    
    return "Unknown"
}