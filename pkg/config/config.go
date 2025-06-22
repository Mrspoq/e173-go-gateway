package config

import (
	"os"
	"strconv"

	"github.com/e173-gateway/e173_go_gateway/pkg/logging" // Import structured logger
)

// AppConfig holds the application configuration.
// For simplicity, we'll use environment variables for now.
// Later, we can integrate Viper for more robust config loading (e.g., from YAML files).
type AppConfig struct {
	ServerPort     string
	DatabaseURL    string // e.g., "postgres://user:password@host:port/dbname?sslmode=disable"
	AsteriskAMIHost string
	AsteriskAMIPort string
	AsteriskAMIUser string
	AsteriskAMIPass string
	GinMode        string // "debug" or "release"
	LogLevel       string // e.g., "debug", "info", "warn", "error"
	LogFormat      string // "json" or "text"
}

// LoadConfig loads configuration from environment variables or defaults.
func LoadConfig() *AppConfig {
	cfg := &AppConfig{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://e173_user:e173_pass@localhost:5432/e173_gateway?sslmode=disable"),
		AsteriskAMIHost: getEnv("ASTERISK_AMI_HOST", "localhost"),
		AsteriskAMIPort: getEnv("ASTERISK_AMI_PORT", "5038"),
		AsteriskAMIUser: getEnv("ASTERISK_AMI_USER", "admin"),
		AsteriskAMIPass: getEnv("ASTERISK_AMI_PASS", "adminpass"), // Example, use secrets management in production
		GinMode:        getEnv("GIN_MODE", "debug"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogFormat:      getEnv("LOG_FORMAT", "text"),
	}

	// Initialize logger early if its config is available, or use a temp logger
	// For now, we assume logger is initialized in main after config is loaded.
	// If logging is needed during config load itself, a more complex setup might be needed.
	logging.Logger.Infof("Configuration partially loaded (sensitive fields may be redacted): %+v", scrubbedConfigForLog(cfg))

	return cfg
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an int or returns a default value.
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// scrubbedConfigForLog returns a copy of the config with sensitive fields redacted for logging.
func scrubbedConfigForLog(cfg *AppConfig) AppConfig {
	safeCfg := *cfg
	if safeCfg.DatabaseURL != "" { // Basic redaction, can be improved
		safeCfg.DatabaseURL = "postgres://user:****@..."
	}
	if safeCfg.AsteriskAMIPass != "" {
		safeCfg.AsteriskAMIPass = "****"
	}
	return safeCfg
}
