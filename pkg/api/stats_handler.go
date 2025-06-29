package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/e173-gateway/e173_go_gateway/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StatsHandler struct {
	modemRepo   repository.ModemRepository
	simCardRepo repository.SIMCardRepository
	cdrRepo     repository.CdrRepository
	gatewayRepo repository.GatewayRepository
	logger      *logrus.Logger
}

func NewStatsHandler(modemRepo repository.ModemRepository, simCardRepo repository.SIMCardRepository, cdrRepo repository.CdrRepository, gatewayRepo repository.GatewayRepository, logger *logrus.Logger) *StatsHandler {
	return &StatsHandler{
		modemRepo:   modemRepo,
		simCardRepo: simCardRepo,
		cdrRepo:     cdrRepo,
		gatewayRepo: gatewayRepo,
		logger:      logger,
	}
}

// GetModemStats returns HTML snippet for modem statistics card
func (h *StatsHandler) GetModemStats(c *gin.Context) {
	ctx := context.Background()
	
	modems, err := h.modemRepo.GetAllModems(ctx)
	if err != nil {
		h.logger.Errorf("Error fetching modems for stats: %v", err)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div class="text-center">
			<div class="inline-flex items-center justify-center w-10 h-10 bg-red-100 dark:bg-red-900 rounded-full mb-2">
				<svg class="w-6 h-6 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>
			</div>
			<h3 class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Modems</h3>
			<p class="mt-1 text-xl font-semibold text-red-600 dark:text-red-400">Error</p>
			<p class="text-xs text-gray-500 dark:text-gray-400">Failed to load</p>
		</div>`)
		return
	}

	onlineCount := 0
	total := len(modems)
	
	for _, modem := range modems {
		if modem.Status == "online" || modem.Status == "active" {
			onlineCount++
		}
	}

	// Determine status color
	statusColor := "green"
	if onlineCount == 0 {
		statusColor = "red"
	} else if float64(onlineCount)/float64(total) < 0.8 {
		statusColor = "amber"
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
	<div class="text-center">
		<div class="inline-flex items-center justify-center w-10 h-10 bg-`+statusColor+`-100 dark:bg-`+statusColor+`-900 rounded-full mb-2">
			<svg class="w-6 h-6 text-`+statusColor+`-600 dark:text-`+statusColor+`-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"></path>
			</svg>
		</div>
		<h3 class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Modems</h3>
		<p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">`+strconv.Itoa(onlineCount)+`/`+strconv.Itoa(total)+`</p>
		<p class="text-xs text-`+statusColor+`-600 dark:text-`+statusColor+`-400">Online</p>
	</div>`)
}

// GetSIMStats returns HTML snippet for SIM statistics card
func (h *StatsHandler) GetSIMStats(c *gin.Context) {
	ctx := context.Background()
	
	sims, err := h.simCardRepo.GetAllSIMCards(ctx)
	if err != nil {
		h.logger.Errorf("Error fetching SIMs for stats: %v", err)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div id="sim-balance-card" hx-get="/api/stats/sims" hx-trigger="every 5s" hx-swap="outerHTML" class="dashboard-card">
			<div class="flex items-center">
				<div class="flex-shrink-0">
					<div class="w-8 h-8 bg-red-100 rounded-md flex items-center justify-center">
						<svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					</div>
				</div>
				<div class="ml-5 w-0 flex-1">
					<dl>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">SIMs Low Balance</dt>
						<dd class="flex items-baseline">
							<div class="text-2xl font-semibold text-red-600">Error</div>
						</dd>
					</dl>
				</div>
			</div>
		</div>`)
		return
	}

	lowBalanceCount := 0
	lowBalanceThreshold := 5.0 // Define threshold for low balance
	
	for _, sim := range sims {
		if sim.Balance.Valid && sim.Balance.Float64 < lowBalanceThreshold {
			lowBalanceCount++
		}
	}

	statusColor := "green"
	if lowBalanceCount > 10 {
		statusColor = "red"
	} else if lowBalanceCount > 5 {
		statusColor = "amber"
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
	<div id="sim-balance-card" hx-get="/api/stats/sims" hx-trigger="every 5s" hx-swap="outerHTML" class="dashboard-card">
		<div class="flex items-center">
			<div class="flex-shrink-0">
				<div class="w-8 h-8 bg-`+statusColor+`-100 rounded-md flex items-center justify-center">
					<svg class="w-5 h-5 text-`+statusColor+`-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
				</div>
			</div>
			<div class="ml-5 w-0 flex-1">
				<dl>
					<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">SIMs Low Balance</dt>
					<dd class="flex items-baseline">
						<div class="text-2xl font-semibold text-gray-900 dark:text-white">`+strconv.Itoa(lowBalanceCount)+`</div>
					</dd>
				</dl>
			</div>
		</div>
	</div>`)
}

