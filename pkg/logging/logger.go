package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a global instance of the Logrus logger.
var Logger = logrus.New()

// InitLogger initializes the global logger with specific settings.
// It takes logLevel (e.g., "debug", "info", "warn", "error") and format ("json" or "text").
func InitLogger(logLevel string, format string) {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Logger.Warnf("Invalid log level '%s', defaulting to 'info': %v", logLevel, err)
		lvl = logrus.InfoLevel
	}
	Logger.SetLevel(lvl)

	Logger.SetOutput(os.Stdout) // Or a file, etc.

	switch format {
	case "json":
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00", // ISO8601
		})
	case "text":
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true, // Optional: for colored output in TTY
		})
	default:
		Logger.Warnf("Invalid log format '%s', defaulting to 'text'", format)
		Logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	Logger.Info("Logger initialized")
}

// Example of adding a field to all log entries from this package
// func GetLoggerWithField(key string, value interface{}) *logrus.Entry {
// 	return Logger.WithField(key, value)
// }
