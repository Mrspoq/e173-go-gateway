package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

func main() {
    // API credentials
    apiKey := "e42f7c9b-2a8e-4b86-a7e4-8f1de2c01f53"
    baseURL := "https://bulkvalidation.wa-validator.xyz/v2/validate/wa_id"
    
    // Test numbers
    testNumbers := []string{"+34674944456"}
    if len(os.Args) > 1 {
        testNumbers = os.Args[1:]
    }
    
    fmt.Println("=== Direct WhatsApp API Test ===")
    fmt.Printf("API: %s\n", baseURL)
    fmt.Println("Testing numbers:")
    
    client := &http.Client{Timeout: 10 * time.Second}
    
    for _, number := range testNumbers {
        // Clean number (remove +)
        cleanNumber := number
        if len(cleanNumber) > 0 && cleanNumber[0] == '+' {
            cleanNumber = cleanNumber[1:]
        }
        
        // Build URL
        url := fmt.Sprintf("%s?number=%s", baseURL, cleanNumber)
        
        // Create request
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            fmt.Printf("❌ Error creating request for %s: %v\n", number, err)
            continue
        }
        
        // Set headers
        req.Header.Set("Authorization", "Bearer "+apiKey)
        req.Header.Set("User-Agent", "E173-Gateway/1.0")
        
        fmt.Printf("\nTesting %s...\n", number)
        
        // Make request
        resp, err := client.Do(req)
        if err != nil {
            fmt.Printf("❌ Request failed: %v\n", err)
            continue
        }
        defer resp.Body.Close()
        
        // Read response
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("❌ Failed to read response: %v\n", err)
            continue
        }
        
        fmt.Printf("Status Code: %d\n", resp.StatusCode)
        fmt.Printf("Raw Response: %s\n", string(body))
        
        // Parse response
        if resp.StatusCode == 200 {
            var result struct {
                Status   bool   `json:"status"`
                Valid    bool   `json:"valid"`
                WaID     string `json:"wa_id"`
                ChatLink string `json:"chat_link"`
            }
            
            if err := json.Unmarshal(body, &result); err != nil {
                fmt.Printf("❌ Failed to parse JSON: %v\n", err)
                continue
            }
            
            if result.Status && result.Valid {
                fmt.Printf("✅ %s HAS WhatsApp\n", number)
                fmt.Printf("   WhatsApp ID: %s\n", result.WaID)
                fmt.Printf("   Chat Link: %s\n", result.ChatLink)
            } else if result.Status && !result.Valid {
                fmt.Printf("❌ %s NO WhatsApp\n", number)
            } else {
                fmt.Printf("⚠️  %s API returned invalid status\n", number)
            }
        }
    }
    
    fmt.Println("\n✅ Test complete!")
}