// GetCallStats returns HTML snippet for call statistics card
func (h *StatsHandler) GetCallStats(c *gin.Context) {
	ctx := context.Background()
	
	// Get recent CDRs to analyze call patterns
	// For live calls, we'll count calls in the last 5 minutes that don't have end times
	cdrs, err := h.cdrRepo.GetRecentCDRs(ctx, 50) // Get last 50 CDRs
	if err != nil {
		h.logger.Errorf("Error fetching CDRs for call stats: %v", err)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div id="live-calls-card" hx-get="/api/stats/calls" hx-trigger="every 3s" hx-swap="outerHTML" class="dashboard-card">
			<div class="flex items-center">
				<div class="flex-shrink-0">
					<div class="w-8 h-8 bg-red-100 rounded-md flex items-center justify-center">
						<svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					</div>
				</div>
				<div class="ml-5 w-0 flex-1">
					<dl>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Live Calls</dt>
						<dd class="flex items-baseline">
							<div class="text-2xl font-semibold text-red-600">Error</div>
						</dd>
					</dl>
				</div>
			</div>
		</div>`)
		return
	}

	// Count active calls (calls without end time in the last 10 minutes)
	activeCalls := 0
	recentCallsCount := 0
	currentTime := time.Now()
	
	for _, cdr := range cdrs {
		// Count calls from the last hour as "recent"
		if cdr.StartTime != nil && currentTime.Sub(*cdr.StartTime) <= time.Hour {
			recentCallsCount++
			
			// If no end time or very recent start time (within 10 minutes), consider it active
			if cdr.EndTime == nil || currentTime.Sub(*cdr.StartTime) <= 10*time.Minute {
				activeCalls++
			}
		}
	}
	
	// Determine status color based on call volume
	statusColor := "blue"
	statusText := "active"
	if activeCalls > 10 {
		statusColor = "green"
		statusText = "high activity"
	} else if activeCalls == 0 {
		statusColor = "gray"
		statusText = "idle"
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
	<div id="live-calls-card" hx-get="/api/stats/calls" hx-trigger="every 3s" hx-swap="outerHTML" class="dashboard-card">
		<div class="flex items-center">
			<div class="flex-shrink-0">
				<div class="w-8 h-8 bg-`+statusColor+`-100 rounded-md flex items-center justify-center">
					<svg class="w-5 h-5 text-`+statusColor+`-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
					</svg>
				</div>
			</div>
			<div class="ml-5 w-0 flex-1">
				<dl>
					<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Live Calls</dt>
					<dd class="flex items-baseline">
						<div class="text-2xl font-semibold text-gray-900 dark:text-white">`+strconv.Itoa(activeCalls)+`</div>
						<div class="ml-2 text-sm text-gray-500 dark:text-gray-400">`+statusText+`</div>
					</dd>
				</dl>
			</div>
		</div>
	</div>`)
}

