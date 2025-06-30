package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

// Mock WhatsApp validator for testing
type mockWhatsAppValidator struct{}

func (m *mockWhatsAppValidator) ValidateNumber(phoneNumber string) (*models.ValidationResult, error) {
	// Simulate WhatsApp validation
	// In reality, this would call the actual WhatsApp Business API
	
	// For testing, let's say 80% of valid numbers have WhatsApp
	rand.Seed(time.Now().UnixNano())
	hasWhatsApp := rand.Float32() < 0.8
	
	result := &models.ValidationResult{
		PhoneNumber: phoneNumber,
		HasWhatsApp: hasWhatsApp,
		Source:      "mock_whatsapp_api",
		Confidence:  0.95,
	}
	
	return result, nil
}

func main() {
	fmt.Println("Morocco Prefix Validation Test - PhoneLib vs Reality")
	fmt.Println("===================================================\n")
	
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
	
	whatsappValidator := &mockWhatsAppValidator{}
	
	fmt.Println("Testing random prefixes from each operator:\n")
	
	// Test statistics
	totalTests := 0
	phoneLibValid := 0
	whatsappChecked := 0
	whatsappValid := 0
	
	// Test each operator
	for opKey, operator := range moroccoData.Operators {
		fmt.Printf("%s (%s)\n", operator.Name, opKey)
		fmt.Println(strings.Repeat("-", 50))
		
		// Pick 5 random prefixes from this operator
		numTests := 5
		if len(operator.Prefixes) < 5 {
			numTests = len(operator.Prefixes)
		}
		
		// Shuffle and pick first N
		prefixesCopy := make([]string, len(operator.Prefixes))
		copy(prefixesCopy, operator.Prefixes)
		rand.Shuffle(len(prefixesCopy), func(i, j int) {
			prefixesCopy[i], prefixesCopy[j] = prefixesCopy[j], prefixesCopy[i]
		})
		
		for i := 0; i < numTests; i++ {
			prefix := prefixesCopy[i]
			
			// Generate random subscriber number
			subscriber := fmt.Sprintf("%06d", rand.Intn(1000000))
			fullNumber := fmt.Sprintf("+%s%s", prefix, subscriber)
			
			totalTests++
			
			// Test with libphonenumber
			info, err := phoneValidator.ValidatePhoneNumber(fullNumber, "MA")
			
			phoneLibResult := "❌ Invalid"
			if err == nil && info.IsValid {
				phoneLibResult = "✅ Valid"
				phoneLibValid++
				
				// Only check WhatsApp for valid numbers
				whatsappResult, err := whatsappValidator.ValidateNumber(fullNumber)
				whatsappChecked++
				
				whatsappStatus := "❌ No WhatsApp"
				if err == nil && whatsappResult.HasWhatsApp {
					whatsappStatus = "✅ Has WhatsApp"
					whatsappValid++
				}
				
				fmt.Printf("  %s: PhoneLib=%s, Type=%s, WhatsApp=%s\n",
					fullNumber, phoneLibResult, info.NumberType, whatsappStatus)
			} else {
				fmt.Printf("  %s: PhoneLib=%s (Should be valid!)\n",
					fullNumber, phoneLibResult)
			}
		}
		fmt.Println()
	}
	
	// Test invalid/unassigned prefixes
	fmt.Println("\nTesting INVALID/UNASSIGNED prefixes:")
	fmt.Println(strings.Repeat("-", 50))
	
	invalidPrefixes := []string{
		"212509",  // Not assigned
		"212519",  // Not assigned
		"212559",  // Not assigned
		"212609",  // Not assigned (we know 600-609 are Inwi)
		"212679",  // Let's check if this is assigned
		"212685",  // Not in any list
		"212710",  // Gap in 71X range  
		"212730",  // Not assigned
		"212740",  // Not assigned
		"212789",  // Not assigned
		"212799",  // Not assigned
		"212800",  // Outside mobile range
	}
	
	invalidAccepted := 0
	
	for _, prefix := range invalidPrefixes {
		subscriber := fmt.Sprintf("%06d", rand.Intn(1000000))
		fullNumber := fmt.Sprintf("+%s%s", prefix, subscriber)
		
		// Check if this prefix is actually in our database
		isActuallyValid := false
		for _, op := range moroccoData.Operators {
			for _, validPrefix := range op.Prefixes {
				if prefix == validPrefix {
					isActuallyValid = true
					break
				}
			}
			if isActuallyValid {
				break
			}
		}
		
		info, err := phoneValidator.ValidatePhoneNumber(fullNumber, "MA")
		
		if err == nil && info.IsValid && info.IsMobile {
			invalidAccepted++
			if isActuallyValid {
				fmt.Printf("  %s: ✅ Valid (prefix IS assigned)\n", fullNumber)
			} else {
				fmt.Printf("  %s: ⚠️  PhoneLib accepts as valid (prefix NOT assigned!)\n", fullNumber)
			}
		} else {
			fmt.Printf("  %s: ✅ Correctly rejected\n", fullNumber)
		}
	}
	
	// Summary
	fmt.Println("\n\nSUMMARY")
	fmt.Println("=======")
	fmt.Printf("Valid Morocco prefixes tested: %d\n", totalTests)
	fmt.Printf("Accepted by libphonenumber: %d (%.1f%%)\n", 
		phoneLibValid, float64(phoneLibValid)*100/float64(totalTests))
	fmt.Printf("WhatsApp checked: %d\n", whatsappChecked)
	fmt.Printf("Have WhatsApp: %d (%.1f%% of checked)\n", 
		whatsappValid, float64(whatsappValid)*100/float64(whatsappChecked))
	fmt.Printf("\nInvalid prefixes tested: %d\n", len(invalidPrefixes))
	fmt.Printf("Incorrectly accepted by libphonenumber: %d\n", invalidAccepted)
	
	// Conclusion
	fmt.Println("\nCONCLUSION:")
	if phoneLibValid == totalTests {
		fmt.Println("✅ Libphonenumber correctly validates all assigned Morocco prefixes")
	} else {
		fmt.Printf("⚠️  Libphonenumber failed to validate %d assigned prefixes\n", 
			totalTests-phoneLibValid)
	}
	
	if invalidAccepted > 0 {
		fmt.Println("⚠️  Libphonenumber accepts some unassigned prefixes as valid")
		fmt.Println("    This confirms it uses broad ranges, not specific prefix assignments")
	}
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}