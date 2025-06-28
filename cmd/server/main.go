package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
	// Import internal packages
	"github.com/e173-gateway/e173_go_gateway/pkg/logging"
	"github.com/e173-gateway/e173_go_gateway/pkg/config"
	"github.com/e173-gateway/e173_go_gateway/pkg/auth"
	"github.com/e173-gateway/e173_go_gateway/pkg/database" // Import database package
	"github.com/e173-gateway/e173_go_gateway/pkg/repository" // Import repository package
	simhandler "github.com/e173-gateway/e173_go_gateway/pkg/api" // Import API handlers
	
	// Import enterprise modules
	enterpriseRepo "github.com/e173-gateway/e173_go_gateway/internal/repository"
	"github.com/e173-gateway/e173_go_gateway/internal/service"
	"github.com/e173-gateway/e173_go_gateway/internal/handlers"
	adapter "github.com/e173-gateway/e173_go_gateway/internal/database"
)

// loadTemplates recursively loads all HTML templates from the templates directory
func loadTemplates(templatesDir string) (*template.Template, error) {
	t := template.New("")
	
	// First, load the base template
	basePath := filepath.Join(templatesDir, "base.tmpl")
	baseContent, err := os.ReadFile(basePath)
	if err != nil {
		return nil, fmt.Errorf("reading base template: %w", err)
	}
	
	_, err = t.New("base.tmpl").Parse(string(baseContent))
	if err != nil {
		return nil, fmt.Errorf("parsing base template: %w", err)
	}
	
	// Then load all other templates
	err = filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Skip base template (already loaded)
		if path == basePath {
			return nil
		}
		
		// Process both .html and .tmpl files
		if !strings.HasSuffix(path, ".html") && !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		
		// Get relative path from templates directory
		relPath, err := filepath.Rel(templatesDir, path)
		if err != nil {
			return err
		}
		
		// Read template file
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		// Parse template with its relative path as name
		_, err = t.New(relPath).Parse(string(b))
		if err != nil {
			return fmt.Errorf("parsing template %s: %w", relPath, err)
		}
		
		logging.Logger.Debugf("Loaded template: %s", relPath)
		
		return nil
	})
	
	return t, err
}

