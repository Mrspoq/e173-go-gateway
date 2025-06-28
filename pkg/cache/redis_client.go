package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/go-redis/redis/v8"
)

// RedisClient wraps the Redis client with application-specific methods
type RedisClient struct {
    client      *redis.Client
    defaultTTL  time.Duration
    keyPrefix   string
}

// Config holds Redis configuration
type Config struct {
    Host         string
    Port         int
    Password     string
    DB           int
    DefaultTTL   time.Duration
    KeyPrefix    string
    MaxRetries   int
    DialTimeout  time.Duration
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

// DefaultConfig returns default Redis configuration
func DefaultConfig() *Config {
    return &Config{
        Host:         "localhost",
        Port:         6379,
        Password:     "",
        DB:           0,
        DefaultTTL:   24 * time.Hour,
        KeyPrefix:    "e173:",
        MaxRetries:   3,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    }
}

// NewRedisClient creates a new Redis client with the given configuration
func NewRedisClient(cfg *Config) (*RedisClient, error) {
    client := redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password:     cfg.Password,
        DB:           cfg.DB,
        MaxRetries:   cfg.MaxRetries,
        DialTimeout:  cfg.DialTimeout,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
    })
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }
    
    return &RedisClient{
        client:     client,
        defaultTTL: cfg.DefaultTTL,
        keyPrefix:  cfg.KeyPrefix,
    }, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
    return r.client.Close()
}

// buildKey creates a namespaced key
func (r *RedisClient) buildKey(key string) string {
    return r.keyPrefix + key
}

// Set stores a value with the default TTL
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
    return r.SetWithTTL(ctx, key, value, r.defaultTTL)
}

// SetWithTTL stores a value with a custom TTL
func (r *RedisClient) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to marshal value: %w", err)
    }
    
    return r.client.Set(ctx, r.buildKey(key), data, ttl).Err()
}

// Get retrieves a value from cache
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := r.client.Get(ctx, r.buildKey(key)).Bytes()
    if err != nil {
        if err == redis.Nil {
            return ErrCacheMiss
        }
        return fmt.Errorf("failed to get value: %w", err)
    }
    
    if err := json.Unmarshal(data, dest); err != nil {
        return fmt.Errorf("failed to unmarshal value: %w", err)
    }
    
    return nil
}

// Exists checks if a key exists
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
    count, err := r.client.Exists(ctx, r.buildKey(key)).Result()
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

// Delete removes a key from cache
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
    fullKeys := make([]string, len(keys))
    for i, key := range keys {
        fullKeys[i] = r.buildKey(key)
    }
    
    return r.client.Del(ctx, fullKeys...).Err()
}

// Increment increments a counter
func (r *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
    return r.client.Incr(ctx, r.buildKey(key)).Result()
}

// IncrementBy increments a counter by a specific amount
func (r *RedisClient) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
    return r.client.IncrBy(ctx, r.buildKey(key), value).Result()
}

// SetNX sets a value only if it doesn't exist (useful for locks)
func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
    data, err := json.Marshal(value)
    if err != nil {
        return false, fmt.Errorf("failed to marshal value: %w", err)
    }
    
    return r.client.SetNX(ctx, r.buildKey(key), data, ttl).Result()
}

// GetMultiple retrieves multiple values at once
func (r *RedisClient) GetMultiple(ctx context.Context, keys []string) (map[string][]byte, error) {
    fullKeys := make([]string, len(keys))
    for i, key := range keys {
        fullKeys[i] = r.buildKey(key)
    }
    
    values, err := r.client.MGet(ctx, fullKeys...).Result()
    if err != nil {
        return nil, err
    }
    
    result := make(map[string][]byte)
    for i, value := range values {
        if value != nil {
            if data, ok := value.(string); ok {
                result[keys[i]] = []byte(data)
            }
        }
    }
    
    return result, nil
}

// SetMultiple sets multiple values at once
func (r *RedisClient) SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
    pipe := r.client.Pipeline()
    
    for key, value := range items {
        data, err := json.Marshal(value)
        if err != nil {
            return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
        }
        pipe.Set(ctx, r.buildKey(key), data, ttl)
    }
    
    _, err := pipe.Exec(ctx)
    return err
}

// GetTTL gets the remaining TTL for a key
func (r *RedisClient) GetTTL(ctx context.Context, key string) (time.Duration, error) {
    return r.client.TTL(ctx, r.buildKey(key)).Result()
}

// Expire sets a new expiration time for a key
func (r *RedisClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
    return r.client.Expire(ctx, r.buildKey(key), ttl).Err()
}

// Flush removes all keys with the configured prefix
func (r *RedisClient) Flush(ctx context.Context) error {
    iter := r.client.Scan(ctx, 0, r.keyPrefix+"*", 0).Iterator()
    var keys []string
    
    for iter.Next(ctx) {
        keys = append(keys, iter.Val())
    }
    
    if err := iter.Err(); err != nil {
        return err
    }
    
    if len(keys) > 0 {
        return r.client.Del(ctx, keys...).Err()
    }
    
    return nil
}

// HealthCheck verifies Redis connectivity
func (r *RedisClient) HealthCheck(ctx context.Context) error {
    return r.client.Ping(ctx).Err()
}

// GetClient returns the underlying Redis client for advanced operations
func (r *RedisClient) GetClient() *redis.Client {
    return r.client
}

// ErrCacheMiss is returned when a key is not found in cache
var ErrCacheMiss = fmt.Errorf("cache miss")