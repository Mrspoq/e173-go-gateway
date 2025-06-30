package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
	
	"github.com/e173-gateway/e173_go_gateway/pkg/validation"
)

type MoroccoData struct {
	Operators map[string]struct {
		Name     string   `json:"name"`
		Prefixes []string `json:"prefixes"`
	} `json:"operators"`
}

// Simple rate limiter
type RateLimiter struct {
	rate     int
	lastCall time.Time
}

func NewRateLimiter(ratePerSecond int) *RateLimiter {
	return &RateLimiter{
		rate:     ratePerSecond,
		lastCall: time.Now(),
	}
}

func (rl *RateLimiter) Wait() {
	elapsed := time.Since(rl.lastCall)
	minInterval := time.Second / time.Duration(rl.rate)
	if elapsed < minInterval {
		time.Sleep(minInterval - elapsed)
	}
	rl.lastCall = time.Now()
}

func main() {
	fmt.Println("Morocco WhatsApp API Test - Production Ready")
	fmt.Println("===========================================\n")
	
	// Get API key
	apiKey := os.Getenv("WHATSAPP_API_KEY")
	if apiKey == "" {
		log.Fatal("WHATSAPP_API_KEY not set. Run: export WHATSAPP_API_KEY=your-api-key-here")
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
	
	// Initialize validators
	phoneValidator, err := validation.NewLibPhoneNumberValidator()
	if err != nil {
		log.Fatalf("Failed to initialize libphonenumber: %v", err)
	}
	
	whatsappValidator := validation.NewPrivateWhatsAppValidator(apiKey)
	
	// Create rate limiter - 5 requests per second
	limiter := NewRateLimiter(5)
	
	// Test known numbers first
	fmt.Println("Testing known numbers:")
	fmt.Println("=====================")
	
	knownNumbers := []struct {
		number   string
		expected string
	}{
		{"+33761698939", "Has WhatsApp (French)"},
		{"+212761698939", "No WhatsApp (Morocco 761)"},
	}
	
	for _, test := range knownNumbers {
		limiter.Wait()
		fmt.Printf("\n%s: ", test.number)
		
		result, err := whatsappValidator.ValidateNumber(test.number)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else if result.HasWhatsApp {
			fmt.Printf("✅ Has WhatsApp")
		} else {
			fmt.Printf("❌ No WhatsApp")
		}
		fmt.Printf(" (Expected: %s)\n", test.expected)
	}
	
	// Test Morocco prefixes
	fmt.Println("\n\nTesting Morocco operator prefixes (3 per operator):")
	fmt.Println("==================================================")
	
	rand.Seed(time.Now().UnixNano())
	
	stats := struct {
		Total       int
		HasWhatsApp int
		NoWhatsApp  int
		Errors      int
		ByOperator  map[string]struct{ Total, HasWA int }
	}{
		ByOperator: make(map[string]struct{ Total, HasWA int }),
	}
	
	startTime := time.Now()
	
	for _, operator := range moroccoData.Operators {
		fmt.Printf("\n%s:\n", operator.Name)
		
		opStats := struct{ Total, HasWA int }{}
		
		// Test 3 numbers per operator
		for i := 0; i < 3 && i < len(operator.Prefixes); i++ {
			limiter.Wait()
			
			// Pick random prefix
			prefix := operator.Prefixes[rand.Intn(len(operator.Prefixes))]
			
			// Generate realistic number (not too random)
			areaCode := 600 + rand.Intn(100)  // 600-699
			subscriber := 1000 + rand.Intn(9000) // 1000-9999
			fullNumber := fmt.Sprintf("+%s%03d%04d", prefix, areaCode, subscriber)
			
			// Validate format first
			info, err := phoneValidator.ValidatePhoneNumber(fullNumber, "MA")
			if err != nil || !info.IsValid {
				fmt.Printf("  ❌ %s: Invalid format\n", fullNumber)
				continue
			}
			
			// Check WhatsApp
			fmt.Printf("  %s: ", fullNumber)
			
			result, err := whatsappValidator.ValidateNumber(fullNumber)
			
			stats.Total++
			opStats.Total++
			
			if err != nil {
				fmt.Printf("❌ API Error: %v\n", err)
				stats.Errors++
			} else if result.HasWhatsApp {
				fmt.Printf("✅ Has WhatsApp\n")
				stats.HasWhatsApp++
				opStats.HasWA++
			} else {
				fmt.Printf("❌ No WhatsApp\n")
				stats.NoWhatsApp++
			}
		}
		
		stats.ByOperator[operator.Name] = opStats
	}
	
	elapsed := time.Since(startTime)
	
	// Summary
	fmt.Println("\n\nSUMMARY")
	fmt.Println("=======")
	fmt.Printf("Total Morocco numbers tested: %d\n", stats.Total)
	fmt.Printf("Have WhatsApp: %d (%.1f%%)\n", 
		stats.HasWhatsApp, float64(stats.HasWhatsApp)*100/float64(max(stats.Total, 1)))
	fmt.Printf("No WhatsApp: %d (%.1f%%)\n", 
		stats.NoWhatsApp, float64(stats.NoWhatsApp)*100/float64(max(stats.Total, 1)))
	fmt.Printf("API errors: %d\n", stats.Errors)
	
	fmt.Println("\nBy Operator:")
	for name, opStats := range stats.ByOperator {
		pct := float64(opStats.HasWA)*100/float64(max(opStats.Total, 1))
		fmt.Printf("  %s: %d/%d have WhatsApp (%.1f%%)\n",
			name, opStats.HasWA, opStats.Total, pct)
	}
	
	fmt.Printf("\nTime: %.2f seconds (%.2f requests/second)\n", 
		elapsed.Seconds(), float64(stats.Total+2)/elapsed.Seconds())
	
	fmt.Println("\n✅ Validation Pipeline Status:")
	fmt.Println("  1. Libphonenumber: Validates E.164 format")
	fmt.Println("  2. Prefix Database: Ensures valid operator (169 assigned prefixes)")
	fmt.Println("  3. WhatsApp API: Real-time account verification")
	fmt.Printf("  4. Rate Limiting: Respecting 5 TPS (actual: %.2f TPS)\n",
		float64(stats.Total+2)/elapsed.Seconds())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}