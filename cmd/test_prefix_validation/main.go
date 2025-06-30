package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
	fmt.Println("Morocco Prefix Validation Deep Test")
	fmt.Println("===================================")
	fmt.Println("Testing if libphonenumber validates actual operator prefixes or just E.164 format\n")
	
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
	
	rand.Seed(time.Now().UnixNano())
	
	// Test 1: Valid prefixes with random suffixes
	fmt.Println("TEST 1: Valid Morocco prefixes with random suffixes")
	fmt.Println("==================================================")
	
	validTests := 0
	totalValid := 0
	
	// Test 5 random prefixes from each operator
	for _, operator := range moroccoData.Operators {
		fmt.Printf("\n%s:\n", operator.Name)
		
		for i := 0; i < 5 && i < len(operator.Prefixes); i++ {
			prefix := operator.Prefixes[rand.Intn(len(operator.Prefixes))]
			
			// Test with different suffixes
			suffixes := []string{"0000", "4500", "9999", "1234", "5678"}
			
			for _, suffix := range suffixes {
				number := fmt.Sprintf("+%s%s", prefix, suffix)
				info, err := validator.ValidatePhoneNumber(number, "MA")
				
				validTests++
				if err == nil && info.IsValid {
					totalValid++
					fmt.Printf("  ✅ %s - Valid Mobile\n", number)
				} else {
					fmt.Printf("  ❌ %s - Invalid\n", number)
				}
			}
		}
	}
	
	fmt.Printf("\nResult: %d/%d numbers validated as correct\n", totalValid, validTests)
	
	// Test 2: Invalid/non-existent prefixes
	fmt.Println("\n\nTEST 2: Invalid/non-existent Morocco prefixes")
	fmt.Println("=============================================")
	
	// Create fake prefixes that don't exist in Morocco
	fakePrefixes := []string{
		"212609",  // Not in any operator list
		"212609",  // Gap in Inwi range
		"212630",  // Actually valid for Inwi - control test
		"212680",  // Actually valid for Orange - control test
		"212690",  // Actually valid for Orange - control test
		"212710",  // Not assigned
		"212730",  // Not assigned
		"212740",  // Not assigned
		"212750",  // Not assigned
		"212780",  // Not assigned
		"212790",  // Not assigned
		"212800",  // Not assigned
		"212900",  // Not assigned
		"212555",  // Landline prefix (should be valid but not mobile)
		"212522",  // Casablanca landline
	}
	
	invalidTests := 0
	invalidValid := 0
	invalidMobile := 0
	
	for _, prefix := range fakePrefixes {
		suffix := "123456"
		number := fmt.Sprintf("+%s%s", prefix, suffix)
		info, err := validator.ValidatePhoneNumber(number, "MA")
		
		invalidTests++
		status := "❌ Invalid"
		if err == nil && info.IsValid {
			invalidValid++
			if info.IsMobile {
				invalidMobile++
				status = "✅ Valid Mobile"
			} else {
				status = "☎️  Valid Fixed Line"
			}
		}
		
		fmt.Printf("%s: %s (Type: %s)\n", number, status, info.NumberType)
	}
	
	fmt.Printf("\nResult: %d/%d fake prefixes accepted as valid (%d as mobile)\n", 
		invalidValid, invalidTests, invalidMobile)
	
	// Test 3: Boundary testing - prefixes just outside valid ranges
	fmt.Println("\n\nTEST 3: Boundary testing - numbers just outside valid ranges")
	fmt.Println("===========================================================")
	
	// Test prefixes that are close to valid ones
	boundaryTests := []struct {
		prefix string
		note   string
	}{
		{"212599", "Just before 600 (Inwi starts at 600)"},
		{"212609", "Between Inwi 608 and 610 (IAM)"},
		{"212729", "Just after Inwi 728"},
		{"212759", "Mobile range but unassigned"},
		{"212769", "Between IAM 767 and Orange 770"},
		{"212778", "Just after Orange 777"},
		{"212799", "End of 7XX mobile range"},
		{"212800", "Outside mobile range"},
		{"212699", "Actually valid - Orange (control)"},
		{"212700", "Actually valid - Inwi (control)"},
	}
	
	for _, test := range boundaryTests {
		number := fmt.Sprintf("+%s123456", test.prefix)
		info, err := validator.ValidatePhoneNumber(number, "MA")
		
		status := "❌ Invalid"
		if err == nil && info.IsValid {
			if info.IsMobile {
				status = "✅ Valid Mobile"
			} else {
				status = "☎️  Valid Fixed"
			}
		}
		
		fmt.Printf("%s: %s - %s\n", number, status, test.note)
	}
	
	// Test 4: Check if libphonenumber knows exact prefix boundaries
	fmt.Println("\n\nTEST 4: Exact prefix boundary knowledge test")
	fmt.Println("===========================================")
	
	// Test specific known gaps in Morocco mobile assignments
	gapTests := []struct {
		start  int
		end    int
		prefix string
	}{
		{709, 719, "212"},  // Gap between Inwi 708 and 720
		{729, 760, "212"},  // Gap after Inwi 728 before IAM 761
		{763, 765, "212"},  // Gap in IAM range
		{768, 769, "212"},  // Gap between IAM 767 and Orange 770
		{778, 799, "212"},  // After Orange 777
	}
	
	fmt.Println("\nTesting gaps in Morocco mobile number assignments:")
	for _, gap := range gapTests {
		gapValid := 0
		gapTotal := 0
		
		for i := gap.start; i <= gap.end && gapTotal < 5; i++ {
			number := fmt.Sprintf("+%s%03d123456", gap.prefix, i)
			info, err := validator.ValidatePhoneNumber(number, "MA")
			
			gapTotal++
			if err == nil && info.IsValid && info.IsMobile {
				gapValid++
				fmt.Printf("  ⚠️  %s - Accepted as valid mobile (should be invalid)\n", number)
			}
		}
		
		if gapValid == 0 {
			fmt.Printf("  ✅ Gap %03d-%03d: Correctly rejected all numbers\n", gap.start, gap.end)
		} else {
			fmt.Printf("  ❌ Gap %03d-%03d: Incorrectly accepted %d/%d numbers\n", 
				gap.start, gap.end, gapValid, gapTotal)
		}
	}
	
	// Summary
	fmt.Println("\n\nSUMMARY")
	fmt.Println("=======")
	fmt.Printf("Valid prefixes test: %d/%d accepted\n", totalValid, validTests)
	fmt.Printf("Invalid prefixes test: %d/%d incorrectly accepted\n", invalidValid, invalidTests)
	fmt.Printf("Mobile detection: %d/%d fake prefixes accepted as mobile\n", invalidMobile, invalidTests)
	
	if totalValid == validTests && invalidValid > invalidTests/2 {
		fmt.Println("\n⚠️  CONCLUSION: libphonenumber appears to validate E.164 format")
		fmt.Println("    but NOT specific operator prefix assignments!")
		fmt.Println("    It accepts many invalid prefixes as valid Morocco mobile numbers.")
	} else if totalValid == validTests && invalidValid == 0 {
		fmt.Println("\n✅ CONCLUSION: libphonenumber correctly validates both format")
		fmt.Println("    AND actual operator prefix assignments!")
	} else {
		fmt.Println("\n❓ CONCLUSION: Mixed results - libphonenumber has partial")
		fmt.Println("    knowledge of Morocco operator prefixes.")
	}
}