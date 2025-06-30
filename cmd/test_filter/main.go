package main

import (
	"fmt"
	"log"
	
	"github.com/e173-gateway/e173_go_gateway/pkg/models"
	"github.com/e173-gateway/e173_go_gateway/pkg/service"
	"github.com/e173-gateway/e173_go_gateway/pkg/validation"
)

// Mock repositories for testing
type mockBlacklistRepo struct{}

func (m *mockBlacklistRepo) Add(number string, reason string) error {
	return nil
}

func (m *mockBlacklistRepo) Remove(number string) error {
	return nil
}

func (m *mockBlacklistRepo) IsBlacklisted(number string) (bool, error) {
	// Example blacklisted numbers
	blacklisted := map[string]bool{
		"+2341234567890": true,
		"+1234567890":    true,
	}
	return blacklisted[number], nil
}

func (m *mockBlacklistRepo) GetAll() ([]models.BlacklistEntry, error) {
	return nil, nil
}

func (m *mockBlacklistRepo) GetByNumber(number string) (*models.BlacklistEntry, error) {
	return nil, nil
}

type mockPrefixRepo struct{}

func (m *mockPrefixRepo) Create(prefix *models.Prefix) error {
	return nil
}

func (m *mockPrefixRepo) GetByID(id string) (*models.Prefix, error) {
	return nil, nil
}

func (m *mockPrefixRepo) GetByPrefix(prefix string) (*models.Prefix, error) {
	return nil, nil
}

func (m *mockPrefixRepo) GetAllActive() ([]models.Prefix, error) {
	// Return some test prefixes
	return []models.Prefix{
		{Prefix: "212", GatewayID: "gw-morocco", IsActive: true},    // Morocco
		{Prefix: "234", GatewayID: "gw-nigeria", IsActive: true},    // Nigeria
		{Prefix: "1", GatewayID: "gw-us", IsActive: true},          // US/Canada
		{Prefix: "44", GatewayID: "gw-uk", IsActive: true},         // UK
	}, nil
}

func (m *mockPrefixRepo) Update(prefix *models.Prefix) error {
	return nil
}

func (m *mockPrefixRepo) Delete(id string) error {
	return nil
}

// Mock WhatsApp validator for testing
type mockWhatsAppValidator struct{}

func (m *mockWhatsAppValidator) ValidateNumber(phoneNumber string) (*models.ValidationResult, error) {
	// Simulate some numbers having WhatsApp
	hasWhatsApp := map[string]bool{
		"+212661234567":  true,  // Morocco mobile
		"+2348012345678": true,  // Nigeria mobile
		"+14155552222":   true,  // US mobile
		"+447911123456":  false, // UK mobile (no WhatsApp)
	}
	
	result := &models.ValidationResult{
		PhoneNumber: phoneNumber,
		HasWhatsApp: hasWhatsApp[phoneNumber],
		Source:      "mock_whatsapp",
	}
	
	return result, nil
}

func main() {
	fmt.Println("E173 Gateway Filter Test")
	fmt.Println("========================")
	
	// Initialize libphonenumber validator
	phoneValidator, err := validation.NewLibPhoneNumberValidator()
	if err != nil {
		log.Fatalf("Failed to initialize libphonenumber: %v", err)
	}
	fmt.Println("✓ Libphonenumber initialized successfully")
	
	// Create filter service with mocks
	filterService := service.NewFilterService(
		&mockBlacklistRepo{},
		&mockPrefixRepo{},
		&mockWhatsAppValidator{},
		phoneValidator,
	)
	
	// Test cases
	testCases := []struct {
		name   string
		source string
		dest   string
	}{
		{"Valid Morocco to Morocco", "+212661234567", "+212661234567"},
		{"Valid Nigeria to US", "+2348012345678", "+14155552222"},
		{"Invalid source format", "12345", "+14155552222"},
		{"Invalid destination format", "+2348012345678", "invalid"},
		{"Blacklisted source", "+2341234567890", "+14155552222"},
		{"No WhatsApp on destination", "+2348012345678", "+447911123456"},
		{"No route for destination", "+2348012345678", "+999123456789"},
	}
	
	fmt.Println("\nRunning test cases:")
	fmt.Println("-------------------")
	
	for _, tc := range testCases {
		fmt.Printf("\nTest: %s\n", tc.name)
		fmt.Printf("Source: %s → Destination: %s\n", tc.source, tc.dest)
		
		call := &models.Call{
			SourceNumber: tc.source,
			DestNumber:   tc.dest,
		}
		
		result, err := filterService.ProcessCall(call)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}
		
		switch result.Action {
		case "route":
			fmt.Printf("✅ ROUTE via gateway %s (prefix: %s)\n", result.GatewayID, result.Prefix)
		case "reject":
			fmt.Printf("🚫 REJECT: %s\n", result.Reason)
		case "blackhole":
			fmt.Printf("⛔ BLACKHOLE: %s\n", result.Reason)
		}
	}
	
	// Test libphonenumber directly
	fmt.Println("\n\nDirect libphonenumber validation tests:")
	fmt.Println("--------------------------------------")
	
	testNumbers := []string{
		"+212661234567",  // Morocco
		"+2348012345678", // Nigeria
		"+14155552222",   // US
		"+447911123456",  // UK
		"invalid",        // Invalid
	}
	
	for _, number := range testNumbers {
		info, err := phoneValidator.ValidatePhoneNumber(number)
		if err != nil {
			fmt.Printf("\n❌ %s: %v\n", number, err)
			continue
		}
		
		fmt.Printf("\n📱 %s:\n", number)
		fmt.Printf("  Valid: %v\n", info.IsValid)
		fmt.Printf("  Mobile: %v\n", info.IsMobile)
		fmt.Printf("  Formatted: %s\n", info.Formatted)
		fmt.Printf("  Country Code: %s\n", info.CountryCode)
		fmt.Printf("  Region: %s\n", info.Region)
		fmt.Printf("  Type: %s\n", info.NumberType)
	}
}