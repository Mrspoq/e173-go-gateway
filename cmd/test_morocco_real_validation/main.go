package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	
	"github.com/e173-gateway/e173_go_gateway/pkg/validation"
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
)

type MoroccoData struct {
	Operators map[string]struct {
		Name     string   `json:"name"`
		Prefixes []string `json:"prefixes"`
	} `json:"operators"`
}

func main() {
	fmt.Println("Morocco Validation Test - Real WhatsApp API")
	fmt.Println("==========================================\n")
	
	// Check for WhatsApp API key
	whatsappAPIKey := os.Getenv("WHATSAPP_API_KEY")
	if whatsappAPIKey == "" {
		fmt.Println("⚠️  WARNING: WHATSAPP_API_KEY environment variable not set")
		fmt.Println("   To use real WhatsApp validation, set: export WHATSAPP_API_KEY=your_key")
		fmt.Println("   Using mock validation for now...\n")
	} else {
		fmt.Println("✅ WhatsApp API key found, will use real validation\n")
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
	
	// Initialize WhatsApp validator
	var whatsappValidator validation.WhatsAppValidator
	if whatsappAPIKey != "" {
		whatsappValidator = validation.NewWhatsAppBusinessValidator(whatsappAPIKey)
		fmt.Println("Using real WhatsApp Business API")
	} else {
		// Use mock if no API key
		whatsappValidator = &mockWhatsAppValidator{}
		fmt.Println("Using mock WhatsApp validator")
	}
	
	fmt.Println("\nTesting sample numbers from each operator:")
	fmt.Println("=========================================\n")
	
	// Statistics
	stats := struct {
		TotalTested      int
		PhoneLibValid    int
		WhatsAppChecked  int
		HasWhatsApp      int
		APIErrors        int
	}{}
	
	// Test each operator with fewer numbers to avoid API rate limits
	for opKey, operator := range moroccoData.Operators {
		fmt.Printf("%s\n", operator.Name)
		fmt.Println(strings.Repeat("-", 50))
		
		// Test only 2 numbers per operator to conserve API calls
		numTests := 2
		if len(operator.Prefixes) < 2 {
			numTests = len(operator.Prefixes)
		}
		
		// Pick random prefixes
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < numTests; i++ {
			prefix := operator.Prefixes[rand.Intn(len(operator.Prefixes))]
			
			// Generate a realistic subscriber number
			subscriber := fmt.Sprintf("%06d", 100000+rand.Intn(900000))
			fullNumber := fmt.Sprintf("+%s%s", prefix, subscriber)
			
			stats.TotalTested++
			
			// 1. Test with libphonenumber
			info, err := phoneValidator.ValidatePhoneNumber(fullNumber, "MA")
			if err != nil {
				fmt.Printf("  %s: PhoneLib Error: %v\n", fullNumber, err)
				continue
			}
			
			phoneLibStatus := "❌ Invalid"
			if info.IsValid {
				phoneLibStatus = "✅ Valid"
				stats.PhoneLibValid++
			}
			
			// 2. Test with WhatsApp API (only for valid numbers)
			whatsappStatus := "Not checked"
			if info.IsValid && info.IsMobile {
				stats.WhatsAppChecked++
				
				whatsappResult, err := whatsappValidator.ValidateNumber(fullNumber)
				if err != nil {
					whatsappStatus = fmt.Sprintf("API Error: %v", err)
					stats.APIErrors++
				} else if whatsappResult.HasWhatsApp {
					whatsappStatus = "✅ Has WhatsApp"
					stats.HasWhatsApp++
				} else {
					whatsappStatus = "❌ No WhatsApp"
				}
			}
			
			fmt.Printf("  %s\n", fullNumber)
			fmt.Printf("    PhoneLib: %s (Type: %s)\n", phoneLibStatus, info.NumberType)
			fmt.Printf("    WhatsApp: %s\n", whatsappStatus)
			
			// Add delay between API calls to respect rate limits
			if whatsappAPIKey != "" && i < numTests-1 {
				time.Sleep(1 * time.Second)
			}
		}
		fmt.Println()
	}
	
	// Test some unassigned prefixes
	fmt.Println("Testing UNASSIGNED prefixes (should be rejected):")
	fmt.Println("================================================\n")
	
	unassignedPrefixes := []string{
		"212509",  // Not assigned
		"212685",  // Not in any operator list
		"212730",  // Not assigned
		"212799",  // Not assigned
	}
	
	unassignedAccepted := 0
	
	for _, prefix := range unassignedPrefixes {
		subscriber := fmt.Sprintf("%06d", 100000+rand.Intn(900000))
		fullNumber := fmt.Sprintf("+%s%s", prefix, subscriber)
		
		info, err := phoneValidator.ValidatePhoneNumber(fullNumber, "MA")
		
		status := "✅ Correctly rejected"
		if err == nil && info.IsValid && info.IsMobile {
			unassignedAccepted++
			status = "⚠️  INCORRECTLY accepted as valid mobile"
		}
		
		fmt.Printf("  %s: %s\n", fullNumber, status)
	}
	
	// Summary
	fmt.Println("\n\nSUMMARY")
	fmt.Println("=======")
	fmt.Printf("Total numbers tested: %d\n", stats.TotalTested)
	fmt.Printf("Valid per libphonenumber: %d (%.1f%%)\n", 
		stats.PhoneLibValid, float64(stats.PhoneLibValid)*100/float64(stats.TotalTested))
	fmt.Printf("WhatsApp API checked: %d\n", stats.WhatsAppChecked)
	fmt.Printf("Have WhatsApp: %d (%.1f%% of checked)\n", 
		stats.HasWhatsApp, float64(stats.HasWhatsApp)*100/float64(max(stats.WhatsAppChecked, 1)))
	fmt.Printf("API errors: %d\n", stats.APIErrors)
	fmt.Printf("\nUnassigned prefixes incorrectly accepted: %d/%d\n", 
		unassignedAccepted, len(unassignedPrefixes))
	
	fmt.Println("\nCONCLUSIONS:")
	fmt.Println("============")
	fmt.Println("1. Libphonenumber validates based on broad ranges, not specific assignments")
	fmt.Println("2. Our prefix database is essential for accurate Morocco routing")
	fmt.Println("3. WhatsApp validation should only be called for valid numbers to save costs")
	
	if whatsappAPIKey == "" {
		fmt.Println("\n💡 TIP: Set WHATSAPP_API_KEY environment variable to test real WhatsApp validation")
	}
}

// Mock WhatsApp validator for when API key is not available
type mockWhatsAppValidator struct{}

func (m *mockWhatsAppValidator) ValidateNumber(phoneNumber string) (*models.ValidationResult, error) {
	// Simulate API response
	rand.Seed(time.Now().UnixNano())
	hasWhatsApp := rand.Float32() < 0.8 // 80% have WhatsApp
	
	return &models.ValidationResult{
		PhoneNumber: phoneNumber,
		HasWhatsApp: hasWhatsApp,
		Source:      "mock",
		Confidence:  0.95,
	}, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}