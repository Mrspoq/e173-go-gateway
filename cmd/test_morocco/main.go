package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	
	"github.com/e173-gateway/e173_go_gateway/pkg/validation"
)

type MoroccoData struct {
	Operators map[string]struct {
		Name     string   `json:"name"`
		Prefixes []string `json:"prefixes"`
	} `json:"operators"`
}

func main() {
	fmt.Println("Morocco Phone Number Validation Test")
	fmt.Println("====================================")
	
	// Load Morocco prefixes data
	data, err := ioutil.ReadFile("/root/e173_go_gateway/data/morocco_mobile_prefixes.json")
	if err != nil {
		log.Fatalf("Failed to read Morocco data: %v", err)
	}
	
	var moroccoData MoroccoData
	if err := json.Unmarshal(data, &moroccoData); err != nil {
		log.Fatalf("Failed to parse Morocco data: %v", err)
	}
	
	// Initialize libphonenumber validator
	validator, err := validation.NewLibPhoneNumberValidator()
	if err != nil {
		log.Fatalf("Failed to initialize libphonenumber: %v", err)
	}
	fmt.Println("✓ Libphonenumber initialized successfully\n")
	
	// Test numbers from each operator
	fmt.Println("Testing sample numbers from each operator:")
	fmt.Println("==========================================")
	
	for operatorCode, operator := range moroccoData.Operators {
		fmt.Printf("\n%s (%s):\n", operator.Name, operatorCode)
		
		// Test first 5 prefixes from each operator
		count := 0
		validCount := 0
		mobileCount := 0
		
		for i, prefix := range operator.Prefixes {
			if i >= 5 {
				break
			}
			
			// Create a test number with the prefix + 6 random digits
			testNumber := fmt.Sprintf("+%s123456", prefix)
			
			// Validate the number
			info, err := validator.ValidatePhoneNumber(testNumber, "MA")
			if err != nil {
				fmt.Printf("  ❌ %s: Error - %v\n", testNumber, err)
				continue
			}
			
			count++
			if info.IsValid {
				validCount++
			}
			if info.IsMobile {
				mobileCount++
			}
			
			status := "❌"
			if info.IsValid {
				status = "✅"
			}
			
			fmt.Printf("  %s %s: Valid=%v Mobile=%v Formatted='%s' Type=%s\n",
				status, testNumber, info.IsValid, info.IsMobile, 
				info.Formatted, info.NumberType)
		}
		
		fmt.Printf("  Summary: %d/%d valid, %d/%d mobile\n", 
			validCount, count, mobileCount, count)
	}
	
	// Test specific known valid Morocco numbers
	fmt.Println("\n\nTesting specific Morocco number formats:")
	fmt.Println("=========================================")
	
	testCases := []struct {
		number      string
		description string
	}{
		{"+212661234567", "IAM mobile (661 prefix)"},
		{"+212 661 234 567", "IAM mobile with spaces"},
		{"00212661234567", "IAM mobile with 00 prefix"},
		{"+212620123456", "Orange mobile (620 prefix)"},
		{"+212700123456", "Inwi mobile (700 prefix)"},
		{"+212523456789", "Landline (Casablanca)"},
		{"+212537123456", "Landline (Rabat)"},
		{"212661234567", "No + prefix"},
		{"0661234567", "Local format"},
		{"+21266123456", "Too short"},
		{"+2126612345678", "Too long"},
	}
	
	for _, tc := range testCases {
		info, err := validator.ValidatePhoneNumber(tc.number, "MA")
		if err != nil {
			fmt.Printf("\n❌ %s (%s)\n   Error: %v\n", tc.number, tc.description, err)
			continue
		}
		
		status := "❌"
		if info.IsValid {
			status = "✅"
		}
		
		fmt.Printf("\n%s %s (%s)\n", status, tc.number, tc.description)
		fmt.Printf("   Valid: %v\n", info.IsValid)
		fmt.Printf("   Mobile: %v\n", info.IsMobile)
		fmt.Printf("   Formatted: %s\n", info.Formatted)
		fmt.Printf("   Country Code: %s\n", info.CountryCode)
		fmt.Printf("   Region: %s\n", info.Region)
		fmt.Printf("   Type: %s\n", info.NumberType)
	}
	
	// Test validation performance
	fmt.Println("\n\nValidation Performance Test:")
	fmt.Println("============================")
	
	// Test IsValid method (simpler interface)
	quickTests := []string{
		"+212661234567",
		"+212620123456", 
		"+212700123456",
	}
	
	for _, number := range quickTests {
		isValid := validator.IsValid(number)
		isMobile := validator.IsValidMobile(number)
		formatted := validator.FormatInternational(number)
		
		fmt.Printf("\nQuick test for %s:\n", number)
		fmt.Printf("  IsValid: %v\n", isValid)
		fmt.Printf("  IsValidMobile: %v\n", isMobile)
		fmt.Printf("  FormatInternational: %s\n", formatted)
	}
	
	// Show prefix distribution
	fmt.Println("\n\nMorocco Operator Prefix Summary:")
	fmt.Println("=================================")
	for _, operator := range moroccoData.Operators {
		fmt.Printf("%s: %d prefixes\n", operator.Name, len(operator.Prefixes))
		
		// Group prefixes by first 5 digits
		groups := make(map[string]int)
		for _, prefix := range operator.Prefixes {
			if len(prefix) >= 5 {
				group := prefix[:5]
				groups[group]++
			}
		}
		
		// Show groups
		var sortedGroups []string
		for group := range groups {
			sortedGroups = append(sortedGroups, group)
		}
		
		fmt.Printf("  Prefix groups: ")
		for i, group := range sortedGroups {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s (%d)", strings.TrimPrefix(group, "212"), groups[group])
		}
		fmt.Println()
	}
}