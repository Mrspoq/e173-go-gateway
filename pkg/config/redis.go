package config

import (
    "os"
    "strconv"
    "time"
    
    "github.com/e173-gateway/e173_go_gateway/pkg/cache"
)

// LoadRedisConfig loads Redis configuration from environment variables
func LoadRedisConfig() *cache.Config {
    config := cache.DefaultConfig()
    
    // Override with environment variables if present
    if host := os.Getenv("REDIS_HOST"); host != "" {
        config.Host = host
    }
    
    if port := os.Getenv("REDIS_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            config.Port = p
        }
    }
    
    if password := os.Getenv("REDIS_PASSWORD"); password != "" {
        config.Password = password
    }
    
    if db := os.Getenv("REDIS_DB"); db != "" {
        if d, err := strconv.Atoi(db); err == nil {
            config.DB = d
        }
    }
    
    if prefix := os.Getenv("REDIS_KEY_PREFIX"); prefix != "" {
        config.KeyPrefix = prefix
    }
    
    if ttl := os.Getenv("REDIS_DEFAULT_TTL"); ttl != "" {
        if duration, err := time.ParseDuration(ttl); err == nil {
            config.DefaultTTL = duration
        }
    }
    
    return config
}