package handlers

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
    
    "github.com/e173-gateway/e173_go_gateway/pkg/analytics"
    "github.com/e173-gateway/e173_go_gateway/pkg/cache"
)

// AnalyticsHandler handles analytics and reporting endpoints
type AnalyticsHandler struct {
    analytics *analytics.Service
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(db *pgxpool.Pool, cacheService *cache.CacheService) *AnalyticsHandler {
    // Create analytics service with cache adapter
    cacheAdapter := analytics.NewCacheAdapter(cacheService)
    analyticsService := analytics.NewService(db, cacheAdapter)
    
    return &AnalyticsHandler{
        analytics: analyticsService,
    }
}

// GetCallAnalytics returns call analytics for a time period
func (h *AnalyticsHandler) GetCallAnalytics(c *gin.Context) {
    // Get time range from query params
    startStr := c.Query("start")
    endStr := c.Query("end")
    
    // Default to last 7 days if not specified
    endTime := time.Now()
    startTime := endTime.AddDate(0, 0, -7)
    
    if startStr != "" {
        if t, err := time.Parse("2006-01-02", startStr); err == nil {
            startTime = t
        }
    }
    
    if endStr != "" {
        if t, err := time.Parse("2006-01-02", endStr); err == nil {
            endTime = t.Add(24 * time.Hour) // Include full day
        }
    }
    
    analytics, err := h.analytics.GetCallAnalytics(c.Request.Context(), startTime, endTime)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get call analytics",
        })
        return
    }
    
    c.JSON(http.StatusOK, analytics)
}

// GetSIMAnalytics returns SIM card analytics
func (h *AnalyticsHandler) GetSIMAnalytics(c *gin.Context) {
    analytics, err := h.analytics.GetSIMAnalytics(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get SIM analytics",
        })
        return
    }
    
    c.JSON(http.StatusOK, analytics)
}

// GetSpamAnalytics returns spam detection analytics
func (h *AnalyticsHandler) GetSpamAnalytics(c *gin.Context) {
    // Get days parameter (default 30)
    days := 30
    if d := c.Query("days"); d != "" {
        if val, err := time.ParseDuration(d + "d"); err == nil {
            days = int(val.Hours() / 24)
        }
    }
    
    analytics, err := h.analytics.GetSpamAnalytics(c.Request.Context(), days)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get spam analytics",
        })
        return
    }
    
    c.JSON(http.StatusOK, analytics)
}

// GetDashboardData returns aggregated dashboard data
func (h *AnalyticsHandler) GetDashboardData(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Get multiple analytics in parallel
    type result struct {
        calls *analytics.CallAnalytics
        sims  *analytics.SIMAnalytics
        spam  *analytics.SpamAnalytics
        err   error
    }
    
    ch := make(chan result, 3)
    
    // Call analytics (last 24 hours)
    go func() {
        endTime := time.Now()
        startTime := endTime.Add(-24 * time.Hour)
        calls, err := h.analytics.GetCallAnalytics(ctx, startTime, endTime)
        ch <- result{calls: calls, err: err}
    }()
    
    // SIM analytics
    go func() {
        sims, err := h.analytics.GetSIMAnalytics(ctx)
        ch <- result{sims: sims, err: err}
    }()
    
    // Spam analytics (last 7 days)
    go func() {
        spam, err := h.analytics.GetSpamAnalytics(ctx, 7)
        ch <- result{spam: spam, err: err}
    }()
    
    // Collect results
    var callAnalytics *analytics.CallAnalytics
    var simAnalytics *analytics.SIMAnalytics
    var spamAnalytics *analytics.SpamAnalytics
    
    for i := 0; i < 3; i++ {
        r := <-ch
        if r.err != nil {
            continue
        }
        if r.calls != nil {
            callAnalytics = r.calls
        }
        if r.sims != nil {
            simAnalytics = r.sims
        }
        if r.spam != nil {
            spamAnalytics = r.spam
        }
    }
    
    // Build dashboard response
    dashboard := gin.H{
        "overview": gin.H{
            "total_calls_24h":    0,
            "active_sims":        0,
            "spam_blocked_7d":    0,
            "revenue_generated":  0.0,
        },
        "call_analytics": callAnalytics,
        "sim_analytics":  simAnalytics,
        "spam_analytics": spamAnalytics,
        "last_updated":   time.Now(),
    }
    
    // Update overview stats if data available
    if callAnalytics != nil {
        dashboard["overview"].(gin.H)["total_calls_24h"] = callAnalytics.TotalCalls
    }
    if simAnalytics != nil {
        dashboard["overview"].(gin.H)["active_sims"] = simAnalytics.ActiveSIMs
    }
    if spamAnalytics != nil {
        dashboard["overview"].(gin.H)["spam_blocked_7d"] = spamAnalytics.SpamCallsBlocked
    }
    
    c.JSON(http.StatusOK, dashboard)
}