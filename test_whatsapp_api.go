package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/joho/godotenv"
    "github.com/e173-gateway/e173_go_gateway/pkg/validation"
    "github.com/e173-gateway/e173_go_gateway/pkg/database"
    "github.com/e173-gateway/e173_go_gateway/pkg/repository"
    "github.com/e173-gateway/e173_go_gateway/pkg/config"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Test phone numbers
    testNumbers := []string{
        "+2348123456789", // Nigeria MTN
        "+2347012345678", // Nigeria Airtel
        "+12125551234",   // US number
        "+34674944456",   // Spain (example from API docs)
    }
    
    if len(os.Args) > 1 {
        testNumbers = os.Args[1:]
    }
    
    fmt.Println("=== Testing WhatsApp API Integration ===")
    
    // Test 1: Direct API call (no caching)
    fmt.Println("\n1. Testing direct API calls:")
    apiKey := "e42f7c9b-2a8e-4b86-a7e4-8f1de2c01f53"
    validator := validation.NewPrivateWhatsAppValidator(apiKey)
    
    for _, number := range testNumbers {
        result, err := validator.ValidateNumber(number)
        if err != nil {
            fmt.Printf("❌ Error validating %s: %v\n", number, err)
            continue
        }
        
        if result.HasWhatsApp {
            fmt.Printf("✅ %s HAS WhatsApp (confidence: %.2f)\n", number, result.Confidence)
        } else {
            fmt.Printf("❌ %s NO WhatsApp (confidence: %.2f)\n", number, result.Confidence)
        }
    }
    
    // Test 2: Database-cached API calls
    fmt.Println("\n2. Testing database-cached API calls:")
    
    // Load config and connect to database
    cfg := config.LoadConfig()
    dbPool, err := database.NewDBPool(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer dbPool.Close()
    
    // Create repository and DB-backed validator
    cacheRepo := repository.NewSimpleWhatsAppValidationRepository(dbPool)
    dbValidator := validation.NewPrivateWhatsAppValidatorDB(apiKey, cacheRepo)
    
    // Test same numbers with caching
    for _, number := range testNumbers {
        result, err := dbValidator.ValidateNumber(number)
        if err != nil {
            fmt.Printf("❌ Error validating %s: %v\n", number, err)
            continue
        }
        
        if result.HasWhatsApp {
            fmt.Printf("✅ %s HAS WhatsApp (confidence: %.2f, source: %s)\n", 
                number, result.Confidence, result.Source)
        } else {
            fmt.Printf("❌ %s NO WhatsApp (confidence: %.2f, source: %s)\n", 
                number, result.Confidence, result.Source)
        }
    }
    
    // Test 3: Check cache stats
    fmt.Println("\n3. Cache statistics:")
    stats := dbValidator.GetCacheStats()
    for key, value := range stats {
        fmt.Printf("   %s: %v\n", key, value)
    }
    
    // Test 4: Real person detection
    fmt.Println("\n4. Testing real person detection:")
    for _, number := range testNumbers {
        isReal, confidence, err := dbValidator.IsLikelyRealPerson(number)
        if err != nil {
            fmt.Printf("❌ Error checking %s: %v\n", number, err)
            continue
        }
        
        if isReal {
            fmt.Printf("👤 %s is LIKELY a real person (confidence: %.2f)\n", number, confidence)
        } else {
            fmt.Printf("🤖 %s is LIKELY automated/spam (confidence: %.2f)\n", number, confidence)
        }
    }
    
    fmt.Println("\n✅ WhatsApp API integration test complete!")
}