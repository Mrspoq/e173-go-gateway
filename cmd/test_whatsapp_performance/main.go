package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	
	"github.com/e173-gateway/e173_go_gateway/pkg/validation"
)

type MoroccoData struct {
	Operators map[string]struct {
		Name     string   `json:"name"`
		Prefixes []string `json:"prefixes"`
	} `json:"operators"`
}

type TestResult struct {
	Number       string
	HasWhatsApp  bool
	ResponseTime time.Duration
	Error        error
	Timestamp    time.Time
}

func main() {
	fmt.Println("Morocco WhatsApp API Performance Test")
	fmt.Println("====================================")
	fmt.Println("Testing: 5 numbers/second for 20 seconds = 100 numbers total\n")
	
	// Get API key
	apiKey := os.Getenv("WHATSAPP_API_KEY")
	if apiKey == "" {
		log.Fatal("WHATSAPP_API_KEY not set")
	}
	
	// Load prefixes
	data, err := ioutil.ReadFile("/root/e173_go_gateway/data/morocco_mobile_prefixes_correct.json")
	if err != nil {
		log.Fatalf("Failed to read Morocco data: %v", err)
	}
	
	var moroccoData MoroccoData
	if err := json.Unmarshal(data, &moroccoData); err != nil {
		log.Fatalf("Failed to parse Morocco data: %v", err)
	}
	
	// Collect all prefixes
	var allPrefixes []string
	for _, op := range moroccoData.Operators {
		allPrefixes = append(allPrefixes, op.Prefixes...)
	}
	
	// Initialize validators
	phoneValidator, err := validation.NewLibPhoneNumberValidator()
	if err != nil {
		log.Fatalf("Failed to initialize libphonenumber: %v", err)
	}
	
	whatsappValidator := validation.NewPrivateWhatsAppValidator(apiKey)
	
	// Generate test numbers
	rand.Seed(time.Now().UnixNano())
	testNumbers := make([]string, 100)
	
	fmt.Println("Generating 100 valid Morocco mobile numbers...")
	for i := 0; i < 100; i++ {
		// Pick random prefix
		prefix := allPrefixes[rand.Intn(len(allPrefixes))]
		
		// Generate 6-digit subscriber number
		subscriber := fmt.Sprintf("%06d", 100000+rand.Intn(900000))
		fullNumber := fmt.Sprintf("+%s%s", prefix, subscriber)
		
		// Verify it's valid
		info, err := phoneValidator.ValidatePhoneNumber(fullNumber, "MA")
		if err != nil || !info.IsValid {
			// Retry if invalid
			i--
			continue
		}
		
		testNumbers[i] = fullNumber
	}
	
	fmt.Println("✅ Generated 100 valid Morocco numbers")
	fmt.Println("\nStarting performance test...")
	fmt.Println("=============================\n")
	
	// Results storage
	results := make([]TestResult, 0, 100)
	var resultsMutex sync.Mutex
	
	// Statistics
	var stats struct {
		sync.Mutex
		Processed    int
		HasWhatsApp  int
		NoWhatsApp   int
		Errors       int
		TotalLatency time.Duration
		MinLatency   time.Duration
		MaxLatency   time.Duration
	}
	stats.MinLatency = time.Hour // Initialize to large value
	
	// Start time
	testStart := time.Now()
	
	// Process 5 numbers per second
	ticker := time.NewTicker(200 * time.Millisecond) // 200ms = 5 per second
	defer ticker.Stop()
	
	numberIndex := 0
	
	for {
		select {
		case <-ticker.C:
			if numberIndex >= len(testNumbers) {
				goto done
			}
			
			// Process number
			number := testNumbers[numberIndex]
			numberIndex++
			
			go func(num string, idx int) {
				// Measure API response time
				start := time.Now()
				result, err := whatsappValidator.ValidateNumber(num)
				latency := time.Since(start)
				
				// Store result
				testResult := TestResult{
					Number:       num,
					ResponseTime: latency,
					Timestamp:    start,
					Error:        err,
				}
				
				if err == nil && result != nil {
					testResult.HasWhatsApp = result.HasWhatsApp
				}
				
				resultsMutex.Lock()
				results = append(results, testResult)
				resultsMutex.Unlock()
				
				// Update stats
				stats.Lock()
				stats.Processed++
				stats.TotalLatency += latency
				
				if latency < stats.MinLatency {
					stats.MinLatency = latency
				}
				if latency > stats.MaxLatency {
					stats.MaxLatency = latency
				}
				
				if err != nil {
					stats.Errors++
					fmt.Printf("[%03d] %s: ❌ Error (%v) - %v\n", idx+1, num, latency, err)
				} else if result.HasWhatsApp {
					stats.HasWhatsApp++
					fmt.Printf("[%03d] %s: ✅ Has WhatsApp - %v\n", idx+1, num, latency)
				} else {
					stats.NoWhatsApp++
					fmt.Printf("[%03d] %s: ❌ No WhatsApp - %v\n", idx+1, num, latency)
				}
				
				// Progress update every 10 numbers
				if stats.Processed%10 == 0 {
					fmt.Printf("\n>>> Progress: %d/100 processed (%.1f%%) <<<\n\n", 
						stats.Processed, float64(stats.Processed))
				}
				stats.Unlock()
			}(number, numberIndex-1)
		}
	}
	
done:
	// Wait for all goroutines to complete
	for stats.Processed < 100 {
		time.Sleep(100 * time.Millisecond)
	}
	
	totalTime := time.Since(testStart)
	
	// Calculate statistics
	avgLatency := stats.TotalLatency / time.Duration(stats.Processed)
	actualTPS := float64(stats.Processed) / totalTime.Seconds()
	
	// Final Report
	fmt.Println("\n\n" + strings.Repeat("=", 60))
	fmt.Println("PERFORMANCE TEST COMPLETE")
	fmt.Println(strings.Repeat("=", 60))
	
	fmt.Printf("\nTest Duration: %.2f seconds\n", totalTime.Seconds())
	fmt.Printf("Numbers Tested: %d\n", stats.Processed)
	fmt.Printf("Target Rate: 5 TPS\n")
	fmt.Printf("Actual Rate: %.2f TPS\n", actualTPS)
	
	fmt.Println("\nRESULTS:")
	fmt.Printf("  ✅ Has WhatsApp: %d (%.1f%%)\n", 
		stats.HasWhatsApp, float64(stats.HasWhatsApp)*100/float64(stats.Processed))
	fmt.Printf("  ❌ No WhatsApp: %d (%.1f%%)\n", 
		stats.NoWhatsApp, float64(stats.NoWhatsApp)*100/float64(stats.Processed))
	fmt.Printf("  ⚠️  API Errors: %d\n", stats.Errors)
	
	fmt.Println("\nLATENCY STATISTICS:")
	fmt.Printf("  Average: %v\n", avgLatency)
	fmt.Printf("  Minimum: %v\n", stats.MinLatency)
	fmt.Printf("  Maximum: %v\n", stats.MaxLatency)
	
	// Latency distribution
	var under100ms, under250ms, under500ms, over500ms int
	for _, r := range results {
		switch {
		case r.ResponseTime < 100*time.Millisecond:
			under100ms++
		case r.ResponseTime < 250*time.Millisecond:
			under250ms++
		case r.ResponseTime < 500*time.Millisecond:
			under500ms++
		default:
			over500ms++
		}
	}
	
	fmt.Println("\nLATENCY DISTRIBUTION:")
	fmt.Printf("  < 100ms: %d (%.1f%%)\n", under100ms, float64(under100ms)*100/float64(len(results)))
	fmt.Printf("  100-250ms: %d (%.1f%%)\n", under250ms, float64(under250ms)*100/float64(len(results)))
	fmt.Printf("  250-500ms: %d (%.1f%%)\n", under500ms, float64(under500ms)*100/float64(len(results)))
	fmt.Printf("  > 500ms: %d (%.1f%%)\n", over500ms, float64(over500ms)*100/float64(len(results)))
	
	fmt.Println("\nCONCLUSIONS:")
	fmt.Printf("✅ API handled %.2f TPS successfully\n", actualTPS)
	fmt.Printf("✅ Average response time: %v\n", avgLatency)
	fmt.Printf("✅ %d%% of requests completed under 250ms\n", 
		(under100ms+under250ms)*100/len(results))
	
	if stats.Errors > 0 {
		fmt.Printf("⚠️  %d errors occurred (%.1f%% error rate)\n", 
			stats.Errors, float64(stats.Errors)*100/float64(stats.Processed))
	}
}