func main() {
	// Load .env file. Errors are ignored if .env is not found.
	if err := godotenv.Load(); err != nil {
		// We can log this if we want, but often it's fine if .env is missing (e.g., in production using real env vars)
		// log.Printf("No .env file found or error loading it: %v", err)
	}

	// Load application configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logging.InitLogger(cfg.LogLevel, cfg.LogFormat)

	logging.Logger.Info("Application configuration loaded successfully.")

	// Initialize database connection
	dbPool, err := database.NewDBPool(cfg.DatabaseURL)
	if err != nil {
		logging.Logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close() // Ensure DB connection is closed on shutdown
	logging.Logger.Info("Successfully connected to the database.")

	// Initialize repositories
	modemRepo := repository.NewPostgresModemRepository(dbPool)
	simCardRepo := repository.NewPostgresSIMCardRepository(dbPool)
	cdrRepo := repository.NewPostgresCdrRepository(dbPool) // Initialize CDR Repository
	gatewayRepo := repository.NewPostgresGatewayRepository(dbPool)
	
	// Initialize JWT service
	var jwtService *auth.JWTService
	if cfg.JWTSecret == "" {
		logging.Logger.Fatal("JWT_SECRET is required in environment variables")
	}
	
	tokenExpiry, err := time.ParseDuration(cfg.JWTExpiry)
	if err != nil {
		logging.Logger.Fatalf("Invalid JWT_EXPIRY format: %v", err)
	}
	
	refreshExpiry, err := time.ParseDuration(cfg.RefreshExpiry)
	if err != nil {
		logging.Logger.Fatalf("Invalid REFRESH_EXPIRY format: %v", err)
	}
	
	jwtService = auth.NewJWTService(cfg.JWTSecret, tokenExpiry, refreshExpiry)
	logging.Logger.Info("JWT authentication service initialized")
	
	// Create sqlx adapter for enterprise repositories
	sqlxDB, err := adapter.CreateSQLXAdapter(dbPool)
	if err != nil {
		logging.Logger.Fatalf("Failed to create sqlx adapter: %v", err)
	}
	defer adapter.CloseAdapter(sqlxDB) // Ensure adapter is closed on shutdown
	
	// Initialize enterprise repositories
	userRepo := enterpriseRepo.NewPostgresUserRepository(dbPool) // User repo uses pgxpool directly
	customerRepo := enterpriseRepo.NewPostgresCustomerRepository(sqlxDB)
	paymentRepo := enterpriseRepo.NewPostgresPaymentRepository(sqlxDB)
	systemRepo := enterpriseRepo.NewPostgresSystemRepository(sqlxDB)

	// Initialize API Handlers
	simAPIHandler := simhandler.NewSIMCardHandler(simCardRepo)
	statsHandler := simhandler.NewStatsHandler(modemRepo, simCardRepo, cdrRepo, gatewayRepo, logging.Logger)
	gatewayHandler := simhandler.NewGatewayHandler(gatewayRepo, logging.Logger)
	
	// Initialize enterprise services
	authService := service.NewPostgresAuthService(userRepo, systemRepo)
	customerService := service.NewPostgresCustomerService(customerRepo, paymentRepo, systemRepo)
	
	// Initialize enterprise handlers
	authHandlers := handlers.NewAuthHandlers(authService, customerService)
	customerHandlers := handlers.NewCustomerHandlers(customerService)

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	router := gin.Default()

	// Load HTML templates from all directories
	templatesDir := "templates"
	tmpl, err := loadTemplates(templatesDir)
	if err != nil {
		logging.Logger.Fatalf("Failed to load templates: %v", err)
	}
	router.SetHTMLTemplate(tmpl)
	
	// Serve static files
	router.Static("/static", "./web/static")

	// Simple health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"timestamp": time.Now(),
		})
	})
	
	// Test route to debug template/routing issues
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test route working - routing is OK")
	})
	
	// Template test route to debug template loading
	router.GET("/template-test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "simple_dashboard.tmpl", gin.H{
			"title": "Template Test",
		})
	})
	
	// Debug endpoint to check template loading
	router.GET("/template-debug/:name", func(c *gin.Context) {
		templateName := c.Param("name")
		c.HTML(http.StatusOK, templateName, gin.H{
			"title": "Debug: " + templateName,
		})
	})

	// Stats endpoint for real-time dashboard
	router.GET("/api/stats", func(c *gin.Context) {
		// TODO: Replace with real database queries
		stats := gin.H{
			"modems": gin.H{
				"total":    12,
				"online":   10,
				"offline":  2,
				"calling":  3,
			},
			"sims": gin.H{
				"total":      12,
				"active":     10,
				"blocked":    1,
				"low_credit": 2,
			},
			"calls": gin.H{
				"today_total":    247,
				"today_minutes":  1853,
				"active_calls":   3,
				"last_24h_total": 3891,
			},
			"system": gin.H{
				"uptime":     "2h 15m",
				"cpu_usage":  "23%",
				"memory":     "1.2GB",
				"disk_free":  "45GB",
			},
			"timestamp": time.Now(),
		}
		c.JSON(http.StatusOK, stats)
	})

	// HTMX Stats Cards partial - Let each card load individually  
	router.GET("/api/stats/cards", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6" hx-get="/api/v1/stats/modems" hx-trigger="load, every 5s" hx-swap="innerHTML"></div>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6" hx-get="/api/v1/stats/sims" hx-trigger="load, every 5s" hx-swap="innerHTML"></div>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6" hx-get="/api/v1/stats/calls" hx-trigger="load, every 5s" hx-swap="innerHTML"></div>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6" hx-get="/api/v1/stats/spam" hx-trigger="load, every 5s" hx-swap="innerHTML"></div>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6" hx-get="/api/v1/stats/gateways" hx-trigger="load, every 5s" hx-swap="innerHTML"></div>`)
	})

	// Direct stats endpoints for backward compatibility (template might call these)
	router.GET("/api/stats/calls", statsHandler.GetCallStats)
	router.GET("/api/stats/spam", statsHandler.GetSpamStats)
	router.GET("/api/stats/modems", statsHandler.GetModemStats)
	router.GET("/api/stats/sims", statsHandler.GetSIMStats)
	router.GET("/api/stats/gateways", statsHandler.GetGatewayStats)

	// CDR ticker endpoint - returns recent call records
	router.GET("/api/cdr/recent", func(c *gin.Context) {
		// TODO: Replace with real CDR queries from database
		recentCDRs := []gin.H{
			{
				"id":          "CDR-001",
				"timestamp":   time.Now().Add(-2 * time.Minute).Format("15:04:05"),
				"from":        "+1234567890",
				"to":          "+0987654321",
				"duration":    "00:02:15",
				"disposition": "ANSWERED",
				"modem":       "G01",
			},
			{
				"id":          "CDR-002", 
				"timestamp":   time.Now().Add(-5 * time.Minute).Format("15:04:05"),
				"from":        "+1111222333",
				"to":          "+4444555666",
				"duration":    "00:01:45",
				"disposition": "ANSWERED",
				"modem":       "G03",
			},
			{
				"id":          "CDR-003",
				"timestamp":   time.Now().Add(-8 * time.Minute).Format("15:04:05"),
				"from":        "+7777888999",
				"to":          "+1010101010",
				"duration":    "00:00:12",
				"disposition": "NO ANSWER",
				"modem":       "G05",
			},
		}
		c.JSON(http.StatusOK, gin.H{"cdrs": recentCDRs})
	})

	// HTMX CDR List partial
	router.GET("/api/cdr/recent/list", func(c *gin.Context) {
		ctx := context.Background()
		
		// Get recent CDRs from the database (limit to last 10)
		cdrs, err := cdrRepo.GetRecentCDRs(ctx, 10)
		if err != nil {
			logging.Logger.Errorf("Error fetching recent CDRs: %v", err)
			c.Header("Content-Type", "text/html")
			c.String(http.StatusOK, `<div class="text-center text-red-600 dark:text-red-400 p-4">
				<svg class="w-8 h-8 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>
				Error loading call records
			</div>`)
			return
		}
		
		c.Header("Content-Type", "text/html")
		html := ``
		
		if len(cdrs) == 0 {
			html = `<div class="text-center text-gray-500 dark:text-gray-400 p-4">
				<svg class="w-8 h-8 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
				</svg>
				No recent calls
			</div>`
		} else {
			for _, cdr := range cdrs {
				disposition := "UNKNOWN"
				if cdr.Disposition != nil {
					disposition = *cdr.Disposition
				}
				statusColor := "green"
				statusBg := "bg-green-100 dark:bg-green-800"
				
				if disposition == "NO ANSWER" || disposition == "BUSY" {
					statusColor = "yellow"
					statusBg = "bg-yellow-100 dark:bg-yellow-800"
				} else if disposition == "FAILED" || disposition == "CONGESTION" {
					statusColor = "red"
					statusBg = "bg-red-100 dark:bg-red-800"
				} else {
					statusColor = "green"
				}
				
				// Format duration from seconds to MM:SS
				durationText := "00:00"
				if cdr.Duration != nil {
					minutes := *cdr.Duration / 60
					seconds := *cdr.Duration % 60
					durationText = fmt.Sprintf("%02d:%02d", minutes, seconds)
				}
				
				// Format timestamp - use CreatedAt as fallback for call time
				timestampText := cdr.CreatedAt.Format("15:04:05")
				if cdr.StartTime != nil {
					timestampText = cdr.StartTime.Format("15:04:05")
				}
				
				// Get caller ID and destination
				callerID := "Unknown"
				if cdr.CallerIDNum != nil {
					callerID = *cdr.CallerIDNum
				}
				
				destination := "Unknown" 
				if cdr.ConnectedLineNum != nil {
					destination = *cdr.ConnectedLineNum
				}
				
				// Get modem name if available
				modemText := "Unknown"
				if cdr.ModemID != nil {
					modemText = fmt.Sprintf("G%02d", *cdr.ModemID)
				}
				
				html += fmt.Sprintf(`
				<div class="flex items-center justify-between p-3 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-lg">
					<div class="flex items-center space-x-4 flex-1 min-w-0">
						<div class="flex-shrink-0">
							<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium %s text-%s-800 dark:text-%s-100">
								%s
							</span>
						</div>
						<div class="flex-1 min-w-0">
							<div class="flex items-center space-x-2">
								<p class="text-sm font-medium text-gray-900 dark:text-white truncate">%s</p>
								<svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5-5 5M6 12h12"></path>
								</svg>
								<p class="text-sm text-gray-600 dark:text-gray-300 truncate">%s</p>
							</div>
							<div class="flex items-center space-x-4 mt-1">
								<p class="text-xs text-gray-500 dark:text-gray-400">%s</p>
								<p class="text-xs text-gray-500 dark:text-gray-400">%s</p>
							</div>
						</div>
					</div>
					<div class="flex-shrink-0 text-right">
						<p class="text-sm font-medium text-gray-900 dark:text-white">%s</p>
					</div>
				</div>`, 
				statusBg, statusColor, statusColor, disposition,
				callerID, destination, timestampText, modemText,
				durationText)
			}
		}
		
		c.String(http.StatusOK, html)
	})

	// Enhanced modem status endpoint with detailed info
	router.GET("/api/modems/status", func(c *gin.Context) {
		ctx := context.Background()
		
		modems, err := modemRepo.GetAllModems(ctx)
		if err != nil {
			logging.Logger.Errorf("Error fetching modems for status: %v", err)
			c.Header("Content-Type", "text/html")
			c.String(http.StatusOK, `<div class="text-center text-red-600 dark:text-red-400 p-4">
				<svg class="w-8 h-8 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
				</svg>
				Error loading modem status
			</div>`)
			return
		}
		
		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusOK, gin.H{"modems": modems})
		} else {
			// Return HTML for HTMX
			c.Header("Content-Type", "text/html")
			html := ``
			
			if len(modems) == 0 {
				html = `<div class="text-center text-gray-500 dark:text-gray-400 p-4">
					<svg class="w-8 h-8 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012-2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
					</svg>
					No modems found
				</div>`
			} else {
				for _, modem := range modems {
					statusColor := "green"
					statusText := modem.Status
					statusBg := "bg-green-100 dark:bg-green-800"
					
					if modem.Status == "offline" || modem.Status == "disconnected" {
						statusColor = "red"
						statusBg = "bg-red-100 dark:bg-red-800"
					} else if modem.Status == "calling" || modem.Status == "busy" {
						statusColor = "blue"
						statusBg = "bg-blue-100 dark:bg-blue-800"
					} else if modem.Status == "idle" || modem.Status == "online" {
						statusColor = "green"
						statusBg = "bg-green-100 dark:bg-green-800"
					} else {
						statusColor = "gray"
						statusBg = "bg-gray-100 dark:bg-gray-800"
					}
					
					// Build signal strength display
					signalText := "Unknown"
					if modem.SignalStrengthDBM != nil {
						signalText = fmt.Sprintf("%.0f dBm", float64(*modem.SignalStrengthDBM))
					}
					
					// Build operator display
					operatorText := "Unknown"
					if modem.NetworkOperatorName != nil {
						operatorText = *modem.NetworkOperatorName
					}
					
					// Build modem name from device path
					modemName := fmt.Sprintf("Modem %s", modem.DevicePath)
					if modem.IMEI != nil {
						modemName = fmt.Sprintf("Modem %s", (*modem.IMEI)[len(*modem.IMEI)-4:]) // Last 4 digits of IMEI
					}
					
					html += fmt.Sprintf(`
					<div class="%s rounded-lg p-4 border-l-4 border-%s-500 mb-4">
						<div class="flex justify-between items-start">
							<div>
								<h4 class="font-medium text-gray-900 dark:text-white">%s</h4>
								<p class="text-sm text-gray-600 dark:text-gray-300">Signal: %s</p>
								<p class="text-sm text-gray-600 dark:text-gray-300">Operator: %s</p>
								<p class="text-sm text-gray-600 dark:text-gray-300">Path: %s</p>
							</div>
							<div class="text-right">
								<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-%s-100 text-%s-800 dark:bg-%s-800 dark:text-%s-100">
									%s
								</span>
							</div>
						</div>
					</div>`, 
					statusBg, statusColor, modemName, signalText, operatorText, modem.DevicePath, 
					statusColor, statusColor, statusColor, statusColor, strings.ToUpper(statusText))
				}
			}
			c.String(http.StatusOK, html)
		}
	})

	// Authentication redirect middleware for HTML pages
	authRedirect := func(c *gin.Context) {
		// Check if user has session cookie
		cookie, err := c.Cookie("session_token")
		if err != nil || cookie == "" {
			// Redirect to login
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		
		// Validate session
		_, err = authService.ValidateSession(cookie)
		if err != nil {
			// Clear invalid cookie and redirect to login
			c.SetCookie("session_token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		
		c.Next()
	}

	router.GET("/", authRedirect, func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard_standalone.tmpl", gin.H{
			"title": "Dashboard - E173 Gateway",
		})
	})

	// API v1 Group
	v1 := router.Group("/api/v1")
	{
		// Modems API endpoint
		v1.GET("/modems", func(c *gin.Context) {
			modems, err := modemRepo.GetAllModems(c.Request.Context())
			if err != nil {
				logging.Logger.Errorf("Error fetching modems: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch modems"})
				return
			}
			
			// Check if this is an HTMX request
			if c.GetHeader("HX-Request") == "true" {
				c.HTML(http.StatusOK, "partials/modem_list.tmpl", modems)
				return
			}
			
			c.JSON(http.StatusOK, modems)
		})

		// SIM Cards API endpoints
		v1.POST("/simcards", simAPIHandler.CreateSIMCard)
		v1.GET("/simcards/:id", simAPIHandler.GetSIMCardByID)
		v1.GET("/simcards", func(c *gin.Context) {
			sims, err := simCardRepo.GetAllSIMCards(c.Request.Context())
			if err != nil {
				logging.Logger.Errorf("Error fetching SIM cards: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SIM cards"})
				return
			}
			
			// Check if this is an HTMX request
			if c.GetHeader("HX-Request") == "true" {
				c.HTML(http.StatusOK, "partials/sim_list.tmpl", sims)
				return
			}
			
			c.JSON(http.StatusOK, sims)
		})
		v1.PUT("/simcards/:id", simAPIHandler.UpdateSIMCard) // New route for updating a SIM
		v1.DELETE("/simcards/:id", simAPIHandler.DeleteSIMCard) // New route for deleting a SIM

		// Enterprise API endpoints
		v1.POST("/enterprises", handlers.WrapHandler(customerHandlers.CreateCustomer))
		v1.GET("/enterprises/:id", handlers.WrapHandler(customerHandlers.GetCustomer))
		v1.GET("/enterprises", handlers.WrapHandler(customerHandlers.ListCustomers))
		v1.PUT("/enterprises/:id", handlers.WrapHandler(customerHandlers.UpdateCustomer))
		v1.DELETE("/enterprises/:id", handlers.WrapHandler(customerHandlers.DeleteCustomer))

		// Authentication API endpoints
		v1.POST("/auth/login", handlers.WrapHandler(authHandlers.Login))
		v1.POST("/auth/register", handlers.WrapHandler(authHandlers.Login)) // Note: register might need separate handler
		v1.GET("/auth/me", handlers.WrapHandler(authHandlers.GetProfile))

		// Customer Management API endpoints
		v1.POST("/customers", handlers.WrapHandler(customerHandlers.CreateCustomer))
		v1.GET("/customers/:id", handlers.WrapHandler(customerHandlers.GetCustomer))
		v1.GET("/customers", handlers.WrapHandler(customerHandlers.ListCustomers))
		v1.PUT("/customers/:id", handlers.WrapHandler(customerHandlers.UpdateCustomer))
		v1.DELETE("/customers/:id", handlers.WrapHandler(customerHandlers.DeleteCustomer))

		// Gateway Management API endpoints
		v1.POST("/gateways", gatewayHandler.CreateGateway)
		v1.GET("/gateways", gatewayHandler.ListGateways)
		v1.GET("/gateways/:id", gatewayHandler.GetGatewayByID)
		v1.PUT("/gateways/:id", gatewayHandler.UpdateGateway)
		v1.DELETE("/gateways/:id", gatewayHandler.DeleteGateway)
		v1.POST("/gateways/heartbeat", gatewayHandler.Heartbeat)
	}

	// V1 Stats endpoints (called by HTMX stats cards)
	v1.GET("/stats/modems", statsHandler.GetModemStats)
	v1.GET("/stats/sims", statsHandler.GetSIMStats)
	v1.GET("/stats/calls", statsHandler.GetCallStats)
	v1.GET("/stats/spam", statsHandler.GetSpamStats)
	v1.GET("/stats/gateways", statsHandler.GetGatewayStats)

	// UI endpoints for HTMX components
	v1.GET("/ui/modems/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "modems/list.html", gin.H{
			"title": "Modem List",
		})
	})
	v1.GET("/ui/modems/refresh", func(c *gin.Context) {
		c.HTML(http.StatusOK, "modems/refresh.html", gin.H{
			"title": "Modem Refresh",
		})
	})
	v1.GET("/ui/cdr/stream", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cdr/stream.html", gin.H{
			"title": "CDR Stream",
		})
	})
	v1.GET("/ui/alerts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "alerts.html", gin.H{
			"title": "Alerts",
		})
	})

	// HTMX modem status - will be expanded (consider moving to v1 or a dedicated htmx group)
	// For now, keeping it separate if it's purely for HTMX direct consumption and not a general API
	// router.GET("/api/modems/status", func(c *gin.Context) {
	// 	// In a real app, fetch data from pkg/modem or pkg/database
	// 	// For now, return a simple HTML snippet
	// 	c.Header("Content-Type", "text/html")
	// 	// Updated to use a dynamic refresh button that re-targets itself within the swapped content
	// 	c.String(http.StatusOK, "<div id=\"modem-status-dynamic\"><p>Modem 1: Online, Signal: Good</p><p>Modem 2: Offline</p><p><small>Status as of: "+time.Now().Format(time.RFC1123)+"</small></p><button hx-get=\"/api/modems/status\" hx-target=\"#modem-status-dynamic\" hx-swap=\"innerHTML\">Refresh Again</button></div>")
	// })

	// Add missing API routes that HTMX templates are calling
	router.GET("/api/sims", func(c *gin.Context) {
		ctx := context.Background()
		
		sims, err := simCardRepo.GetAllSIMCards(ctx)
		if err != nil {
			logging.Logger.Errorf("Error fetching SIMs: %v", err)
			c.Header("Content-Type", "text/html")
			c.String(http.StatusOK, `<div class="text-center text-red-600 dark:text-red-400 p-4">Error loading SIMs</div>`)
			return
		}
		
		html := ""
		for _, sim := range sims {
			statusColor := "green"
			statusBg := "bg-green-100 dark:bg-green-900"
			if sim.Status == "blocked" {
				statusColor = "red"
				statusBg = "bg-red-100 dark:bg-red-900"
			} else if sim.Status == "low_credit" {
				statusColor = "yellow"
				statusBg = "bg-yellow-100 dark:bg-yellow-900"
			}
			
			html += fmt.Sprintf(`
			<div class="flex items-center justify-between p-3 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-lg">
				<div class="flex items-center space-x-4 flex-1">
					<div class="flex-shrink-0">
						<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium %s text-%s-800 dark:text-%s-100">
							%s
						</span>
					</div>
					<div class="flex-1">
						<p class="text-sm font-medium text-gray-900 dark:text-white">%s</p>
						<p class="text-xs text-gray-500 dark:text-gray-400">%s</p>
					</div>
				</div>
			</div>`, 
			statusBg, statusColor, statusColor, strings.ToUpper(sim.Status),
			sim.MSISDN.String, sim.OperatorName.String)
		}
		
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})

	// Public routes (no auth required)
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login_standalone.tmpl", gin.H{
			"title": "Login - E173 Gateway",
		})
	})

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		
		// Get client info
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()
		
		// Attempt login using auth service
		_, session, err := authService.Login(username, password, ipAddress, userAgent)
		if err != nil {
			c.HTML(http.StatusOK, "login_standalone.tmpl", gin.H{
				"title": "Login",
				"error": "Invalid username or password",
			})
			return
		}
		
		// Set session cookie
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("session_token", session.SessionToken, int(time.Until(session.ExpiresAt).Seconds()), "/", "", false, true)
		
		// For HTMX requests, redirect to dashboard
		if c.GetHeader("HX-Request") == "true" {
			c.Header("HX-Redirect", "/")
			c.Status(http.StatusOK)
			return
		}
		
		// For regular requests, redirect to dashboard
		c.Redirect(http.StatusFound, "/")
	})

	// Logout route to clear token
	router.POST("/logout", func(c *gin.Context) {
		c.SetCookie("session_token", "", -1, "/", "", false, true)
		
		if c.GetHeader("HX-Request") == "true" {
			c.Header("HX-Redirect", "/login")
			c.Status(http.StatusOK)
			return
		}
		
		c.Redirect(http.StatusFound, "/login")
	})

	// Authentication API endpoints (public routes)
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", handlers.WrapHandler(authHandlers.Login))
		authGroup.POST("/logout", handlers.WrapHandler(authHandlers.Logout))
		authGroup.GET("/profile", handlers.WrapHandler(authHandlers.AuthMiddleware(authHandlers.GetProfile)))
		authGroup.POST("/change-password", handlers.WrapHandler(authHandlers.AuthMiddleware(authHandlers.ChangePassword)))
		authGroup.GET("/sessions", handlers.WrapHandler(authHandlers.AuthMiddleware(authHandlers.GetSessions)))
		authGroup.POST("/revoke-sessions", handlers.WrapHandler(authHandlers.AuthMiddleware(authHandlers.RevokeSessions)))
	}

	// Customer API endpoints (protected routes)
	customerGroup := router.Group("/api/customers")
	customerGroup.Use(handlers.WrapMiddleware(authHandlers.AuthMiddleware))
	{
		customerGroup.GET("", handlers.WrapHandler(customerHandlers.ListCustomers))
		customerGroup.POST("", handlers.WrapHandler(customerHandlers.CreateCustomer))
		customerGroup.GET("/:id", handlers.WrapHandler(customerHandlers.GetCustomer))
		customerGroup.PUT("/:id", handlers.WrapHandler(customerHandlers.UpdateCustomer))
		customerGroup.DELETE("/:id", handlers.WrapHandler(customerHandlers.DeleteCustomer))
		customerGroup.GET("/:id/balance", handlers.WrapHandler(customerHandlers.GetBalance))
		customerGroup.POST("/:id/balance", handlers.WrapHandler(customerHandlers.UpdateBalance))
		customerGroup.POST("/:id/suspend", handlers.WrapHandler(customerHandlers.SuspendCustomer))
		customerGroup.POST("/:id/reactivate", handlers.WrapHandler(customerHandlers.ReactivateCustomer))
		customerGroup.GET("/stats", handlers.WrapHandler(customerHandlers.GetCustomerStats))
	}

	// Admin template routes
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin/login.html", gin.H{
				"title": "Admin Login",
			})
		})
		
		// Protected admin routes
		adminProtected := adminGroup.Group("/")
		adminProtected.Use(handlers.WrapMiddleware(authHandlers.AuthMiddleware))
		adminProtected.Use(handlers.WrapRoleMiddleware(authHandlers.RoleMiddleware, "admin"))
		{
			adminProtected.GET("/dashboard", func(c *gin.Context) {
				c.HTML(http.StatusOK, "admin/dashboard.html", gin.H{
					"title": "Admin Dashboard",
				})
			})
		}
	}

	// Customer Management Frontend Routes (temporarily unprotected for testing)
	customerUIGroup := router.Group("/customers")
	// customerUIGroup.Use(handlers.WrapMiddleware(authHandlers.AuthMiddleware)) // Temporarily disabled
	{
		// Customer List
		customerUIGroup.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "customers_standalone.tmpl", gin.H{
				"title": "Customer Management",
			})
		})
		
		// Create Customer
		customerUIGroup.GET("/create", func(c *gin.Context) {
			c.HTML(http.StatusOK, "customers/create.html", gin.H{
				"title": "Create Customer",
			})
		})
		
		// Edit Customer (requires customer data)
		customerUIGroup.GET("/:id/edit", func(c *gin.Context) {
			customerID := c.Param("id")
			// In a real implementation, you'd fetch customer data here
			
			c.HTML(http.StatusOK, "customers/edit.html", gin.H{
				"title": "Edit Customer",
				"customer_id": customerID,
			})
		})
		
		// Balance Management (requires customer data)
		customerUIGroup.GET("/:id/balance", func(c *gin.Context) {
			customerID := c.Param("id")
			// In a real implementation, you'd fetch customer data here
			
			c.HTML(http.StatusOK, "customers/balance.html", gin.H{
				"title": "Balance Management",
				"customer_id": customerID,
			})
		})
	}

	// Gateway Management Frontend Routes (protected)
	gatewayUIGroup := router.Group("/gateways")
	gatewayUIGroup.Use(handlers.WrapMiddleware(authHandlers.AuthMiddleware))
	{
		// Gateway List
		gatewayUIGroup.GET("", gatewayHandler.GetGatewayListUI)
		
		// Create Gateway
		gatewayUIGroup.GET("/create", gatewayHandler.GetGatewayCreateUI)
		
		// Edit Gateway
		gatewayUIGroup.GET("/:id/edit", gatewayHandler.GetGatewayEditUI)
	}

	// Modem Management Frontend Routes (temporarily unprotected for testing)
	modemUIGroup := router.Group("/modems")
	// modemUIGroup.Use(handlers.WrapMiddleware(authHandlers.AuthMiddleware)) // Temporarily disabled for testing
	{
		// Modem List
		modemUIGroup.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "modems_standalone.tmpl", gin.H{
				"title": "Modem Management",
			})
		})
	}

	// SIM Management Frontend Routes (temporarily unprotected for testing)
	simUIGroup := router.Group("/sims")
	// simUIGroup.Use(handlers.WrapMiddleware(authHandlers.AuthMiddleware)) // Temporarily disabled for testing
	{
		// SIM List
		simUIGroup.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "sims_standalone.tmpl", gin.H{
				"title": "SIM Card Management",
			})
		})
	}

	// CDR Frontend Routes (temporarily unprotected for testing)
	cdrUIGroup := router.Group("/cdrs")
	// cdrUIGroup.Use(handlers.WrapMiddleware(authHandlers.AuthMiddleware)) // Temporarily disabled for testing
	{
		// CDR List
		cdrUIGroup.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "cdrs/list.tmpl", gin.H{
				"title": "Call Detail Records",
			})
		})
	}

	// Blacklist Frontend Routes (temporarily unprotected for testing)
	blacklistUIGroup := router.Group("/blacklist")
	// blacklistUIGroup.Use(handlers.WrapMiddleware(authHandlers.AuthMiddleware)) // Temporarily disabled for testing
	{
		// Blacklist List
		blacklistUIGroup.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "blacklist/list.tmpl", gin.H{
				"title": "Blacklist Management",
			})
		})
	}

	// Settings Frontend Routes (protected, admin only)
	settingsUIGroup := router.Group("/settings")
	settingsUIGroup.Use(auth.JWTMiddleware(jwtService))
	// settingsUIGroup.Use(handlers.WrapRoleMiddleware(authHandlers.RoleMiddleware, "admin")) // Temporarily disabled for testing
	{
		// Settings Page
		settingsUIGroup.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "settings_standalone.tmpl", gin.H{
				"title": "System Settings",
			})
		})
	}

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	logging.Logger.Infof("Starting server on %s in %s mode", serverAddr, cfg.GinMode)
	if err := router.Run(serverAddr); err != nil {
		logging.Logger.Fatalf("Failed to run server: %v", err)
	}
}