// GetSpamStats returns HTML snippet for spam statistics card
func (h *StatsHandler) GetSpamStats(c *gin.Context) {
	ctx := context.Background()
	
	// Get recent CDRs to analyze spam patterns
	cdrs, err := h.cdrRepo.GetRecentCDRs(ctx, 200) // Get more CDRs for better spam analysis
	if err != nil {
		h.logger.Errorf("Error fetching CDRs for spam stats: %v", err)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div id="spam-blocked-card" hx-get="/api/stats/spam" hx-trigger="every 10s" hx-swap="outerHTML" class="dashboard-card">
			<div class="flex items-center">
				<div class="flex-shrink-0">
					<div class="w-8 h-8 bg-red-100 rounded-md flex items-center justify-center">
						<svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					</div>
				</div>
				<div class="ml-5 w-0 flex-1">
					<dl>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Spam Blocked Today</dt>
						<dd class="flex items-baseline">
							<div class="text-2xl font-semibold text-red-600">Error</div>
						</dd>
					</dl>
				</div>
			</div>
		</div>`)
		return
	}

	// Analyze CDRs for spam patterns
	spamCount := 0
	shortCallCount := 0
	frequentCallerMap := make(map[string]int)
	currentTime := time.Now()
	todayStart := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	
	for _, cdr := range cdrs {
		// Only analyze calls from today
		if cdr.StartTime == nil || (*cdr.StartTime).Before(todayStart) {
			continue
		}
		
		// Count frequency of each caller
		if cdr.CallerIDNum != nil {
			frequentCallerMap[*cdr.CallerIDNum]++
		}
		
		// Detect short calls (potential spam pattern)
		if cdr.StartTime != nil && cdr.EndTime != nil {
			duration := (*cdr.EndTime).Sub(*cdr.StartTime)
			if duration < 5*time.Second { // Calls shorter than 5 seconds
				shortCallCount++
			}
		}
		
		// Check if call was marked as FAILED or BUSY (potential spam indicators)
		if cdr.Disposition != nil && (*cdr.Disposition == "FAILED" || *cdr.Disposition == "BUSY") {
			// Additional spam indicator
		}
	}
	
	// Count potential spam based on frequency (more than 10 calls from same number today)
	for _, callCount := range frequentCallerMap {
		if callCount > 10 {
			spamCount += callCount - 5 // Consider excess calls as spam
		}
	}
	
	// Add short calls to spam count (likely robocalls)
	spamCount += shortCallCount

	// Determine status color
	statusColor := "green"
	statusText := "blocked today"
	if spamCount > 50 {
		statusColor = "red"
		statusText = "high spam activity"
	} else if spamCount > 20 {
		statusColor = "amber"
		statusText = "moderate spam"
	} else if spamCount == 0 {
		statusText = "no spam detected"
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
	<div id="spam-blocked-card" hx-get="/api/stats/spam" hx-trigger="every 10s" hx-swap="outerHTML" class="dashboard-card">
		<div class="flex items-center">
			<div class="flex-shrink-0">
				<div class="w-8 h-8 bg-`+statusColor+`-100 rounded-md flex items-center justify-center">
					<svg class="w-5 h-5 text-`+statusColor+`-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L18.364 5.636M5.636 18.364l12.728-12.728"></path>
					</svg>
				</div>
			</div>
			<div class="ml-5 w-0 flex-1">
				<dl>
					<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Potential Spam Today</dt>
					<dd class="flex items-baseline">
						<div class="text-2xl font-semibold text-gray-900 dark:text-white">`+strconv.Itoa(spamCount)+`</div>
						<div class="ml-2 text-sm text-gray-500 dark:text-gray-400">`+statusText+`</div>
					</dd>
				</dl>
			</div>
		</div>
	</div>`)
}

// GetGatewayStats returns HTML snippet for gateway statistics card
func (h *StatsHandler) GetGatewayStats(c *gin.Context) {
	ctx := context.Background()
	
	total, online, _, err := h.gatewayRepo.GetGatewayStats(ctx)
	if err != nil {
		h.logger.Errorf("Error fetching gateway stats: %v", err)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div id="gateway-stats-card" hx-get="/api/stats/gateways" hx-trigger="every 5s" hx-swap="outerHTML" class="dashboard-card">
			<div class="flex items-center">
				<div class="flex-shrink-0">
					<div class="w-8 h-8 bg-red-100 rounded-md flex items-center justify-center">
						<svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
						</svg>
					</div>
				</div>
				<div class="ml-5 w-0 flex-1">
					<dl>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Gateways</dt>
						<dd class="flex items-baseline">
							<div class="text-2xl font-semibold text-red-600">Error</div>
						</dd>
					</dl>
				</div>
			</div>
		</div>`)
		return
	}

	// Determine status color
	statusColor := "blue"
	statusIcon := "M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"
	if online == 0 && total > 0 {
		statusColor = "red"
	} else if total > 0 && float64(online)/float64(total) < 0.8 {
		statusColor = "amber"
	} else if online > 0 {
		statusColor = "green"
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, `
	<div id="gateway-stats-card" hx-get="/api/stats/gateways" hx-trigger="every 5s" hx-swap="outerHTML" class="dashboard-card">
		<div class="flex items-center">
			<div class="flex-shrink-0">
				<div class="w-8 h-8 bg-`+statusColor+`-100 rounded-md flex items-center justify-center">
					<svg class="w-5 h-5 text-`+statusColor+`-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="`+statusIcon+`"></path>
					</svg>
				</div>
			</div>
			<div class="ml-5 w-0 flex-1">
				<dl>
					<dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">Gateways</dt>
					<dd class="flex items-baseline">
						<div class="text-2xl font-semibold text-gray-900 dark:text-white">`+strconv.Itoa(online)+` / `+strconv.Itoa(total)+`</div>
						<div class="ml-2 text-sm text-gray-500 dark:text-gray-400">online</div>
					</dd>
				</dl>
			</div>
		</div>
	</div>`)
}
