package main

import (
	"fmt"
	"log"
	
	"github.com/e173-gateway/e173_go_gateway/pkg/validation"
)

func main() {
	fmt.Println("Morocco Number Length Test")
	fmt.Println("==========================\n")
	
	validator, err := validation.NewLibPhoneNumberValidator()
	if err != nil {
		log.Fatalf("Failed to initialize libphonenumber: %v", err)
	}
	
	// Test different lengths for a known valid prefix (661 - IAM)
	prefix := "212661"
	
	fmt.Println("Testing different lengths for prefix +212661 (IAM):")
	fmt.Println("---------------------------------------------------")
	
	// Test with different suffix lengths
	suffixes := []string{
		"",          // No suffix
		"1",         // 1 digit
		"12",        // 2 digits
		"123",       // 3 digits
		"1234",      // 4 digits
		"12345",     // 5 digits
		"123456",    // 6 digits - standard Morocco mobile
		"1234567",   // 7 digits
		"12345678",  // 8 digits
		"123456789", // 9 digits
	}
	
	for _, suffix := range suffixes {
		number := fmt.Sprintf("+%s%s", prefix, suffix)
		info, err := validator.ValidatePhoneNumber(number, "MA")
		
		totalDigits := len(prefix) + len(suffix) - 3 // -3 for country code
		
		if err != nil {
			fmt.Printf("❌ %s (length: %d) - Error: %v\n", number, totalDigits, err)
			continue
		}
		
		status := "❌"
		if info.IsValid {
			status = "✅"
		}
		
		fmt.Printf("%s %s (length: %d) - Valid: %v, Mobile: %v, Type: %s\n",
			status, number, totalDigits, info.IsValid, info.IsMobile, info.NumberType)
	}
	
	// Test actual Morocco number formats
	fmt.Println("\n\nTesting real Morocco number formats:")
	fmt.Println("------------------------------------")
	
	realNumbers := []struct {
		number      string
		description string
	}{
		{"+212661234567", "Standard format - 9 digits after country code"},
		{"+2126612345678", "10 digits after country code"},
		{"+212661234567890", "13 digits after country code"},
		{"+212609123456", "Invalid prefix 609 but correct length"},
		{"+212710123456", "Gap prefix 710 but correct length"},
		{"+212778123456", "Gap prefix 778 but correct length"},
		{"+212600123456", "Valid Inwi prefix 600"},
		{"+212620123456", "Valid Orange prefix 620"},
		{"+212770123456", "Valid Orange prefix 770"},
	}
	
	for _, test := range realNumbers {
		info, err := validator.ValidatePhoneNumber(test.number, "MA")
		
		if err != nil {
			fmt.Printf("❌ %s - %s\n   Error: %v\n", test.number, test.description, err)
			continue
		}
		
		status := "❌"
		if info.IsValid {
			status = "✅"
		}
		
		fmt.Printf("%s %s - %s\n", status, test.number, test.description)
		fmt.Printf("   Valid: %v, Mobile: %v, Type: %s, Formatted: %s\n",
			info.IsValid, info.IsMobile, info.NumberType, info.Formatted)
	}
	
	// Test if libphonenumber knows about Morocco mobile number plan
	fmt.Println("\n\nLibphonenumber's Morocco mobile number understanding:")
	fmt.Println("----------------------------------------------------")
	
	// Test edge cases
	testPrefixes := []string{
		"600", "609",  // Inwi range
		"610", "619",  // IAM range  
		"620", "629",  // Orange range
		"630", "639",  // Mixed
		"700", "709",  // Inwi 70X
		"710", "719",  // Gap
		"770", "779",  // Orange 77X
		"780", "789",  // Gap
		"800", "900",  // Outside mobile
	}
	
	for i := 0; i < len(testPrefixes); i += 2 {
		start := testPrefixes[i]
		end := testPrefixes[i+1]
		
		// Test start and end of range
		for _, prefix := range []string{start, end} {
			number := fmt.Sprintf("+212%s123456", prefix)
			isValid := validator.IsValid(number)
			isMobile := validator.IsValidMobile(number)
			
			mobileStatus := ""
			if isValid && isMobile {
				mobileStatus = "MOBILE"
			} else if isValid && !isMobile {
				mobileStatus = "FIXED/OTHER"
			} else {
				mobileStatus = "INVALID"
			}
			
			fmt.Printf("212%s: %s\n", prefix, mobileStatus)
		}
		fmt.Println()
	}
}