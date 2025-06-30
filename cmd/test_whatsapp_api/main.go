package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	
	"github.com/e173-gateway/e173_go_gateway/pkg/validation"
)

type MoroccoData struct {
	Operators map[string]struct {
		Name     string   `json:"name"`
		Prefixes []string `json:"prefixes"`
	} `json:"operators"`
}

func main() {
	fmt.Println("Morocco WhatsApp Validation Test")
	fmt.Println("================================\n")
	
	// Try to get API key from multiple sources
	apiKey := os.Getenv("WHATSAPP_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("WA_VALIDATOR_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("PRIVATE_WHATSAPP_KEY")
	}
	
	if apiKey == "" {
		fmt.Println("WhatsApp API key not found in environment variables.")
		fmt.Println("Checked: WHATSAPP_API_KEY, WA_VALIDATOR_API_KEY, PRIVATE_WHATSAPP_KEY")
		fmt.Print("\nEnter your wa-validator.xyz API key (or press Enter to skip): ")
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(input)
		
		if apiKey == "" {
			fmt.Println("\nNo API key provided. Exiting...")
			return
		}
	} else {
		fmt.Println("✅ API key found in environment variables")
	}
	
	// Load Morocco prefixes
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
	
	// Use the private WhatsApp validator
	whatsappValidator := validation.NewPrivateWhatsAppValidator(apiKey)
	
	fmt.Println("\nUsing Private WhatsApp API: https://bulkvalidation.wa-validator.xyz/v2/validate/wa_id")
	fmt.Println("\nTesting Morocco numbers by operator:")
	fmt.Println("=====================================\n")
	
	// Test statistics
	stats := struct {
		TotalTested     int
		ValidNumbers    int
		HasWhatsApp     int
		NoWhatsApp      int
		APIErrors       int
	}{}
	
	rand.Seed(time.Now().UnixNano())
	
	// Test 2 numbers from each operator
	for _, operator := range moroccoData.Operators {
		fmt.Printf("\n%s:\n", operator.Name)
		fmt.Println(strings.Repeat("-", 40))
		
		// Test 2 random prefixes
		for i := 0; i < 2 && i < len(operator.Prefixes); i++ {
			prefix := operator.Prefixes[rand.Intn(len(operator.Prefixes))]
			
			// Generate realistic number
			subscriber := fmt.Sprintf("%06d", 200000+rand.Intn(800000))
			fullNumber := fmt.Sprintf("+%s%s", prefix, subscriber)
			
			stats.TotalTested++
			
			// First validate with libphonenumber
			info, err := phoneValidator.ValidatePhoneNumber(fullNumber, "MA")
			if err != nil || !info.IsValid {
				fmt.Printf("  %s: ❌ Invalid number format\n", fullNumber)
				continue
			}
			
			stats.ValidNumbers++
			
			// Test with WhatsApp API
			fmt.Printf("  %s: ", fullNumber)
			
			result, err := whatsappValidator.ValidateNumber(fullNumber)
			if err != nil {
				fmt.Printf("❌ API Error: %v\n", err)
				stats.APIErrors++
			} else if result.HasWhatsApp {
				fmt.Printf("✅ Has WhatsApp (confidence: %.2f)\n", result.Confidence)
				stats.HasWhatsApp++
			} else {
				fmt.Printf("❌ No WhatsApp\n")
				stats.NoWhatsApp++
			}
			
			// Small delay to respect rate limits
			time.Sleep(500 * time.Millisecond)
		}
	}
	
	// Test some specific numbers if you have them
	fmt.Println("\n\nTest specific numbers (optional):")
	fmt.Println("=================================")
	fmt.Println("Enter phone numbers to test (one per line, empty line to finish):")
	
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		number := strings.TrimSpace(input)
		
		if number == "" {
			break
		}
		
		// Ensure it starts with +
		if !strings.HasPrefix(number, "+") {
			number = "+" + number
		}
		
		// Validate format first
		info, err := phoneValidator.ValidatePhoneNumber(number, "MA")
		if err != nil || !info.IsValid {
			fmt.Printf("  ❌ Invalid number format: %s\n", number)
			continue
		}
		
		// Check WhatsApp
		result, err := whatsappValidator.ValidateNumber(number)
		if err != nil {
			fmt.Printf("  ❌ API Error: %v\n", err)
		} else if result.HasWhatsApp {
			fmt.Printf("  ✅ %s has WhatsApp (confidence: %.2f)\n", number, result.Confidence)
		} else {
			fmt.Printf("  ❌ %s does not have WhatsApp\n", number)
		}
		
		time.Sleep(500 * time.Millisecond)
	}
	
	// Summary
	fmt.Println("\n\nSUMMARY")
	fmt.Println("=======")
	fmt.Printf("Total numbers tested: %d\n", stats.TotalTested)
	fmt.Printf("Valid Morocco numbers: %d\n", stats.ValidNumbers)
	fmt.Printf("Have WhatsApp: %d (%.1f%%)\n", 
		stats.HasWhatsApp, float64(stats.HasWhatsApp)*100/float64(max(stats.ValidNumbers, 1)))
	fmt.Printf("No WhatsApp: %d (%.1f%%)\n", 
		stats.NoWhatsApp, float64(stats.NoWhatsApp)*100/float64(max(stats.ValidNumbers, 1)))
	fmt.Printf("API errors: %d\n", stats.APIErrors)
	
	fmt.Println("\nValidation Pipeline Summary:")
	fmt.Println("1. ✅ Libphonenumber validates format and length")
	fmt.Println("2. ✅ Our prefix database ensures valid operator assignment") 
	fmt.Println("3. ✅ WhatsApp API confirms active WhatsApp account")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}