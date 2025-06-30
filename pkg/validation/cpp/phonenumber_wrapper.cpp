#include "phonenumber_wrapper.h"
#include <phonenumbers/phonenumberutil.h>
#include <phonenumbers/phonenumber.pb.h>
#include <phonenumbers/geocoding/phonenumber_offline_geocoder.h>
#include <cstring>
#include <string>
#include <memory>

using namespace i18n::phonenumbers;

static std::unique_ptr<PhoneNumberUtil> phone_util;
static std::unique_ptr<PhoneNumberOfflineGeocoder> geocoder;

// Helper function to convert string to C string
char* string_to_cstring(const std::string& str) {
    char* cstr = new char[str.length() + 1];
    std::strcpy(cstr, str.c_str());
    return cstr;
}

// Initialize the library
int phone_lib_init() {
    try {
        phone_util.reset(PhoneNumberUtil::GetInstance());
        geocoder.reset(new PhoneNumberOfflineGeocoder());
        return 1;
    } catch (...) {
        return 0;
    }
}

// Validate a phone number
int validate_phone_number(const char* region_code, const char* phone_number, PhoneValidationResult* result) {
    if (!phone_util || !result) {
        return 0;
    }
    
    // Initialize result
    memset(result, 0, sizeof(PhoneValidationResult));
    
    try {
        std::string number_str(phone_number);
        std::string region(region_code ? region_code : "US");
        
        PhoneNumber parsed_number;
        PhoneNumberUtil::ErrorType error = phone_util->Parse(number_str, region, &parsed_number);
        
        if (error != PhoneNumberUtil::NO_PARSING_ERROR) {
            result->is_valid = 0;
            result->error_msg = string_to_cstring("Failed to parse number");
            return 1;
        }
        
        // Check if valid
        result->is_valid = phone_util->IsValidNumber(parsed_number) ? 1 : 0;
        result->is_possible = phone_util->IsPossibleNumber(parsed_number) ? 1 : 0;
        
        // Get number type
        PhoneNumberUtil::PhoneNumberType type = phone_util->GetNumberType(parsed_number);
        result->is_mobile = (type == PhoneNumberUtil::MOBILE || 
                           type == PhoneNumberUtil::FIXED_LINE_OR_MOBILE) ? 1 : 0;
        
        // Format number
        std::string formatted;
        phone_util->Format(parsed_number, PhoneNumberUtil::INTERNATIONAL, &formatted);
        result->formatted_number = string_to_cstring(formatted);
        
        // Get country code
        result->country_code = string_to_cstring(std::to_string(parsed_number.country_code()));
        
        // Get national number
        result->national_number = string_to_cstring(std::to_string(parsed_number.national_number()));
        
        // Get region
        std::string region_code_str;
        phone_util->GetRegionCodeForNumber(parsed_number, &region_code_str);
        result->region = string_to_cstring(region_code_str);
        
        // Carrier info not available in basic installation
        result->carrier = nullptr;
        
        // Get number type string
        const char* type_str = "UNKNOWN";
        switch (type) {
            case PhoneNumberUtil::FIXED_LINE:
                type_str = "FIXED_LINE";
                break;
            case PhoneNumberUtil::MOBILE:
                type_str = "MOBILE";
                break;
            case PhoneNumberUtil::FIXED_LINE_OR_MOBILE:
                type_str = "FIXED_LINE_OR_MOBILE";
                break;
            case PhoneNumberUtil::TOLL_FREE:
                type_str = "TOLL_FREE";
                break;
            case PhoneNumberUtil::PREMIUM_RATE:
                type_str = "PREMIUM_RATE";
                break;
            case PhoneNumberUtil::SHARED_COST:
                type_str = "SHARED_COST";
                break;
            case PhoneNumberUtil::VOIP:
                type_str = "VOIP";
                break;
            case PhoneNumberUtil::PERSONAL_NUMBER:
                type_str = "PERSONAL_NUMBER";
                break;
            case PhoneNumberUtil::PAGER:
                type_str = "PAGER";
                break;
            case PhoneNumberUtil::UAN:
                type_str = "UAN";
                break;
            case PhoneNumberUtil::VOICEMAIL:
                type_str = "VOICEMAIL";
                break;
            default:
                type_str = "UNKNOWN";
        }
        result->number_type = string_to_cstring(type_str);
        
        return 1;
    } catch (const std::exception& e) {
        result->error_msg = string_to_cstring(e.what());
        return 0;
    } catch (...) {
        result->error_msg = string_to_cstring("Unknown error");
        return 0;
    }
}

// Free result strings
void free_validation_result(PhoneValidationResult* result) {
    if (!result) return;
    
    delete[] result->formatted_number;
    delete[] result->country_code;
    delete[] result->national_number;
    delete[] result->carrier;
    delete[] result->region;
    delete[] result->number_type;
    delete[] result->error_msg;
    
    memset(result, 0, sizeof(PhoneValidationResult));
}

// Check if number is valid
int is_valid_number(const char* region_code, const char* phone_number) {
    if (!phone_util) return 0;
    
    try {
        std::string number_str(phone_number);
        std::string region(region_code ? region_code : "US");
        
        PhoneNumber parsed_number;
        if (phone_util->Parse(number_str, region, &parsed_number) != PhoneNumberUtil::NO_PARSING_ERROR) {
            return 0;
        }
        
        return phone_util->IsValidNumber(parsed_number) ? 1 : 0;
    } catch (...) {
        return 0;
    }
}

// Check if number is valid mobile number
int is_valid_mobile_number(const char* region_code, const char* phone_number) {
    if (!phone_util) return 0;
    
    try {
        std::string number_str(phone_number);
        std::string region(region_code ? region_code : "US");
        
        PhoneNumber parsed_number;
        if (phone_util->Parse(number_str, region, &parsed_number) != PhoneNumberUtil::NO_PARSING_ERROR) {
            return 0;
        }
        
        PhoneNumberUtil::PhoneNumberType type = phone_util->GetNumberType(parsed_number);
        return (type == PhoneNumberUtil::MOBILE || 
                type == PhoneNumberUtil::FIXED_LINE_OR_MOBILE) ? 1 : 0;
    } catch (...) {
        return 0;
    }
}

// Format number in international format
char* format_international(const char* region_code, const char* phone_number) {
    if (!phone_util) return nullptr;
    
    try {
        std::string number_str(phone_number);
        std::string region(region_code ? region_code : "US");
        
        PhoneNumber parsed_number;
        if (phone_util->Parse(number_str, region, &parsed_number) != PhoneNumberUtil::NO_PARSING_ERROR) {
            return nullptr;
        }
        
        std::string formatted;
        phone_util->Format(parsed_number, PhoneNumberUtil::INTERNATIONAL, &formatted);
        return string_to_cstring(formatted);
    } catch (...) {
        return nullptr;
    }
}

// Get carrier name for number (not available in basic installation)
char* get_carrier(const char* region_code, const char* phone_number) {
    // Carrier lookup not available without carrier mapper
    return nullptr;
}

// Free string returned by library
void free_string(char* str) {
    delete[] str;
}