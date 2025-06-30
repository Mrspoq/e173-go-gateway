package validation

// #cgo CFLAGS: -I./cpp
// #cgo LDFLAGS: -L./cpp -lphonenumber_wrapper -lphonenumber -lgeocoding -lprotobuf -lboost_system -lboost_thread -licui18n -licuuc -licudata -lstdc++
// #include "phonenumber_wrapper.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"strings"
	"sync"
	"unsafe"
)

// LibPhoneNumberValidator uses Google's official libphonenumber C++ library
type LibPhoneNumberValidator struct {
	initialized bool
	mu          sync.Mutex
}

// Initialize the libphonenumber library (call once)
var libPhoneNumberOnce sync.Once
var libPhoneNumberInitErr error

// NewLibPhoneNumberValidator creates a new validator using the official libphonenumber
func NewLibPhoneNumberValidator() (*LibPhoneNumberValidator, error) {
	validator := &LibPhoneNumberValidator{}
	
	// Initialize library once
	libPhoneNumberOnce.Do(func() {
		if C.phone_lib_init() != 1 {
			libPhoneNumberInitErr = fmt.Errorf("failed to initialize libphonenumber")
		}
	})
	
	if libPhoneNumberInitErr != nil {
		return nil, libPhoneNumberInitErr
	}
	
	validator.initialized = true
	return validator, nil
}

// ValidatePhoneNumber validates and parses a phone number
func (l *LibPhoneNumberValidator) ValidatePhoneNumber(phoneNumber string, defaultRegion ...string) (*PhoneNumberInfo, error) {
	if !l.initialized {
		return nil, fmt.Errorf("libphonenumber not initialized")
	}
	
	region := "US"
	if len(defaultRegion) > 0 && defaultRegion[0] != "" {
		region = defaultRegion[0]
	}
	
	// Convert to C strings
	cPhoneNumber := C.CString(phoneNumber)
	cRegion := C.CString(region)
	defer C.free(unsafe.Pointer(cPhoneNumber))
	defer C.free(unsafe.Pointer(cRegion))
	
	// Create result struct
	var result C.PhoneValidationResult
	
	// Validate the number
	if C.validate_phone_number(cRegion, cPhoneNumber, &result) != 1 {
		errorMsg := "unknown error"
		if result.error_msg != nil {
			errorMsg = C.GoString(result.error_msg)
		}
		C.free_validation_result(&result)
		return nil, fmt.Errorf("validation failed: %s", errorMsg)
	}
	
	// Convert result to Go struct
	info := &PhoneNumberInfo{
		Original:       phoneNumber,
		IsValid:        result.is_valid == 1,
		IsMobile:       result.is_mobile == 1,
	}
	
	if result.formatted_number != nil {
		info.Formatted = C.GoString(result.formatted_number)
	}
	if result.country_code != nil {
		info.CountryCode = C.GoString(result.country_code)
	}
	if result.national_number != nil {
		info.NationalNumber = C.GoString(result.national_number)
	}
	if result.carrier != nil {
		info.Carrier = C.GoString(result.carrier)
	} else {
		info.Carrier = "UNKNOWN"
	}
	if result.region != nil {
		info.Region = C.GoString(result.region)
	}
	if result.number_type != nil {
		info.NumberType = C.GoString(result.number_type)
	}
	
	// Free the result
	C.free_validation_result(&result)
	
	return info, nil
}

// IsValid implements PhoneNumberValidator interface
func (l *LibPhoneNumberValidator) IsValid(phoneNumber string) bool {
	if !l.initialized {
		return false
	}
	
	// Try to detect region from number
	region := l.detectRegion(phoneNumber)
	
	cPhoneNumber := C.CString(phoneNumber)
	cRegion := C.CString(region)
	defer C.free(unsafe.Pointer(cPhoneNumber))
	defer C.free(unsafe.Pointer(cRegion))
	
	return C.is_valid_number(cRegion, cPhoneNumber) == 1
}

// IsValidMobile checks if a number is a valid mobile number
func (l *LibPhoneNumberValidator) IsValidMobile(phoneNumber string) bool {
	if !l.initialized {
		return false
	}
	
	region := l.detectRegion(phoneNumber)
	
	cPhoneNumber := C.CString(phoneNumber)
	cRegion := C.CString(region)
	defer C.free(unsafe.Pointer(cPhoneNumber))
	defer C.free(unsafe.Pointer(cRegion))
	
	return C.is_valid_mobile_number(cRegion, cPhoneNumber) == 1
}

// FormatInternational formats a number in international format
func (l *LibPhoneNumberValidator) FormatInternational(phoneNumber string) string {
	if !l.initialized {
		return phoneNumber
	}
	
	region := l.detectRegion(phoneNumber)
	
	cPhoneNumber := C.CString(phoneNumber)
	cRegion := C.CString(region)
	defer C.free(unsafe.Pointer(cPhoneNumber))
	defer C.free(unsafe.Pointer(cRegion))
	
	cFormatted := C.format_international(cRegion, cPhoneNumber)
	if cFormatted == nil {
		return phoneNumber
	}
	
	formatted := C.GoString(cFormatted)
	C.free_string(cFormatted)
	
	return formatted
}

// detectRegion tries to detect the region from the phone number
func (l *LibPhoneNumberValidator) detectRegion(phoneNumber string) string {
	// Clean the number
	cleaned := strings.TrimSpace(phoneNumber)
	
	// Check for common country codes
	if strings.HasPrefix(cleaned, "+212") || strings.HasPrefix(cleaned, "00212") {
		return "MA" // Morocco
	}
	if strings.HasPrefix(cleaned, "+234") || strings.HasPrefix(cleaned, "00234") {
		return "NG" // Nigeria
	}
	if strings.HasPrefix(cleaned, "+1") {
		return "US" // US/Canada
	}
	if strings.HasPrefix(cleaned, "+44") {
		return "GB" // UK
	}
	if strings.HasPrefix(cleaned, "+33") {
		return "FR" // France
	}
	if strings.HasPrefix(cleaned, "+49") {
		return "DE" // Germany
	}
	if strings.HasPrefix(cleaned, "+254") {
		return "KE" // Kenya
	}
	if strings.HasPrefix(cleaned, "+27") {
		return "ZA" // South Africa
	}
	if strings.HasPrefix(cleaned, "+91") {
		return "IN" // India
	}
	if strings.HasPrefix(cleaned, "+86") {
		return "CN" // China
	}
	
	// Default to US
	return "US"
}