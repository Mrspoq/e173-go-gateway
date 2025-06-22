package api

import (
	"context"
	"net/http"
	"time"

	"github.com/e173-gateway/e173_go_gateway/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UIHandler struct {
	modemRepo   repository.ModemRepository
	simCardRepo repository.SIMCardRepository
	cdrRepo     repository.CdrRepository
	logger      *logrus.Logger
}

func NewUIHandler(modemRepo repository.ModemRepository, simCardRepo repository.SIMCardRepository, cdrRepo repository.CdrRepository, logger *logrus.Logger) *UIHandler {
	return &UIHandler{
		modemRepo:   modemRepo,
		simCardRepo: simCardRepo,
		cdrRepo:     cdrRepo,
		logger:      logger,
	}
}

// GetModemList returns HTML snippet for modem list
func (h *UIHandler) GetModemList(c *gin.Context) {
	ctx := context.Background()
	
	modems, err := h.modemRepo.GetAllModems(ctx)
	if err != nil {
		h.logger.Errorf("Error fetching modems for UI: %v", err)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div class="text-center py-4 text-red-600">
			<p>Error loading modems</p>
			<button hx-get="/ui/modems/list" hx-target="#modem-list" hx-swap="innerHTML" 
					class="mt-2 text-sm text-indigo-600 hover:text-indigo-900">Retry</button>
		</div>`)
		return
	}

	if len(modems) == 0 {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div class="text-center py-4 text-gray-500 dark:text-gray-400">
			<p>No modems detected</p>
			<p class="text-sm mt-1">Check hardware connections and try refreshing</p>
		</div>`)
		return
	}

	html := ""
	for _, modem := range modems {
		statusClass := "modem-offline"
		statusText := "Offline"
		statusIcon := `<svg class="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
		</svg>`

		if modem.Status == "online" || modem.Status == "active" {
			statusClass = "modem-online"
			statusText = "Online"
			statusIcon = `<svg class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
			</svg>`
		} else if modem.Status == "warning" {
			statusClass = "modem-warning"
			statusText = "Warning"
			statusIcon = `<svg class="w-4 h-4 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
			</svg>`
		}

		imei := "N/A"
		if modem.IMEI != nil {
			imei = *modem.IMEI
		}

		operator := "Unknown"
		if modem.NetworkOperatorName != nil {
			operator = *modem.NetworkOperatorName
		}

		signal := "N/A"
		if modem.SignalStrengthDBM != nil {
			signal = string(rune(*modem.SignalStrengthDBM)) + " dBm"
		}

		html += `
		<div class="modem-card ` + statusClass + ` mb-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-md">
			<div class="flex items-center justify-between">
				<div class="flex items-center">
					` + statusIcon + `
					<div class="ml-3">
						<p class="text-sm font-medium text-gray-900 dark:text-white">` + modem.DevicePath + `</p>
						<p class="text-xs text-gray-500 dark:text-gray-400">IMEI: ` + imei + `</p>
					</div>
				</div>
				<div class="text-right">
					<p class="text-sm text-gray-900 dark:text-white">` + statusText + `</p>
					<p class="text-xs text-gray-500 dark:text-gray-400">` + operator + `</p>
					<p class="text-xs text-gray-500 dark:text-gray-400">` + signal + `</p>
				</div>
			</div>
		</div>`
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// GetCDRStream returns HTML snippet for recent CDR entries
func (h *UIHandler) GetCDRStream(c *gin.Context) {
	// For now, return mock CDR data since we need to implement real CDR fetching
	// This will be enhanced when we have actual call data coming through AMI
	
	currentTime := time.Now()
	
	// Generate some mock recent activity
	html := `
	<tr class="animate-pulse">
		<td class="px-4 py-2 text-sm text-gray-900 dark:text-gray-100">` + currentTime.Format("15:04:05") + `</td>
		<td class="px-4 py-2 text-sm text-gray-900 dark:text-gray-100">+1234567890</td>
		<td class="px-4 py-2 text-sm text-gray-900 dark:text-gray-100">+0987654321</td>
		<td class="px-4 py-2 text-sm">
			<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800 dark:bg-blue-800 dark:text-blue-100">
				In Progress
			</span>
		</td>
	</tr>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// GetAlerts returns HTML snippet for system alerts
func (h *UIHandler) GetAlerts(c *gin.Context) {
	// For now, return empty alerts
	// This will be enhanced when we implement alerting system
	
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, "")
}
